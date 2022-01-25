// @Description 日志存放到redis
// @Author yixia
// @Copyright 2021 sndks.com. All rights reserved.
// @LastModify 2021/1/14 5:21 下午

package redis

import (
	"github.com/go-redis/redis"
)

// RedisLog redis 日志存储
type RedisLog struct {
	LogKey string
	client *redis.Client
}

// NewRedisLog
// 	@Description: 初始化redis 写组件实例
//	@Param logKey redis key
//	@Param client redis 实例
// 	@return *RedisLog
func NewRedisLog(logKey string, client *redis.Client) *RedisLog {
	return &RedisLog{LogKey: logKey, client: client}
}

// 	SetLogKey
// 	@Description 新的实例写的
// 	@Receiver l
// 	@Return *RedisLog
func (l *RedisLog) SetLogKey(logKey string) *RedisLog {
	newInstance := l.clone()
	newInstance.LogKey = logKey
	return newInstance
}

func (l *RedisLog) clone() *RedisLog {
	copy := *l
	return &copy
}

// Write
// 	@Description: 写入到日志, 通过 lpush 命令写到redis中
// 	@receiver l RedisLog
//	@Param p 内容
// 	@return n 字节数组位置
// 	@return err 错误
func (l *RedisLog) Write(p []byte) (n int, err error) {
	return len(p), l.client.LPush(l.LogKey, p).Err()
}

// Sync
// 	@Description 刷新缓存中数据
// 	@Receiver l RedisLog
// 	@Return error 错误
func (l *RedisLog) Sync() error {
	return nil
}

// Close
// 	@Description 关闭redis
// 	@Receiver l RedisLog
// 	@Return error 错误
func (l RedisLog) Close() error {
	return l.client.Close()
}
