// @Description metric

package metric

import (
	"github.com/LuoHongLiang0921/kuaigo/pkg"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/server/governor"
	"math"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// TypeHTTP ...
	TypeHTTP = "http"
	// TypeGRPCUnary ...
	TypeGRPCUnary = "unary"
	// TypeGRPCStream ...
	TypeGRPCStream = "stream"
	// TypeRedis ...
	TypeRedis = "redis"
	TypeGorm  = "gorm"
	// TypeRocketMQ ...
	TypeRocketMQ = "rocketmq"
	// TypeWebsocket ...
	TypeWebsocket = "ws"

	// TypeMySQL ...
	TypeMySQL = "mysql"

	// CodeJob
	CodeJobSuccess = "ok"
	// CodeJobFail ...
	CodeJobFail = "fail"
	// CodeJobReentry ...
	CodeJobReentry = "reentry"

	// CodeCache
	CodeCacheMiss = "miss"
	// CodeCacheHit ...
	CodeCacheHit = "hit"

	// Namespace
	DefaultNamespace = "tabby"
)

var (
	// ServerHandleCounter ...
	ServerHandleCounter = CounterVecOpts{
		Namespace: DefaultNamespace,
		Name:      "server_handle_total",
		Labels:    []string{"type", "method", "peer", "code"},
	}.Build()

	// ServerHandleHistogram ...
	ServerHandleHistogram = HistogramVecOpts{
		Namespace: DefaultNamespace,
		Name:      "server_handle_seconds",
		Labels:    []string{"type", "method", "peer"},
	}.Build()

	// ClientHandleCounter ...
	ClientHandleCounter = CounterVecOpts{
		Namespace: DefaultNamespace,
		Name:      "client_handle_total",
		Labels:    []string{"type", "name", "method", "peer", "code"},
	}.Build()

	// ClientHandleHistogram ...
	ClientHandleHistogram = HistogramVecOpts{
		Namespace: DefaultNamespace,
		Name:      "client_handle_seconds",
		Labels:    []string{"type", "name", "method", "peer"},
	}.Build()

	// JobHandleCounter ...
	JobHandleCounter = CounterVecOpts{
		Namespace: DefaultNamespace,
		Name:      "job_handle_total",
		Labels:    []string{"type", "name", "code"},
	}.Build()

	// JobHandleHistogram ...
	JobHandleHistogram = HistogramVecOpts{
		Namespace: DefaultNamespace,
		Name:      "job_handle_seconds",
		Labels:    []string{"type", "name"},
	}.Build()

	LibHandleHistogram = HistogramVecOpts{
		Namespace: DefaultNamespace,
		Name:      "lib_handle_seconds",
		Labels:    []string{"type", "method", "address"},
	}.Build()
	// LibHandleCounter ...
	LibHandleCounter = CounterVecOpts{
		Namespace: DefaultNamespace,
		Name:      "lib_handle_total",
		Labels:    []string{"type", "method", "address", "code"},
	}.Build()

	LibHandleSummary = SummaryVecOpts{
		Namespace: DefaultNamespace,
		Name:      "lib_handle_stats",
		Labels:    []string{"name", "status"},
	}.Build()

	// CacheHandleCounter ...
	CacheHandleCounter = CounterVecOpts{
		Namespace: DefaultNamespace,
		Name:      "cache_handle_total",
		Labels:    []string{"type", "name", "action", "code"},
	}.Build()

	// CacheHandleHistogram ...
	CacheHandleHistogram = HistogramVecOpts{
		Namespace: DefaultNamespace,
		Name:      "cache_handle_seconds",
		Labels:    []string{"type", "name", "action"},
	}.Build()

	// BuildInfoGauge ...
	BuildInfoGauge = GaugeVecOpts{
		Namespace: DefaultNamespace,
		Name:      "build_info",
		Labels:    []string{"name", "aid", "mode", "region", "zone", "app_version", "tabby_version", "start_time", "build_time", "go_version"},
	}.Build()
)

// prometheus related param
var (
	Namespace              = "tabby"
	Subsystem              = "go"
	helpDescriptionMap     = make(map[string]string)
	histogramDefaultBucket = []float64{0.003, 0.005, 0.01, 0.03, 0.05, 0.07, 0.09, 0.1, 0.15, 0.2, 0.25, 0.3, 0.5, 0.7, 1, 1.5, 2, math.Inf(+1)}
	histogramBuckets       = make(map[string][]float64)
	constLabels            = make(map[string]string)
)

func init() {
	BuildInfoGauge.WithLabelValues(
		pkg.GetAppName(),
		pkg.GetAppID(),
		pkg.GetAppMode(),
		pkg.GetAppRegion(),
		pkg.GetAppZone(),
		pkg.GetAppVersion(),
		pkg.GetTabbyVersion(),
		pkg.GetStartTime(),
		pkg.GetBuildTime(),
		pkg.GetGoVersion(),
	).Set(float64(time.Now().UnixNano() / 1e6))

	governor.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})
}

type metricType string

const (
	cv metricType = "counterVec"
	gv metricType = "gaugeVec"
	hv metricType = "histogramVec"
	sv metricType = "summaryVec"
)

type metric struct {
	mt  metricType
	mu  *sync.RWMutex
	opt *Option
	bag map[string]interface{}
}

// metric Option
type Option struct {
	Namespace   string    // 命名空间
	Subsystem   string    // 子系统
	Name        string    // 名字
	Help        string    // 描述
	Labels      []string  // 标签
	Buckets     []float64 // histogram buckets
	ConstLabels map[string]string
	mt          metricType //prometheus 类型
}

//  newMetric 返回根据Option 生成的
func (o Option) Build() *metric {
	if o.ConstLabels == nil {
		o.ConstLabels = constLabels
	}
	switch o.mt {
	case hv:
		if o.Buckets == nil {
			o.Buckets = histogramDefaultBucket
		}
	}
	return &metric{
		mt:  o.mt,
		mu:  &sync.RWMutex{},
		opt: &o,
		bag: make(map[string]interface{}),
	}
}

// getDescription result can not be empty
func getDescription(name string) string {
	if v := helpDescriptionMap[name]; v != "" {
		return v
	}

	return name
}
func (m *metric) gen(name string, labels []string) interface{} {
	switch m.mt {
	case cv:
		counterVec := prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace:   m.opt.Namespace,
				Subsystem:   m.opt.Subsystem,
				Name:        name,
				Help:        m.opt.Help,
				ConstLabels: m.opt.ConstLabels,
			},
			labels,
		)

		err := prometheus.Register(counterVec)
		if err != nil {
			return nil
		}

		return counterVec
	case gv:
		gaugeVec := prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   m.opt.Namespace,
				Subsystem:   m.opt.Subsystem,
				Name:        name,
				Help:        m.opt.Help,
				ConstLabels: m.opt.ConstLabels,
			}, labels,
		)

		err := prometheus.Register(gaugeVec)
		if err != nil {
			return nil
		}

		return gaugeVec
	case hv:
		opts := prometheus.HistogramOpts{
			Namespace:   m.opt.Namespace,
			Subsystem:   m.opt.Subsystem,
			Name:        name,
			Help:        m.opt.Help,
			ConstLabels: m.opt.ConstLabels,
			Buckets:     m.opt.Buckets,
		}

		if v, ok := histogramBuckets[name]; ok {
			opts.Buckets = v
		}

		histogramVec := prometheus.NewHistogramVec(opts, labels)
		err := prometheus.Register(histogramVec)
		if err != nil {
			return nil
		}

		return histogramVec
	case sv:
		summaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Namespace:   m.opt.Namespace,
			Subsystem:   m.opt.Subsystem,
			Name:        name,
			Help:        m.opt.Help,
			ConstLabels: m.opt.ConstLabels,
		}, labels)

		err := prometheus.Register(summaryVec)
		if err != nil {
			return nil
		}

		return summaryVec
	default:
		return nil
	}
}

func (m *metric) loadOrStore(name string, labels []string) interface{} {
	m.mu.RLock()

	if v, ok := m.bag[name]; ok {
		m.mu.RUnlock()
		return v
	}

	m.mu.RUnlock()

	// try again
	m.mu.Lock()
	if v, ok := m.bag[name]; ok {
		m.mu.Unlock()
		return v
	}

	v := m.gen(name, labels)
	m.bag[name] = v
	m.mu.Unlock()
	return v
}

// genLabels
func genLabels(kv interface{}) ([]string, []string) {
	var lbNames, lbValues []string
	switch v := kv.(type) {
	case []string:
		if l := len(v) % 2; l != 0 {
			v = v[:l-1]
		}
		for i, l := 0, len(v); i < l; i = i + 2 {
			lbNames = append(lbNames, v[i])
			lbValues = append(lbValues, v[i+1])
		}
	case map[string]string:
		for k := range v {
			lbNames = append(lbNames, k)
		}

		sort.Strings(lbNames)

		for i := range lbNames {
			lbValues = append(lbValues, v[lbNames[i]])
		}
	}
	return lbNames, lbValues
}
