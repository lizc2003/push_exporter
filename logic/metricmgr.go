package logic

import (
	"github.com/prometheus/client_golang/prometheus"
	"os"
	"sync"
)

type MetricMgr struct {
	mtx        sync.Mutex
	mtxGauge   sync.Mutex
	mtxCounter sync.Mutex
	metrics    map[string]*MetricValue
	gauges     map[string]*prometheus.GaugeVec
	counters   map[string]*CounterVec
	Builtin    *BuiltinMetrics
}

func NewMetricMgr() *MetricMgr {
	hostName, _ := os.Hostname()
	return &MetricMgr{
		metrics:  make(map[string]*MetricValue),
		gauges:   make(map[string]*prometheus.GaugeVec),
		counters: make(map[string]*CounterVec),
		Builtin:  NewBuiltinMetrics(hostName),
	}
}

func (this *MetricMgr) AddMetric(m *MetricValue, now int64) {
	if m.Metric == "" {
		return
	}

	if m.Step == 0 {
		m.Step = 5 // default 5s
	}
	m.ExpireTime = now + m.Step + 2
	switch m.Type {
	case TypeCounter:
	default:
		m.Type = TypeGauge
	}

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
	var gVec *prometheus.GaugeVec
	var err error
	var ok bool

	labelNames, labelVals := m.GetLabels()
	key := m.MetricKey(labelNames)

	this.mtxGauge.Lock()
	if gVec, ok = this.gauges[key]; !ok {
		gVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: m.Metric,
			Help: "",
		}, labelNames)
		this.gauges[key] = gVec
	}
	ret, err = gVec.GetMetricWithLabelValues(labelVals...)
	if err == nil {
		ret.Set(m.Value)
	}
	this.mtxGauge.Unlock()

	return ret, err
}

func (this *MetricMgr) CreateCounter(m *MetricValue) (Counter, error) {
	var ret Counter
	var aVec *CounterVec
	var err error
	var ok bool

	labelNames, labelVals := m.GetLabels()
	key := m.MetricKey(labelNames)

	this.mtxCounter.Lock()
	if aVec, ok = this.counters[key]; !ok {
		aVec = NewCounterVec(m.Metric, "", labelNames)
		this.counters[key] = aVec
	}
	ret, err = aVec.GetMetricWithLabelValues(labelVals...)
	if err == nil {
		ret.Set(m.Value)
	}
	this.mtxCounter.Unlock()

	return ret, err
}
