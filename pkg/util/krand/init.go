// @Description 随机包初始化

package krand

import (
	"math/rand"
	"sync"
	"time"
)

var (
	timeBase     = time.Date(1582, time.October, 15, 0, 0, 0, 0, time.UTC).Unix()
	hardwareAddr []byte
	clockSeq     uint32
	randInstance *rand.Rand
	r            = rand.New(rand.NewSource(time.Now().UnixNano()))
	mu           sync.Mutex
)
