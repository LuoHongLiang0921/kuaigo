package gorm

import (
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/server/governor"
	"github.com/LuoHongLiang0921/kuaigo/pkg/metric"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var (
	_logger = klog.KuaigoLogger.With(klog.FieldMod("gorm"))
)

func init() {
	type gormStatus struct {
		Gorms map[string]interface{} `json:"gorms"`
	}
	var rets = gormStatus{
		Gorms: make(map[string]interface{}, 0),
	}
	governor.HandleFunc("/debug/gorm/stats", func(w http.ResponseWriter, r *http.Request) {
		rets.Gorms = Stats()
		_ = jsoniter.NewEncoder(w).Encode(rets)
	})
	go monitor()
}

// monitor
// 	@Description 每十秒获取数据库运行状况
func monitor() {
	for {
		time.Sleep(time.Second * 10)
		Range(func(name string, db *DB) bool {
			stats := db.DB().Stats()
			metric.LibHandleSummary.Observe(float64(stats.Idle), name, "idle")
			metric.LibHandleSummary.Observe(float64(stats.InUse), name, "inuse")
			metric.LibHandleSummary.Observe(float64(stats.WaitCount), name, "wait")
			metric.LibHandleSummary.Observe(float64(stats.OpenConnections), name, "conns")
			metric.LibHandleSummary.Observe(float64(stats.MaxOpenConnections), name, "max_open_conns")
			metric.LibHandleSummary.Observe(float64(stats.MaxIdleClosed), name, "max_idle_closed")
			metric.LibHandleSummary.Observe(float64(stats.MaxLifetimeClosed), name, "max_lifetime_closed")
			return true
		})
	}
}
