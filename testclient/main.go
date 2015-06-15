package main

import (
	"time"

	"golang.org/x/net/context"
	"gopkg.in/spacemonkeygo/monitor.v1/trace"
	"gopkg.in/spacemonkeygo/monitor.v1/trace/gen-go/zipkin"
)

func main() {
	trace.Configure(1, true, &zipkin.Endpoint{
		Ipv4:        0,
		Port:        8080,
		ServiceName: "go-zipkin-testclient",
	})

	if c, e := trace.NewScribeCollector("0.0.0.0:9410"); e != nil {
		panic(e)
	} else {
		trace.RegisterTraceCollector(c)
	}

	for {
		ctx := context.Background()
		if e := firstCall(ctx); e != nil {
			panic(e)
		}
		time.Sleep(1 * time.Second)
	}
}

func firstCall(c context.Context) (e error) {
	defer trace.Trace(&c)(&e)
	time.Sleep(3000 * time.Microsecond)
	return secondCall(c)
}

func secondCall(c context.Context) (e error) {
	defer trace.Trace(&c)(&e)
	time.Sleep(500 * time.Microsecond)
	return thirdCall(c)
}

func thirdCall(c context.Context) (e error) {
	defer trace.Trace(&c)(&e)
	time.Sleep(100 * time.Microsecond)
	return nil
}
