package main

import (
	"flag"
	"fmt"
	"github.com/go-kit/log/level"
	"github.com/lizc2003/push_exporter/logic"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/exporter-toolkit/web"
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
	//http.Handle("/metrics", promhttp.InstrumentMetricHandler(prometheus.DefaultRegisterer, exporterHandler))
	http.Handle("/metrics", exporterHandler)

	http.Handle("/push", logic.NewPushHandler(metricMgr))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(landingPage)
	})

	webConfigFile := ""
	toolkitFlags := web.FlagConfig{
		WebListenAddresses: &[]string{fmt.Sprintf(":%d", getListenPort())},
		WebConfigFile:      &webConfigFile,
	}
	srv := &http.Server{}
	if err := web.ListenAndServe(srv, &toolkitFlags, logger); err != nil {
		level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
		os.Exit(1)
	}
}

func getListenPort() int {
	port := flag.Int("p", 1999, "port")
	flag.Parse()
	return *port
}
