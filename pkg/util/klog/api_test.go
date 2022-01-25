package klog_test

import (
	"context"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog/console"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog/file"
	"github.com/LuoHongLiang0921/kuaigo/test"
	"sync"
	"testing"


	"github.com/stretchr/testify/assert"
)

func TestInfo(t *testing.T) {
	ctx := klog.WithCommonLog(context.Background(), klog.Common{AppId: 104, ServiceSource: "task", ServiceName: "push-service"})
	klog.WithContext(ctx).Info("test")
}

func TestWithContext(t *testing.T) {
	ctx := context.Background()
	//RunningLogger.SetLevel(zap.DebugLevel)
	//KuaigoLogger.SetLevel(zap.DebugLevel)
	//ErrorLogger.SetLevel(zap.DebugLevel)
	//AccessLogger.SetLevel(zap.DebugLevel)
	//TaskLogger.SetLevel(zap.DebugLevel)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx := klog.AccessLoggerContext(context.Background())
		klog.WithContext(ctx).Info("test access test2")
		klog.WithContext(ctx).Infof("test access %v", "test2 ")

		klog.WithContext(ctx).Debugf("test access %v", "test2 ")
		klog.WithContext(ctx).Debug("test access test2")

		klog.WithContext(ctx).Error("test access test2")
		klog.WithContext(ctx).Errorf("test access %v", "test2 ")
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx := klog.ErrorLoggerContext(context.Background())
		klog.WithContext(ctx).Info("test error test2")
		klog.WithContext(ctx).Infof("test error %v", "test3 ")

		klog.WithContext(ctx).Debugf("test error %v", "test3 ")
		klog.WithContext(ctx).Debug("test error test2")

		klog.WithContext(ctx).Error("test error test2")
		klog.WithContext(ctx).Errorf("test error %v", "test3 ")
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx := klog.RunningLoggerContext(context.Background())
		klog.WithContext(ctx).Info("test running test4")
		klog.WithContext(ctx).Infof("test running %v", "test4 ")

		klog.WithContext(ctx).Debug("test running test4")
		klog.WithContext(ctx).Debugf("test running %v", "test4 ")

		klog.WithContext(ctx).Error("test running test4")
		klog.WithContext(ctx).Errorf("test running %v", "test4 ")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx := klog.TaskLoggerContext(ctx, klog.WithServiceSource(klog.ServiceSourceTask))
		klog.WithContext(ctx).Info("test task 1")
		klog.WithContext(ctx).Infof("test task %v", "1")

		klog.WithContext(ctx).Debug("test task")
		klog.WithContext(ctx).Debugf("test task %v", "1")

		klog.WithContext(ctx).Error("test task 1")
		klog.WithContext(ctx).Errorf("test task %v", "1")

		klog.WithContext(ctx).Warn("test task warn")
		klog.WithContext(ctx).Warnf("test task warnf %v", 1)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx := klog.TaskLoggerContext(ctx, klog.WithServiceSource(klog.ServiceSourceTask))
		klog.RunningLogger.WithContext(ctx).Info("test running to task 1")
		klog.RunningLogger.WithContext(ctx).Infof("test running to task %v", "1")
		klog.RunningLogger.WithContext(ctx).Debug("test running to task")
		klog.RunningLogger.WithContext(ctx).Debugf("test running to task %v", "1")
		klog.RunningLogger.WithContext(ctx).Error("test running to task 1")
		klog.RunningLogger.WithContext(ctx).Errorf("test running to task %v", "1")
		klog.RunningLogger.WithContext(ctx).Warn("test running to task warn")
		klog.RunningLogger.WithContext(ctx).Warnf("test running to task warnf %v", 1)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		klog.KuaigoLogger.WithContext(ctx).Info("test tabby 1")
		klog.KuaigoLogger.WithContext(ctx).Infof("test tabby %v", "1")
		klog.KuaigoLogger.WithContext(ctx).Debug("test tabby")
		klog.KuaigoLogger.WithContext(ctx).Debugf("test tabby %v", "1")
		klog.KuaigoLogger.WithContext(ctx).Error("test tabby 1")
		klog.KuaigoLogger.WithContext(ctx).Errorf("test tabby %v", "1")
		klog.KuaigoLogger.WithContext(ctx).Warn("test tabby warn")
		klog.KuaigoLogger.WithContext(ctx).Warnf("test tabby warnf %v", 1)
	}()
	wg.Wait()
	klog.FlushAll()
}

func TestError(t *testing.T) {
	rootCtx := context.Background()
	ctx := klog.TaskLoggerContext(rootCtx, klog.WithServiceSource(klog.ServiceSourceTask))
	klog.WithContext(ctx).Errorf("test err %v", "alarm")
	klog.WithContext(ctx).Error("test err alarm")
	klog.WithContext(ctx).Info("test info alarm")
	klog.FlushAll()
	select {}
}

func TestInitLogger(t *testing.T) {
	// 需要设置 CONFIG_FILE_ADDR=test/testdata/test_conf.yaml
	if assert.NoError(t, test.InitTestForFile()) {
		console.RegisterOutputCreatorHandler()
		file.RegisterOutputCreatorHandler()
		//InitLogger("v2", "tabby-test")

		//redis.RegisterOutputCreatorHandler()
		running := klog.StdConfig("running").WithConfigVersion("v2").WithServiceName("tabby-test").Build()
		//for i := 0; i < 10000; i++ {

		running.Infof("test info 1")
		running.Error("test info 1,error")
		running.Debug("test info 1,error")
		running.Debugf("test info 1,error")
		running.Panic("test info 1,panic")

		running.Flush()

	}
	//RunningLogger.Flush()
}
