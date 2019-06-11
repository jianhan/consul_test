package main

import (
	"context"
	"fmt"
	"time"

	proto "github.com/jianhan/consul_test/greeting/proto"
	micro "github.com/micro/go-micro"
	mprometheus "github.com/micro/go-plugins/wrapper/monitoring/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error {
	rsp.Greeting = "Hello " + req.Name
	return nil
}
func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
)

func main() {
	recordMetrics()
	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name("greeter"),
		micro.Version("v1"),
		micro.WrapHandler(mprometheus.NewHandlerWrapper()),
	)

	// Init will parse the command line flags.
	service.Init()
	t := "test"

	fmt.Println(t)
	// Register handler
	proto.RegisterGreeterHandler(service.Server(), new(Greeter))

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
