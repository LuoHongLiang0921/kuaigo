// @description
// @author yixia
// Copyright 2021 sndks.com. All rights reserved.
// @datetime 2021/1/14 5:21 下午
// @lastmodify 2021/1/14 5:21 下午

package ktime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimer(t *testing.T) {
	var testWheel = NewRashTimer(1 * time.Millisecond)
	t1 := testWheel.NewTimer(500 * time.Millisecond)

	before := time.Now()
	<-t1.C
	after := time.Now()

	assert.True(t, after.Sub(before) < time.Millisecond*600)
}
