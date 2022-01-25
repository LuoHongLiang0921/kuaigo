package gorm

import (
	"sync"
)

var instances = sync.Map{}

// Range
// 	@Description 遍历所有db实例
//	@Param fn 执行函数
func Range(fn func(name string, db *DB) bool) {
	instances.Range(func(key, val interface{}) bool {
		return fn(key.(string), val.(*DB))
	})
}

// Stats
// 	@Description 获取所有db 统计信息
// 	@Return stats 统计信息
func Stats() (stats map[string]interface{}) {
	stats = make(map[string]interface{})
	instances.Range(func(key, val interface{}) bool {
		name := key.(string)
		db := val.(*DB)

		stats[name] = db.DB().Stats()
		return true
	})

	return
}
