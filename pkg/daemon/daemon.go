package daemon

import (
	"github.com/lavrs/docker-monitoring-service/pkg/daemon/router"
	m "github.com/lavrs/docker-monitoring-service/pkg/metrics"
	"net/http"
	"time"
)

// Run daemon
func Run(port string, ucli, uci int) error {
	metrics := m.Get()
	metrics.SetUCListInterval(time.Duration(ucli) * time.Second)
	metrics.SetUCMetricsInterval(time.Duration(uci) * time.Second)
	go metrics.Collect()

	fsrv := &http.Server{
		Handler: router.App(),
		Addr:    ":" + port,
	}
	return fsrv.ListenAndServe()
}
