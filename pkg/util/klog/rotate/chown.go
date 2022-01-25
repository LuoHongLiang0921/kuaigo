// @Description
// @Author yixia
// @Copyright 2021 sndks.com. All rights reserved.
// @LastModify 2021/1/14 5:21 下午

// +build !linux

package rotate

import (
	"os"
)

func chown(_ string, _ os.FileInfo) error {
	return nil
}
