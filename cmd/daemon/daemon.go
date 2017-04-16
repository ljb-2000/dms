package main

import (
	"github.com/lavrs/docker-monitoring-service/pkg/echo"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "p, port",
			Value: "8080",
			Usage: "daemon port",
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.NArg() > 0 {
			err := cli.ShowAppHelp(c)
			if err != nil {
				panic(err)
			}
		}

		echo.Echo(c.String("port"))

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
