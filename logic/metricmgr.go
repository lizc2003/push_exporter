package logic

import (
	"github.com/prometheus/client_golang/prometheus"
	"os"
	"sync"
)

type MetricMgr struct {
	mtx      sync.Mutex
	mtxGauge sync.Mutex
	metrics  map[string]*MetricValue
	gauges   map[string]*prometheus.GaugeVec
	Builtin  *BuiltinMetrics
}

func NewMetricMgr() *MetricMgr {
	hostName, _ := os.Hostname()
	return &MetricMgr{
		metrics: make(map[string]*MetricValue),
		gauges:  make(map[string]*prometheus.GaugeVec),
		Builtin: NewBuiltinMetrics(hostName),
	}
}

func (this *MetricMgr) AddMetric(m *MetricValue, now int64) {
	m.ExpireTime = now + m.Step
	key := m.Key()

	this.mtx.Lock()
	this.metrics[key] = m
	this.mtx.Unlock()
}

func (this *MetricMgr) GetMetrics() map[string]*MetricValue {
	var ret map[string]*MetricValue

	this.mtx.Lock()
	if len(this.metrics) > 0 {
		ret = this.metrics
		this.metrics = make(map[string]*MetricValue)
	}
	this.mtx.Unlock()
	return ret
}

func (this *MetricMgr) CreateGauge(m *MetricValue) (prometheus.Gauge, error) {
	var ret prometheus.Gauge
	var err error
	var gvec *prometheus.GaugeVec
	var ok bool

	labelNames, labelVals := m.GetLabels()
	key := m.GaugeKey(labelNames)

	this.mtxGauge.Lock()
	if gvec, ok = this.gauges[key]; !ok {
		gvec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: m.Metric,
			Help: "",
		}, labelNames)
		this.gauges[key] = gvec
	}
	ret, err = gvec.GetMetricWithLabelValues(labelVals...)
	if err == nil {
		ret.Set(m.Value)
	}
	this.mtxGauge.Unlock()

	return ret, err
}
