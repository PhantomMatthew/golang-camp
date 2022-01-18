// Package redis provides a Redis broker
package redis2

import (
	"context"
	ejson "encoding/json"
	"go.uber.org/zap"
	"newline.com/newline/src/common/config"
	"newline.com/newline/src/common/log"

	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/codec"
	"github.com/micro/go-micro/v2/codec/json"
	"github.com/micro/go-micro/v2/config/cmd"
)

//var ErrNil = errors.New("go-redis: nil returned")
//type Error string

func init() {
	cmd.DefaultBrokers["redis"] = NewBroker
}

// publication is an internal publication for the Redis broker.
type publication struct {
	consumerGroup string
	stream        string
	topic         string
	message       *broker.Message
	err           error
}

// Topic returns the topic this publication applies to.
func (p *publication) Topic() string {
	return p.topic
}

// Message returns the broker message of the publication.
func (p *publication) Message() *broker.Message {
	return p.message
}

// Ack sends an acknowledgement to the broker. However this is not supported
// is Redis and therefore this is a no-op.
func (p *publication) Ack() error {
	return nil
}

func (p *publication) Error() error {
	return p.err
}

// subscriber proxies and handles Redis messages as broker publications.
type subscriber struct {
	codec         codec.Marshaler
	client        *redis.Client
	consumerGroup string
	stream        string
	topic         string
	handle        broker.Handler
	opts          broker.SubscribeOptions
}

func (s *subscriber) deleteLastMessages(receiverStreamName string, groupName string, id string) {
	// 查看这个stream的消费组状态
	xinfoCmd := s.client.XInfoGroups(receiverStreamName)
	xinfoGroups, err := xinfoCmd.Result()

	if err == nil && len(xinfoGroups) > 0 {
		// 获得排除当前组后的group组
		var otherGroups []redis.XInfoGroups
		for _, item := range xinfoGroups {
			if item.Name != groupName {
				otherGroups = append(otherGroups, item)
			}
		}
		if len(otherGroups) > 0 {
			// 判断当前消费组接收ID小于其他分组的lastid就删除
			minId, _ := strconv.Atoi(strings.Split(otherGroups[0].LastDeliveredID, "-")[0])
			for i := 1; i < len(otherGroups); i++ {
				newValue, _ := strconv.Atoi(strings.Split(otherGroups[i].LastDeliveredID, "-")[0])
				if minId > newValue {
					minId = newValue
				}
			}
			currentId, _ := strconv.Atoi(strings.Split(id, "-")[0])
			if currentId < minId {
				//当前组为最后消费的组
				delResult := s.client.XDel(receiverStreamName, id)
				//fmt.Println(delResult)
				log.GetLogger().Info("message delete result", zap.String("stream name", receiverStreamName), zap.Int64("result", delResult.Val()))
			}
		}

	}
}

// recv loops to receive new messages from Redis and handle them
// as publications.
func (s *subscriber) recv(stream, consumerGroup, consumer string) {
	// Close the connection once the subscriber stops receiving.
	//defer s.client.Conn().Close()
	//defer s.client.Pipeline().Close()
	defer s.client.Conn().Close()

	log.GetLogger().Info("check stream is existing or not", zap.String("stream name", stream), zap.Int64("stream_exists_result", s.client.Exists(stream).Val()))

	//Check stream exists or not
	if isExists := s.client.Exists(stream).Val(); isExists == 0 {
		id := s.client.XAdd(&redis.XAddArgs{
			Stream:       stream,
			MaxLen:       0,
			MaxLenApprox: 0,
			ID:           "*",
			Values:       nil,
		}).Val()

		s.client.XDel(stream, id).Val()
		//rv := s.client.XDel(stream, id).Val()
		//if rv > 0 {
		//	s.client.XGroupCreate(stream, consumerGroup, "$")
		//
		//}
	}

	check_backlog := config.Get("broker.checkBackLog").Bool()
	var lastId = "0-0"

	for {
		//log.Println("in first for loop of recv")
		var myId string
		if check_backlog {
			myId = lastId
		} else {
			myId = ">"
		}

		args := redis.XReadGroupArgs{
			Group:    consumerGroup,
			Consumer: consumer,
			Streams:  []string{stream, myId},
			Count:    200,
			Block:    time.Duration(5) * time.Second,
		}

		log.GetLogger().Info("xreadgroup", zap.String("stream", stream), zap.String("consumer group", consumerGroup))
		result := s.client.XReadGroup(&args)

		ss := result.Val()
		log.GetLogger().Debug("stream_result", zap.String("stream", stream), zap.String("consumer_group", consumerGroup), zap.Any("read_group_result", ss))

		if ss == nil {
			log.GetLogger().Debug("Stream no message", zap.String("stream", stream))
			continue
		}

		// 说明已经处理完了之前 未确认 的消息，开始处理新消息，lastId改为 >
		if len(ss[0].Messages) == 0 {
			//log.Println("Stream " + stream + " read end")
			log.GetLogger().Info("Stream read end", zap.String("stream", stream))

			check_backlog = false
		}

		for _, item := range ss {
			messages := item.Messages

			for _, msg := range messages {
				id := msg.ID
				values := msg.Values

				if msg.Values == nil {
					log.GetLogger().Info("msg values is nil", zap.String("stream", stream), zap.String("consumer_group", consumerGroup))
					// Acknowledge and delete this kind of message
					xack_result := s.client.XAck(stream, consumerGroup, id)
					log.GetLogger().Info("xack_result", zap.Int64("value", xack_result.Val()))
					// 验证是否消费组的所有人都消费后的逻辑
					s.deleteLastMessages(stream, consumerGroup, id)

					lastId = id

					continue
				}

				if len(reflect.ValueOf(values).MapKeys()) > 0 {
					first_key := reflect.ValueOf(values).MapKeys()[0].String()
					//log.Println("result:",id, values)
					if values[first_key] != nil {
						//log.Println("values stream is not nil")
						log.GetLogger().Info("values stream is not nil")

						tempMessage := values[first_key].(string)
						//log.Println("message is " + tempMessage)
						p := publication{
							topic: s.stream,
							message: &broker.Message{
								Header: map[string]string{
									"stream": stream,
									"key":    first_key,
								},
								Body: []byte(tempMessage),
							},
						}

						// Handle error? Retry?
						if p.err = s.handle(&p); p.err != nil {
							//log.Fatal("handle message error")
							log.GetLogger().Info("handle message error")

						} else {
							log.GetLogger().Info("message consuming", zap.String("stream", stream), zap.String("consumer group", consumerGroup), zap.String("id", id))
							// 消息消费确认
							xack_result := s.client.XAck(stream, consumerGroup, id)
							log.GetLogger().Info("xack_result", zap.Int64("value", xack_result.Val()))
							// 验证是否消费组的所有人都消费后的逻辑
							s.deleteLastMessages(stream, consumerGroup, id)

							lastId = id
							//time.Sleep(time.Second * 0)
						}
					} else {
						//log.Println("message key is not found")
						log.GetLogger().Info("message key is not found")
					}

				}

			}

		}
	}

}

// Options returns the subscriber options.
func (s *subscriber) Options() broker.SubscribeOptions {
	return s.opts
}

// Topic returns the topic of the subscriber.
func (s *subscriber) Topic() string {
	return s.topic
}

func (s *subscriber) Stream() string {
	return s.stream
}

// Unsubscribe unsubscribes the subscriber and frees the connection.
func (s *subscriber) Unsubscribe() error {
	// TODO to refractor this later
	//s.client.Do("unsubscribe", s.Topic())
	//s.client.Pipeline().Close()
	return nil
}

// broker implementation for Redis.
type RedisBroker struct {
	addr               string
	client             *redis.Client
	multiStreamOptions []MultiStreamOption
	publishOptions     []PublishOption
	opts               broker.Options
	bopts              *brokerOptions
}

type MultiStreamOption struct {
	Stream        string `json:"stream"`
	ConsumerGroup string `json:"consumer_group"`
	Consumer      string `json:"consumer"`
}

func (b *RedisBroker) SetMultiStreamOptions() (values *[]MultiStreamOption) {
	var multiStreamOption []MultiStreamOption
	get_from_cfg := config.Get("broker.subscribe").String()

	if err := ejson.Unmarshal([]byte(get_from_cfg), &multiStreamOption); err != nil {
		//log.Println(err)
		log.GetLogger().Error("Unmarshal json error", zap.Error(err))
	}
	for _, m := range multiStreamOption {
		//log.Println(m.Stream, m.ConsumerGroup, m.Consumer)
		log.GetLogger().Info("Message Info", zap.String("stream", m.Stream), zap.String("ConsumerGroup", m.ConsumerGroup), zap.String("Consumer", m.Consumer))
		b.multiStreamOptions = append(b.multiStreamOptions, m)
	}

	return &b.multiStreamOptions
}

type PublishOption struct {
	Stream        string `json:"stream"`
	ConsumerGroup string `json:"consumer_group"`
}

func (b *RedisBroker) SetPublishOptions() (values *[]PublishOption) {
	var publishOptions []PublishOption

	get_from_cfg := config.Get("broker.publish").String()

	if err := ejson.Unmarshal([]byte(get_from_cfg), &publishOptions); err != nil {
		//log.Println(err)
		log.GetLogger().Info("SetPublishOptions unmarshal json", zap.Error(err))
	}

	for _, o := range publishOptions {
		//log.Println(o.Stream, o.ConsumerGroup)
		log.GetLogger().Info("SetPublishOptions Info", zap.String("stream", o.Stream), zap.String("consumer_group", o.ConsumerGroup))
		b.publishOptions = append(b.publishOptions, o)
	}
	return &b.publishOptions
}

// String returns the name of the broker implementation.
func (b *RedisBroker) String() string {
	return "redis"
}

// Options returns the options defined for the broker.
func (b *RedisBroker) Options() broker.Options {
	return b.opts
}

// Address returns the address the broker will use to create new connections.
// This will be set only after Connect is called.
func (b *RedisBroker) Address() string {
	return b.addr
}

// Init sets or overrides broker options.
func (b *RedisBroker) Init(opts ...broker.Option) error {
	//TODO to check for etcd and other registry configuration
	if b.client != nil {
		return errors.New("redis: cannot init while connected")
	}

	for _, o := range opts {
		o(&b.opts)
	}

	b.multiStreamOptions = []MultiStreamOption{}
	b.publishOptions = []PublishOption{}

	return nil
}

// Connect establishes a connection to Redis which provides the
// pub/sub implementation.
func (b *RedisBroker) Connect() error {
	if b.client != nil {
		return nil
	}

	var addr string

	if len(b.opts.Addrs) == 0 || b.opts.Addrs[0] == "" {
		addr = "127.0.0.1:6379"
	} else {
		addr = b.opts.Addrs[0]

		//if !strings.HasPrefix("redis://", addr) {
		//	addr = "redis://" + addr
		//}
	}

	b.addr = addr
	//log.Println(b.addr)
	addresses := strings.Split(b.addr, "@")

	var clientAddress string
	var clientPassword string

	if len(addresses) == 2 {
		clientAddress = addresses[1]
		clientPassword = addresses[0]
	} else if len(addresses) == 1 {
		clientAddress = addresses[0]
	}
	if len(addresses) == 0 {
		log.GetLogger().Error("Redis address error")
		return errors.New("Redis address error")
	} else {
		log.GetLogger().Debug("Client Info", zap.String("address", clientAddress), zap.String("password", clientPassword))
	}

	b.client = redis.NewClient(
		&redis.Options{
			Network:   "",
			Addr:      clientAddress,
			Dialer:    nil,
			OnConnect: nil,
			Password:  clientPassword,
			DB:        0,
			PoolSize:  10,
		})
	return nil
}

// Disconnect closes the connection pool.
func (b *RedisBroker) Disconnect() error {
	//err := b.client.Pipeline().Close()
	err := b.client.Conn().Close()
	b.client = nil
	b.addr = ""
	return err
}

func (b *RedisBroker) publish(stream string, values map[string]interface{}) string {
	//// TODO to check it
	//exists := b.client.Exists(stream)
	//if exists.Val() ==  0 {
	//	id := b.client.XAdd(&redis.XAddArgs{
	//		Stream:       stream,
	//		MaxLen:       0,
	//		MaxLenApprox: 0,
	//		ID:           "*",
	//		Values:       nil,
	//	}).Val()
	//
	//	b.client.XDel(stream, id)
	//}

	result := b.client.XAdd(&redis.XAddArgs{
		Stream:       stream,
		MaxLen:       5 * 100000,
		MaxLenApprox: 0,
		ID:           "*",
		Values:       values,
	})

	return result.Val()

}

// Publish publishes a message.
func (b *RedisBroker) Publish(topic string, msg *broker.Message, opts ...broker.PublishOption) error {

	var values map[string]interface{}
	values = make(map[string]interface{})
	values[topic] = string(msg.Body[:])

	b.publish(topic, values)

	err := b.client.Conn().Close()

	return err
}

// Subscribe returns a subscriber for the topic and handler.
func (b *RedisBroker) Subscribe(topic string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	log.GetLogger().Debug("redis subscribe")
	var options broker.SubscribeOptions
	for _, o := range opts {
		o(&options)
	}

	var consumer_group, consumer string

	for _, streamSettings := range b.multiStreamOptions {
		if topic == streamSettings.Stream {
			consumer_group = streamSettings.ConsumerGroup
			consumer = streamSettings.Consumer
		}
	}

	log.GetLogger().Info("Topic ConsumerGroup Info", zap.String("topic", topic), zap.String("consumer_group", consumer_group))

	s := subscriber{
		codec:         b.opts.Codec,
		client:        b.client,
		consumerGroup: consumer_group,
		stream:        topic,
		topic:         topic,
		handle:        handler,
		opts:          options,
	}

	log.GetLogger().Info("xgroup create", zap.String("stream", topic), zap.String("consumer group", consumer_group))
	s.client.XGroupCreate(topic, consumer_group, "$")

	go s.recv(topic, consumer_group, consumer)

	return &s, nil
}

// NewBroker returns a new broker implemented using the Redis pub/sub
// protocol. The connection address may be a fully qualified IANA address such
// as: redis://user:secret@localhost:6379/0?foo=bar&qux=baz
func NewBroker(opts ...broker.Option) broker.Broker {
	// Default options.
	bopts := &brokerOptions{
		maxIdle:        DefaultMaxIdle,
		maxActive:      DefaultMaxActive,
		idleTimeout:    DefaultIdleTimeout,
		connectTimeout: DefaultConnectTimeout,
		readTimeout:    DefaultReadTimeout,
		writeTimeout:   DefaultWriteTimeout,
	}

	// Initialize with empty broker options.
	options := broker.Options{
		Codec:   json.Marshaler{},
		Context: context.WithValue(context.Background(), optionsKey, bopts),
	}

	for _, o := range opts {
		o(&options)
	}

	return &RedisBroker{
		opts:  options,
		bopts: bopts,
	}
}
