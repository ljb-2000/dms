package daemon

import (
	"github.com/lavrs/docker-monitoring-service/pkg/daemon/router"
	m "github.com/lavrs/docker-monitoring-service/pkg/metrics"
	"net/http"
	"time"
)

// Run start daemon
func Run(port string, ucli, uci int) error {
	// set update intervals
	m.Get().SetUCListInterval(time.Duration(ucli) * time.Second)
	m.Get().SetUCMetricsInterval(time.Duration(uci) * time.Second)

	// start collect metrics
	go m.Get().Collect()

	fsrv := &http.Server{
		Handler: router.App(),
		Addr:    ":" + port,
	}
	return fsrv.ListenAndServe()
}
