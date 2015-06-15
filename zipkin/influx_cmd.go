package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/mattkanwisher/distributedtrace"
)

var InfluxCmd = cli.Command{
	Name:    "influx",
	Aliases: []string{"inf", "i"},
	Usage:   "collect all spans with the same traceId and combine them into influx rows.",
	Before:  InfluxBefore,
	Action:  InfluxAction,
}

// TODO: Allow --host --user --pass style of flags
func InfluxBefore(c *cli.Context) (e error) {
	if len(c.Args()) == 0 {
		return fmt.Errorf("influx address not specified.")
	}

	var u *url.URL
	if u, e = url.Parse(c.Args()[0]); e != nil {
		return e
	} else if u == nil {
		return fmt.Errorf("invalid influx url.")
	}

	if len(u.Path) < 2 {
		return fmt.Errorf("influx url must contains a database name as the path segment.")
	} else if u.User == nil {
		return fmt.Errorf("influx url requires authentication info.")
	} else if u.User.Username() == "" {
		return fmt.Errorf("influx url must contains a username.")
	} else if pass, ok := u.User.Password(); !ok || pass == "" {
		return fmt.Errorf("influx url must contains a password.")
	}

	return nil
}

func InfluxAction(c *cli.Context) {
	config := buildServerConfig(c)

	u, _ := url.Parse(c.Args()[0]) // TODO: Avoid re-parse.
	user := u.User.Username()
	pass, _ := u.User.Password()

	path, series := u.Path[1:], "zipkin"
	splitIdx := strings.Index(path, "/")
	if splitIdx > -1 {
		path, series = path[:splitIdx], path[splitIdx+1:]
	}

	// TODO: Allow configuring series name.
	output, e := zipkin.NewInfluxOutput(config, u.Host, path, user, pass, series)
	if e != nil {
		panic(e)
	}

	startServer(c, config, output)
}
