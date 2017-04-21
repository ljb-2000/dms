package daemon

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	m "github.com/lavrs/docker-monitoring-service/pkg/metrics"
	"net/http"
	"time"
)

var metrics = m.NewMetrics()

func Run(port string, ucltime, uctime int) error {
	router := httprouter.New()

	metrics.SetUCLTime(time.Second * time.Duration(ucltime))
	metrics.SetUCTime(time.Second * time.Duration(uctime))
	go metrics.Collect()

	router.GET("/metrics/:id", getMetrics)

	return http.ListenAndServe(":"+port, router)
}

func getMetrics(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.WriteHeader(200)

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")

	metricsJSON, err := json.Marshal(metrics.Get(p.ByName("id")))
	if err != nil {
		panic(err)
	}

	w.Write(metricsJSON)
}
