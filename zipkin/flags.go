package main

import (
	"time"

	"github.com/codegangsta/cli"
)

var listenFlag = cli.StringFlag{
	Name:  "l",
	Value: "0.0.0.0:9410",
	Usage: "specify address to listen for zipkin spans / scribe packets.",
}

var timeoutFlag = cli.DurationFlag{
	Name:  "timeout",
	Value: 5 * time.Second,
	Usage: "the amount of time to wait for spans per TraceId before flushing.",
}
