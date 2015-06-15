package main

import (
	"github.com/codegangsta/cli"
)

type Cmd interface {
	Info() *cli.Command
	Action()
}
