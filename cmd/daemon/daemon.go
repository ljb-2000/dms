package main

import (
	"github.com/lavrs/docker-monitoring-service/pkg/daemon"
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
	}

	app.Action = func(c *cli.Context) error {
		if c.NArg() > 0 {
			err := cli.ShowAppHelp(c)
			if err != nil {
				return err
			}
			return nil
		}

		err := daemon.Run(c.String("p"))
		if err != nil {
			panic(err)
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
