package logic

import "github.com/prometheus/client_golang/prometheus"

type BuiltinMetrics struct {
	Up prometheus.Gauge
}

func NewBuiltinMetrics(hostName string) *BuiltinMetrics {
	ret := &BuiltinMetrics{
		Up: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   Namespace,
			Subsystem:   "exporter",
			Name:        "up",
			Help:        "Whether the exporter is up.",
			ConstLabels: prometheus.Labels{InstanceName: hostName},
		}),
	}
	ret.Up.Set(1)

	return ret
}
