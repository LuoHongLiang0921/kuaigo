// @Description
// @Author shiyibo
// @Copyright 2021 sndks.com. All rights reserved.
// @Datetime 2021/7/29 2:45 下午

package property

import "testing"

func BenchmarkBuildInExpress(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BuildInExpress(true)
	}
}

func BenchmarkBuildInExpressFalse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BuildInExpress(false)
	}
}
