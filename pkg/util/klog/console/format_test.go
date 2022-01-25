package console

import (
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog/property"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog/zap"
	"testing"


	"github.com/stretchr/testify/assert"

)

func TestFormat_String(t *testing.T) {
	property.RegisterProperty("key", "k")
	property.RegisterProperty("anther", "anther")
	strVal := property.BuildInExpress(true)
	f := Format("ai=%{ai} %{msg} %{ti} key=%{key} " + strVal)
	fer := NewFormatter(f)
	result := fer.String(
		zap.String("ai", "t"),
		zap.String("msg", "msg-test"),
		zap.String("key", "test k"),
		zap.String("anther", "test k"),
	)
	fmt.Println(result)
	assert.Equal(t, "ai=t msg-test  key=test k", result)
}

func BenchmarkFormat_string(b *testing.B) {
	property.RegisterProperty("key", "k")
	f := Format("ai=%{ai} %{msg} %{ti} key=%{key}")
	fer := NewFormatter(f)
	for i := 0; i < b.N; i++ {
		fer.String(
			zap.String("ai", "t"),
			zap.String("msg", "msg-test"),
			zap.String("key", "test k"),
		)
	}
}

func BenchmarkFormat_stringall(b *testing.B) {
	property.RegisterProperty("key", "k")
	f := Format("appId=%{appId} serviceName=%{serviceName} serviceSource=%{serviceSource} traceId=%{traceId} fileName=%{fileName} line=%{line} requestIp=%{requestIp} requestUri=%{requestUri} error=%{error}")
	fer := NewFormatter(f)
	for i := 0; i < b.N; i++ {
		fer.String(
			zap.String("ai", "t"),
			zap.String("msg", "msg-test"),
			zap.String("key", "test k"),
			zap.String("appId", "test k"),
			zap.String("serviceName", "test k"),
			zap.String("serviceSource", "test k"),
			zap.String("traceId", "test k"),
		)
	}
}
