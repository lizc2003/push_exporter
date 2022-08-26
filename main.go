package main

import (
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/lizc2003/push_exporter/logic"
	"net/http"
	"os"
)

func main() {
	var landingPage = []byte(`<html>
<head><title>Push exporter</title></head>
<body>
<h1>Push exporter</h1>
<p><a href='metrics'>Metrics</a></p>
</body>
</html>
`)

	logger := promlog.New(&promlog.Config{})
	metricMgr := logic.NewMetricMgr()

	exporterHandler := logic.NewExporterHandler(metricMgr)
	http.Handle("/metrics", promhttp.InstrumentMetricHandler(prometheus.DefaultRegisterer, exporterHandler))

	http.Handle("/push", logic.NewPushHandler(metricMgr))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(landingPage)
	})

	srv := &http.Server{Addr: ":1999"}
	if err := web.ListenAndServe(srv, "", logger); err != nil {
		level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
		os.Exit(1)
	}
}
