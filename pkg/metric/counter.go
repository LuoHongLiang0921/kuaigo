// @Description 计数类型 ，并且只能增长和重置。例如：一个网站的总访问量，机器的运行时长

package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

// CounterVecOpts ...
type CounterVecOpts struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
}

// Build ...
func (opts CounterVecOpts) Build() *counterVec {
	vec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: opts.Namespace,
			Subsystem: opts.Subsystem,
			Name:      opts.Name,
			Help:      opts.Help,
		}, opts.Labels)
	prometheus.MustRegister(vec)
	return &counterVec{
		CounterVec: vec,
	}
}

// NewCounterVec ...
func NewCounterVec(name string, labels []string) *counterVec {
	return CounterVecOpts{
		Namespace: DefaultNamespace,
		Name:      name,
		Help:      name,
		Labels:    labels,
	}.Build()
}

type counterVec struct {
	*prometheus.CounterVec
}

// Inc ...
func (counter *counterVec) Inc(labels ...string) {
	counter.WithLabelValues(labels...).Inc()
}

// Add ...
func (counter *counterVec) Add(v float64, labels ...string) {
	counter.WithLabelValues(labels...).Add(v)
}

//variableCounterVec 标签可变
type variableCounterVec struct {
	m *metric
}

// CounterVec  namespace为 tabby，subsystem为 go
var CounterVec = &variableCounterVec{
	m: Option{
		Namespace: "tabby",
		Subsystem: "go",
		mt:        cv,
	}.Build(),
}

//NewHistogramVec ...
func NewVariableCounterVec(r Option) *variableCounterVec {
	return &variableCounterVec{m: r.Build()}
}

// Inc ...
func (counter *variableCounterVec) Inc(name string, kv interface{}) {
	lbNames, lbValues := genLabels(kv)
	//lbNames
	v := counter.m.loadOrStore(name, lbNames)
	if v != nil {
		vv := v.(*prometheus.CounterVec)
		vv.WithLabelValues(lbValues...).Inc()
	}
}

// Add ...
func (counter *variableCounterVec) Add(name string, kv interface{}, v float64) {
	lbNames, lbValues := genLabels(kv)
	//lbNames
	v1 := counter.m.loadOrStore(name, lbNames)
	if v1 != nil {
		vv := v1.(*prometheus.CounterVec)
		vv.WithLabelValues(lbValues...).Add(v)
	}
}

// WithHelp 设置help 信息
func (counter *variableCounterVec) WithHelp(help string) *variableCounterVec {
	counter.m.opt.Help = help
	return counter
}

// WithName 设置name
func (counter *variableCounterVec) WithName(name string) *variableCounterVec {
	counter.m.opt.Name = name
	return counter
}
