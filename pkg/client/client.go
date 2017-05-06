package client

import (
	"github.com/lavrs/dms/pkg/context"
	"github.com/lavrs/dms/pkg/logger"
	"github.com/urfave/cli"
	"os"
)

func Run() {
	const (
		debug = "d"
		addr  = "a"
	)

	app := cli.NewApp()

	app.Name = "dms"
	app.Usage = "Docker monitoring service"
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  debug + ", debug",
			Usage: "set debug mode",
		},
		cli.StringFlag{
			Name:  addr + ", addr",
			Value: "http://localhost:4222/api",
			Usage: "set daemon address",
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
		context.Get().Debug = c.Bool(debug)
		// set daemon address
		context.Get().Address = c.String(addr)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Panic(err)
	}
}
