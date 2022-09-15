package logic

import (
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"math"
	"sync/atomic"
)

type Counter struct {
	valBits    uint64
	desc       *prometheus.Desc
	labelPairs []*dto.LabelPair
}

func (this Counter) Desc() *prometheus.Desc {
	return this.desc
}

func (this Counter) Write(out *dto.Metric) error {
	val := math.Float64frombits(atomic.LoadUint64(&this.valBits))
	out.Label = this.labelPairs
	out.Counter = &dto.Counter{Value: &val}
	return nil
}

func (this *Counter) Set(val float64) {
	atomic.StoreUint64(&this.valBits, math.Float64bits(val))
}

type CounterVec struct {
	*prometheus.MetricVec
}

func NewCounterVec(name, help string, labelNames []string) *CounterVec {
	desc := prometheus.NewDesc(name, help, labelNames, nil)
	return &CounterVec{
		MetricVec: prometheus.NewMetricVec(desc, func(lvs ...string) prometheus.Metric {
			if len(lvs) != len(labelNames) {
				panic("inconsistent label cardinality")
			}
			return Counter{desc: desc, labelPairs: prometheus.MakeLabelPairs(desc, lvs)}
		}),
	}
}

func (v *CounterVec) GetMetricWithLabelValues(lvs ...string) (Counter, error) {
	metric, err := v.MetricVec.GetMetricWithLabelValues(lvs...)
	return metric.(Counter), err
}
