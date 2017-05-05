package client

import (
	"github.com/lavrs/dms/pkg/context"
	"github.com/lavrs/dms/pkg/logger"
	"github.com/urfave/cli"
	"os"
)

func Run() {
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

		// set debug mode if use flag "d"
		context.Get().Debug = c.Bool("d")

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Panic(err)
	}
}
