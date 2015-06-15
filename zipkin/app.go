package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var log = logrus.StandardLogger()

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "zipkin"
	app.Usage = "proxy for piping zipkin traces into something else."
	app.Commands = []cli.Command{
		InfluxCmd,
	}
	app.Flags = []cli.Flag{
		listenFlag,
	}

	return app
}
