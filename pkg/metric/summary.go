// @Description 摘要 summary和histogram类似也会产生三个指标，分别是_count，_sum，和{quantile}

package metric

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// SummaryVecOpts ...
type SummaryVecOpts struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
}

type summaryVec struct {
	*prometheus.SummaryVec
}

// Build ...
func (opts SummaryVecOpts) Build() *summaryVec {
	vec := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: opts.Namespace,
			Subsystem: opts.Subsystem,
			Name:      opts.Name,
			Help:      opts.Help,
		}, opts.Labels)
	prometheus.MustRegister(vec)
	return &summaryVec{
		SummaryVec: vec,
	}
}

// Observe ...
func (summary *summaryVec) Observe(v float64, labels ...string) {
	summary.WithLabelValues(labels...).Observe(v)
}

//variableSummaryVec 标签可变
type variableSummaryVec struct {
	m *metric
}

// SummaryVec  namespace为 tabby，subsystem为 go
var SummaryVec = &variableSummaryVec{
	m: Option{
		Namespace: "tabby",
		Subsystem: "go",
		mt:        sv,
	}.Build(),
}

//NewVariableSummaryVec ...
func NewVariableSummaryVec(r Option) *variableSummaryVec {
	return &variableSummaryVec{m: r.Build()}
}

// Timing ...
// kv 可以是[]string, map[string]string
func (s *variableSummaryVec) Timing(name string, kv interface{}, startAt time.Time) {
	lbNames, lbValues := genLabels(kv)
	//lbNames
	v := s.m.loadOrStore(name, lbNames)
	if v != nil {
		vv := v.(*prometheus.SummaryVec)
		vv.WithLabelValues(lbValues...).Observe(time.Since(startAt).Seconds())
	}
}

// WithHelp 设置help 信息
func (s *variableSummaryVec) WithHelp(help string) *variableSummaryVec {
	s.m.opt.Help = help
	return s
}
