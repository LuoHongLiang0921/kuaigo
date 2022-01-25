// @Description 直方图 一个histogram会生成三个指标，分别是_count，_sum，_bucket。

package metric

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// HistogramVecOpts ...
type HistogramVecOpts struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
	Buckets   []float64
}

type histogramVec struct {
	*prometheus.HistogramVec
}

// Build ...
func (opts HistogramVecOpts) Build() *histogramVec {
	vec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: opts.Namespace,
			Subsystem: opts.Subsystem,
			Name:      opts.Name,
			Help:      opts.Help,
			Buckets:   opts.Buckets,
		}, opts.Labels)
	prometheus.MustRegister(vec)
	return &histogramVec{
		HistogramVec: vec,
	}
}

// Observe ...
func (histogram *histogramVec) Observe(v float64, labels ...string) {
	histogram.WithLabelValues(labels...).Observe(v)
}

//variableHistogramVec 标签可变
type variableHistogramVec struct {
	m *metric
}

// HistogramVec  namespace为 tabby，subsystem为 go
var HistogramVec = &variableHistogramVec{
	m: Option{
		Namespace: "tabby",
		Subsystem: "go",
		mt:        hv,
	}.Build(),
}

//NewVariableHistogramVec ...
func NewVariableHistogramVec(r Option) *variableHistogramVec {
	return &variableHistogramVec{m: r.Build()}
}

// Timing ...
// kv 可以是[]string, map[string]string
func (histogram *variableHistogramVec) Timing(name string, kv interface{}, startAt time.Time) {
	lbNames, lbValues := genLabels(kv)
	//lbNames
	v := histogram.m.loadOrStore(name, lbNames)
	if v != nil {
		vv := v.(*prometheus.HistogramVec)
		vv.WithLabelValues(lbValues...).Observe(time.Since(startAt).Seconds())
	}
}

// WithHelp 设置help 信息
func (histogram *variableHistogramVec) WithHelp(help string) *variableHistogramVec {
	histogram.m.opt.Help = help
	return histogram
}

// WithName 设置name
func (histogram *variableHistogramVec) WithName(name string) *variableHistogramVec {
	histogram.m.opt.Name = name
	return histogram
}
