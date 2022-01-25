package gorm

import (
	"context"
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/metric"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/kcolor"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"strconv"
	"time"
)

// Handler ...
type Handler func(*Scope)

// Interceptor ...
type Interceptor func(context.Context, *DSN, string, *Config) func(next Handler) Handler

// debugInterceptor
// 	@Description debug 拦截器
//	@Param ctx 上下文
//	@Param dsn 字符串
//	@Param op 操作
//	@Param options 配置项
// 	@Return func(Handler) Handler 装饰后的debug 拦截器
func debugInterceptor(ctx context.Context, dsn *DSN, op string, options *Config) func(Handler) Handler {
	return func(next Handler) Handler {
		return func(scope *Scope) {
			fmt.Printf("%-50s[%s] => %s\n", kcolor.Green(dsn.Addr+"/"+dsn.DBName), time.Now().Format("04:05.000"), kcolor.Green("Send: "+logSQL(scope.SQL, scope.SQLVars, true)))
			next(scope)
			if scope.HasError() {
				fmt.Printf("%-50s[%s] => %s\n", kcolor.Red(dsn.Addr+"/"+dsn.DBName), time.Now().Format("04:05.000"), kcolor.Red("Erro: "+scope.DB().Error.Error()))
			} else {
				fmt.Printf("%-50s[%s] => %s\n", kcolor.Green(dsn.Addr+"/"+dsn.DBName), time.Now().Format("04:05.000"), kcolor.Green("Affected: "+strconv.Itoa(int(scope.DB().RowsAffected))))
			}
		}
	}
}

// metricInterceptor
// 	@Description 指标拦截器
//	@Param ctx 上下文
//	@Param dsn DSN 结构
//	@Param op 操作
//	@Param options 配置项
// 	@Return func(Handler) Handler 装饰后的指标拦截器
func metricInterceptor(ctx context.Context, dsn *DSN, op string, options *Config) func(Handler) Handler {
	return func(next Handler) Handler {
		return func(scope *Scope) {
			beg := time.Now()
			next(scope)
			cost := time.Since(beg)

			// error metric
			if scope.HasError() {
				metric.LibHandleCounter.WithLabelValues(metric.TypeGorm, dsn.DBName+"."+scope.TableName(), dsn.Addr, "ERR").Inc()
				// todo sql语句，需要转换成脱密状态才能记录到日志
				if scope.DB().Error != ErrRecordNotFound {
					options.logger.WithContext(ctx).Error("mysql err", klog.FieldErr(scope.DB().Error), klog.FieldName(dsn.DBName+"."+scope.TableName()), klog.FieldMethod(op))
				} else {
					options.logger.WithContext(ctx).Warn("record not found", klog.FieldErr(scope.DB().Error), klog.FieldName(dsn.DBName+"."+scope.TableName()), klog.FieldMethod(op))
				}
			} else {
				metric.LibHandleCounter.Inc(metric.TypeGorm, dsn.DBName+"."+scope.TableName(), dsn.Addr, "OK")
			}

			metric.LibHandleHistogram.WithLabelValues(metric.TypeGorm, dsn.DBName+"."+scope.TableName(), dsn.Addr).Observe(cost.Seconds())

			if options.SlowThreshold > time.Duration(0) && options.SlowThreshold < cost {
				options.logger.WithContext(ctx).Error(
					"slow",
					klog.FieldErr(errSlowCommand),
					klog.FieldMethod(op),
					klog.FieldExtMessage(logSQL(scope.SQL, scope.SQLVars, options.DetailSQL)),
					klog.FieldAddr(dsn.Addr),
					klog.FieldName(dsn.DBName+"."+scope.TableName()),
					klog.FieldCost(cost),
				)
			}
		}
	}
}

func logSQL(sql string, args []interface{}, containArgs bool) string {
	if containArgs {
		return bindSQL(sql, args)
	}
	return sql
}
