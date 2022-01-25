// @Description
// @Author yixia
// @Copyright 2021 sndks.com. All rights reserved.
// @LastModify 2021/1/14 5:21 下午

package rotate_test

import (
	"log"

	rotate2 "git.bbobo.com/framework/tabby/pkg/util/xlog/rotate"
)

// To use rotate with the standard library's log package, just pass it into
// the SetOutput function when your application starts.
func Example() {
	log.SetOutput(&rotate2.Logger{
		Filename:   "/var/log/myapp/foo.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   // days
		Compress:   true, // disabled by default
	})
}
