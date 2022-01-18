// @description
// @author yixia
// Copyright 2021 sndks.com. All rights reserved.
// @datetime 2021/1/14 5:21 下午
// @lastmodify 2021/1/14 5:21 下午

package trace_test

import (
	"context"
	"fmt"
	"time"

	"git.bbobo.com/framework/tabby/pkg/trace"
)

func ExampleTraceFunc() {
	// 1. 从配置文件中初始化
	process1 := func(ctx context.Context) {
		span, ctx := trace.StartSpanFromContext(ctx, "process1")
		defer span.Finish()

		// todo something
		fmt.Println("err", ctx.Err())
		time.Sleep(time.Second)
	}

	process2 := func(ctx context.Context) {
		span, ctx := trace.StartSpanFromContext(ctx, "process2")
		defer span.Finish()
		process1(ctx)
	}

	process3 := func(ctx context.Context) {
		span, ctx := trace.StartSpanFromContext(ctx, "process3")
		defer span.Finish()
		process2(ctx)
	}

	process3(context.Background())
}