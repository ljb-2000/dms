package daemon

import (
	"github.com/lavrs/dms/pkg/context"
	m "github.com/lavrs/dms/pkg/daemon/metrics"
	"github.com/lavrs/dms/pkg/daemon/router"
	"github.com/lavrs/dms/pkg/logger"
	"github.com/urfave/cli"
	"net/http"
	"os"
	"time"
)

// Run start daemon
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

		// set debug mode if use flag "d"
		context.Get().Debug = c.Bool("d")
		// set daemon address
		context.Get().Address = ":" + c.String("p")

		// set update intervals
		m.Get().SetUCListInterval(time.Duration(c.Int("ucli")) * time.Second)
		m.Get().SetUCMetricsInterval(time.Duration(c.Int("uci")) * time.Second)

		// start collect metrics
		go m.Get().Collect()

		// listen and serve
		fsrv := &http.Server{
			Handler: router.App(),
			Addr:    context.Get().Address,
		}
		return fsrv.ListenAndServe()
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Panic(err)
	}
}
