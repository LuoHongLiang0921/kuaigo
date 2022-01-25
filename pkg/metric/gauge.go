// @Description 仪表盘 测量值，或瞬时记录值，可以增加，也可以减少。例如：一个视频的同时观看人数，当前运行的进程数

package metric

import "github.com/prometheus/client_golang/prometheus"

// GaugeVecOpts ...
type GaugeVecOpts struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
}

type gaugeVec struct {
	*prometheus.GaugeVec
}

// Build ...
func (opts GaugeVecOpts) Build() *gaugeVec {
	vec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: opts.Namespace,
			Subsystem: opts.Subsystem,
			Name:      opts.Name,
			Help:      opts.Help,
		}, opts.Labels)
	prometheus.MustRegister(vec)
	return &gaugeVec{
		GaugeVec: vec,
	}
}

// NewGaugeVec ...
func NewGaugeVec(name string, labels []string) *gaugeVec {
	return GaugeVecOpts{
		Namespace: DefaultNamespace,
		Name:      name,
		Help:      name,
		Labels:    labels,
	}.Build()
}

// Inc ...
func (gv *gaugeVec) Inc(labels ...string) {
	gv.WithLabelValues(labels...).Inc()
}

// Add ...
func (gv *gaugeVec) Add(v float64, labels ...string) {
	gv.WithLabelValues(labels...).Add(v)
}

// Set ...
func (gv *gaugeVec) Set(v float64, labels ...string) {
	gv.WithLabelValues(labels...).Set(v)
}

//variableCounterVec 标签可变
type variableGaugeVec struct {
	m *metric
}

// GaugeVec  namespace为 tabby，subsystem为 go
var GaugeVec = &variableGaugeVec{
	m: Option{
		Namespace: "tabby",
		Subsystem: "go",
		mt:        gv,
	}.Build(),
}

//NewVariableGaugeVec ...
func NewVariableGaugeVec(r Option) *variableGaugeVec {
	return &variableGaugeVec{m: r.Build()}
}

// Inc ...
func (gv *variableGaugeVec) Inc(name string, kv interface{}) {
	lbNames, lbValues := genLabels(kv)
	//lbNames
	v := gv.m.loadOrStore(name, lbNames)
	if v != nil {
		vv := v.(*prometheus.GaugeVec)
		vv.WithLabelValues(lbValues...).Inc()
	}
}

// Add ...
func (gv *variableGaugeVec) Add(name string, kv interface{}, v float64) {
	lbNames, lbValues := genLabels(kv)
	//lbNames
	v1 := gv.m.loadOrStore(name, lbNames)
	if v1 != nil {
		vv := v1.(*prometheus.GaugeVec)
		vv.WithLabelValues(lbValues...).Add(v)
	}
}

// Set ...
func (gv *variableGaugeVec) Set(name string, kv interface{}, v float64) {
	lbNames, lbValues := genLabels(kv)
	//lbNames
	v1 := gv.m.loadOrStore(name, lbNames)
	if v1 != nil {
		vv := v1.(*prometheus.GaugeVec)
		vv.WithLabelValues(lbValues...).Set(v)
	}
}

// WithHelp 设置help 信息
func (gv *variableGaugeVec) WithHelp(help string) *variableGaugeVec {
	gv.m.opt.Help = help
	return gv
}

// WithName 设置name
func (gv *variableGaugeVec) WithName(name string) *variableGaugeVec {
	gv.m.opt.Name = name
	return gv
}
