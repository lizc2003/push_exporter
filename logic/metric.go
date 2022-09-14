package logic

import (
	"fmt"
	"strings"
)

type MetricValue struct {
	Metric     string  `json:"metric"`
	Instance   string  `json:"endpoint"`
	Job        string  `json:"job"`
	Tags       string  `json:"tags"`
	Value      float64 `json:"value"`
	Type       string  `json:"type"`
	Step       int64   `json:"step"`
	ExpireTime int64   `json:"-"`
}

func (this *MetricValue) Key() string {
	return strings.Join([]string{this.Metric, this.Instance, this.Job, this.Tags}, "|")
}

func (this *MetricValue) GaugeKey(labelNames []string) string {
	return this.Metric + "|" + strings.Join(labelNames, "|")
}

func (this *MetricValue) GetLabels() ([]string, []string) {
	labelNames := []string{InstanceName, JobName}
	labelVals := []string{this.Instance, this.Job}
	if this.Tags != "" {
		ss := strings.Split(this.Tags, ",")
		for _, s := range ss {
			ll := strings.Split(s, "=")
			if len(ll) == 2 {
				labelNames = append(labelNames, ll[0])
				labelVals = append(labelVals, ll[1])
			}
		}
	}
	return labelNames, labelVals
}

func (this *MetricValue) IsValid(now int64) bool {
	return this.ExpireTime >= now
}

func (this *MetricValue) String() string {
	return fmt.Sprintf(
		"<Metric:%s, Instance:%s, Job:%s, Tags:%s, Value:%f, Type:%s, Step:%d>",
		this.Metric,
		this.Instance,
		this.Job,
		this.Tags,
		this.Value,
		this.Type,
		this.Step,
	)
}
