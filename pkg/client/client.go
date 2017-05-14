package client

import (
	"github.com/lavrs/dms/pkg/client/cmd"
	"github.com/lavrs/dms/pkg/context"
	"github.com/lavrs/dms/pkg/logger"
	"github.com/urfave/cli"
	"os"
)

func Run() {
	const (
		// flags
		debug = "d"
		addr  = "a"

		// commands
		stopped  = "stopped"
		launched = "launched"
		metrics  = "metrics"
		logs     = "logs"
	)

	app := cli.NewApp()

	app.Name = "dms"
	app.Usage = "Docker monitoring service"
	app.Version = "0.1.0"

	app.Before = func(c *cli.Context) error {
		// set debug mode if use flag "d"
		context.Get().Debug = c.Bool(debug)
		// set daemon address
		context.Get().Address = c.String(addr)

		return nil
	}

	app.Commands = []cli.Command{
		// command for view stopped containers
		{
			Name:  stopped,
			Usage: "view stopped containers",
			Action: func(c *cli.Context) error {
				cmd.StoppedContainersCmd()
				return nil
			},
		},
		// command for view launched containers
		{
			Name:  launched,
			Usage: "view launched containers",
			Action: func(c *cli.Context) error {
				cmd.LaunchedContainersCmd()
				return nil
			},
		},
		// command for view container logs
		{
			Name:  logs,
			Usage: "view container logs",
			Action: func(c *cli.Context) error {
				cmd.ContainersLogsCmd(c.Args().First())
				return nil
			},
		},
		// command for view container(s) metrics
		{
			Name:  metrics,
			Usage: "view container(s) metrics",
			Action: func(c *cli.Context) error {
				cmd.ContainersMetricsCmd(c.Args())
				return nil
			},
		},
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  debug + ", debug",
			Usage: "set debug mode",
		},
		cli.StringFlag{
			Name:  addr + ", addr",
			Value: "http://localhost:4222",
			Usage: "set daemon address",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Panic(err)
	}
}
