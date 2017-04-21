package main

import (
	"github.com/lavrs/docker-monitoring-service/pkg/context"
	"github.com/lavrs/docker-monitoring-service/pkg/daemon"
	"github.com/lavrs/docker-monitoring-service/pkg/logger"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()

	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "p, port",
			Value: "4222",
			Usage: "set daemon port",
		},
		cli.IntFlag{
			Name:  "uct, upd-container-time",
			Value: 1,
			Usage: "set update container interval",
		},
		cli.IntFlag{
			Name:  "uclt, upd-container-list-time",
			Value: 3,
			Usage: "set update container list interval",
		},
		cli.BoolFlag{
			Name:  "d, debug",
			Usage: "set debug mode",
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.NArg() > 0 {
			err := cli.ShowAppHelp(c)
			if err != nil {
				return err
			}
			return nil
		}

		ctx := context.Get()
		ctx.Debug = c.Bool("d")

		err := daemon.Run(c.String("p"), c.Int("uclt"), c.Int("uct"))
		if err != nil {
			return err
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Panic(err)
	}
}
