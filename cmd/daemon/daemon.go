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

	app.Name = "dms"
	app.Usage = "Docker monitoring service"
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "p, port",
			Value: "4222",
			Usage: "set daemon port",
		},
		cli.IntFlag{
			Name:  "ucli, upd-container-list-interval",
			Value: 3,
			Usage: "set update container list interval",
		},
		cli.IntFlag{
			Name:  "uci, upd-container-interval",
			Value: 1,
			Usage: "set update container metrics interval",
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

        // set debug mode
		context.Get().Debug = c.Bool("d")

		return daemon.Run(c.String("p"), c.Int("ucli"), c.Int("uci"))
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Panic(err)
	}
}
