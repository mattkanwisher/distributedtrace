package main

import (
	stdlog "log"

	"github.com/codegangsta/cli"
	"github.com/mattkanwisher/distributedtrace"
)

func buildServerConfig(c *cli.Context) *zipkin.Config {
	return &zipkin.Config{
		ListenAddress: c.GlobalString(listenFlag.Name),
		TraceTimeout:  c.GlobalDuration(timeoutFlag.Name),
		Logger:        stdlog.New(log.Writer(), "", 0),
	}
}

func startServer(c *cli.Context, config *zipkin.Config, output zipkin.Output) {
	log.Infof("starting zipkin server at: %s", config.ListenAddress)
	server, e := zipkin.NewServer(config, output)
	if e != nil {
		panic(e)
	}

	e = server.Start()
	if e != nil {
		panic(e)
	}

	log.Infof("shutdown.")
}
