package main

import (
	"log"

	tracer "github.com/micro-in-cn/tutorials/microservice-in-micro/part8/plugins/tracer/jaeger"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part8/plugins/tracer/opentracing/stdhttp"
	"github.com/micro/micro/v2/cmd"
	"github.com/micro/micro/v2/plugin"
	"github.com/opentracing/opentracing-go"
)

func init() {
	//plugin.Register(cors.NewPlugin())

	plugin.Register(plugin.NewPlugin(
		plugin.WithName("tracer"),
		plugin.WithHandler(
			stdhttp.TracerWrapper,
		),
	))
}

const name = "gateway"

func main() {
	stdhttp.SetSamplingFrequency(50)
	t, io, err := tracer.NewTracer(name, "")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	cmd.Init()
}
