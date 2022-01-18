package mq

import "github.com/micro/go-micro/v2/broker"

var RedisBroker broker.Broker

func Init(broker broker.Broker)  {
	RedisBroker = broker
}