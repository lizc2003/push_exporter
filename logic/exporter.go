package logic

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

type Exporter struct {
	ctx context.Context
	mgr *MetricMgr
}

func NewExporter(ctx context.Context, mgr *MetricMgr) *Exporter {
	return &Exporter{
		ctx: ctx,
		mgr: mgr,
	}
}

func (this *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- this.mgr.Builtin.Up.Desc()
}

func (this *Exporter) Collect(ch chan<- prometheus.Metric) {
	mgr := this.mgr
	ch <- mgr.Builtin.Up

	now := time.Now().Unix()
	metrics := mgr.GetMetrics()
	for _, m := range metrics {
		if m.IsValid(now) {
			g, err := mgr.CreateGauge(m)
			if err == nil {
				ch <- g
			}
		}
	}
}

/////////////////////////////////////////////////////////////
func NewExporterHandler(mgr *MetricMgr) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		registry := prometheus.NewRegistry()
		registry.MustRegister(NewExporter(ctx, mgr))
		gatherers := prometheus.Gatherers{
			//prometheus.DefaultGatherer,
			registry,
		}
		h := promhttp.HandlerFor(gatherers, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
	}
}
