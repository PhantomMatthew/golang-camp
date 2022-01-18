package main

import (
	"github.com/micro/go-micro/v2/client"
	grpc_client "github.com/micro/go-micro/v2/client/grpc"
	"newline.com/newline/src/common/utils"
	"newline.com/newline/src/srv/customer/jobs"
	"newline.com/newline/src/srv/customer/subscriber"
	//"github.com/micro/go-micro/v2/server"
	//"github.com/micro/go-micro/v2/server/grpc"
	"os"
	"time"

	"go.uber.org/zap"
	"newline.com/newline/src/common/config"
	"newline.com/newline/src/common/log"

	"github.com/micro/go-micro/v2/broker"
	"newline.com/newline/src/srv/customer/handler"
	customer "newline.com/newline/src/srv/customer/proto/customer"

	"github.com/micro/go-micro/v2"

	"newline.com/newline/src/srv/customer/components"
	tracer "newline.com/newline/src/srv/plugins/tracer/jaeger"

	//"github.com/afex/hystrix-go/hystrix"
	//breaker "github.com/micro/go-plugins/wrapper/breaker/hystrix/v2"
	limiter "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"
	opentracingPlugin "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"

	"newline.com/newline/src/go-plugins/broker/redis2"

	"newline.com/newline/src/common/mq"
)

func init() {
	components.Init()
}

func main() {
	// Broker iniialization
	url := config.Get("broker.host").String()
	if url == "" {
		log.GetLogger().Error("no redis broker url")
	}

	b := redis2.NewBroker(broker.Addrs(url))
	b.Init()
	//b.(*redis2.RedisBroker).SetConsumerGroup(config.Get("persistence.redis.customerConsumerGroup").String())
	streamOptions := b.(*redis2.RedisBroker).SetMultiStreamOptions()
	log.GetLogger().Info("Redis Broker SetMultiStreamOptions", zap.Any("multiStreamOptions", streamOptions))

	//zaplog.Info("sub stream options", zap.String("streamOptions", streamOptions))

	if err := b.Connect(); err != nil {
		log.GetLogger().Error("Broker connect error", zap.Error(err))
	}
	defer b.Disconnect()

	//type Broker struct {
	//	Stream        string `json:"stream"`
	//	ConsumerGroup string `json:"consumer_group"`
	//	Consumer      string `json:"consumer"`
	//}
	// 获取消息队列配置
	brokers := config.Get("broker.subscribe").Array()
	for _, v := range brokers {
		if v.Map()["enable"].Bool() {
			b.Subscribe(v.Map()["stream"].Str, subscriber.NewReceiveMsgHandler().ReceiveMsgHandler)
		}
	}

	//b.Subscribe("bbcl_youzan", subscriber.NewReceiveMsgHandler().ReceiveMsgHandler)
	////b.Subscribe("weimob", subscriber.NewReceiveMsgHandler().ReceiveMsgHandler)
	//b.Subscribe("bbcl_weapp", subscriber.NewReceiveMsgHandler().ReceiveMsgHandler)
	//b.Subscribe("bbcl_wechat", subscriber.NewReceiveMsgHandler().ReceiveMsgHandler)

	mq.Init(b)

	// Tracer initialization
	traceAddr := os.Getenv("MICRO_BOOK_TRACER_ADDR")
	t, io, err := tracer.NewTracer(components.GetServiceName(), traceAddr)
	if err != nil {
		log.GetLogger().Error("Tracer error", zap.Error(err))
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	const QPS = 2000

	serviceName := config.Get("name.customer").String()
	if serviceName != "" {
		log.GetLogger().Debug("get service name successfully")
	} else {
		serviceName = "go.micro.service.customerdefault"
	}

	// New Service
	//newServer := server.NewServer(grpc.MaxMsgSize(100 * 1024 * 1024))
	service := micro.NewService(
		micro.Name(serviceName),
		//micro.Address(":9999"),
		micro.RegisterTTL(time.Second*15),
		micro.RegisterInterval(time.Second*10),
		//micro.Client(client.NewClient(client.RequestTimeout(15 * time.Second)),
		micro.Version("latest"),
		//micro.WrapClient(breaker.NewClientWrapper()),
		micro.WrapHandler(
			limiter.NewHandlerWrapper(QPS),
			opentracingPlugin.NewHandlerWrapper(opentracing.GlobalTracer()),
		),
		micro.Broker(b),
		micro.Registry(utils.GetRegistry()),
	)

	// TODO: to refactor in configuration
	// Add hystrix settings as configuration
	//hystrix.DefaultMaxConcurrent = 5000
	//hystrix.DefaultTimeout = 200

	// Initialise service
	service.Init()
	//customerInfo := models.CustBasicInfo{UnionId: "123123"}
	//components.MainDB.Create(&customerInfo)
	// Register Handler
	//customer.RegisterCustomerHandler(service.Server(), new(handler.Customer))
	s := client.NewClient(grpc_client.MaxRecvMsgSize(10 * 10 * 1024))
	_ = customer.NewCustomerService("go.micro.service.customerdefault", s)
	_ = customer.RegisterCustomerHandler(service.Server(), handler.NewCustHandler())

	//dao := dao.CustEmInfoDao{}
	//dao.SearchCustEmInfoList(model.SearchCustEmParam{}, models.PaginationParam{
	//	PageIndex: 1,
	//	PageSize:  10,
	//})
	if config.Get("app.is_prod").Bool() {
		jobs.StartCustomerJob()
	}
	// Register Struct as Subscriber
	//micro.RegisterSubscriber(config.Get("persistence.redis.customerStream").String(), service.Server(), new(subscriber.Customer).Handle)
	//Copy from NoSQLBooster for MongoDB free edition. This message does not appear if you are using a registered version.

	//bll.NewCustBll().UpsertCustomer(nil, result, "em")
	//bll.NewCustBll().SyncTag(nil, []request.Tag{
	//	{Id: "1", Name: "test"},
	//	{Id: "2", Name: "test2"},
	//}, 1, "yz")

	//new(publisher.CustomerPublisher).Publish()
	log.GetLogger().Info("CustomerServiceStart")
	// Run service
	if err := service.Run(); err != nil {
		log.GetLogger().Error("Service Run error", zap.Error(err))
	}

}
