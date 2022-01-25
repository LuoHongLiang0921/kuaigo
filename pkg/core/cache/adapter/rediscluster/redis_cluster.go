package rediscluster

import (
	"context"
	"errors"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/cache/config"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

type redisClusterAdapter struct {
	ctx       context.Context
	config    *config.CacheConfig
	client    redis.Cmdable
	openOnce  sync.Once
	hasOpen   bool
	openError error
}

// NewRedisClusterAdapter
// 	@Description Mysql数据库操作类构建函数
//  @Param ctx 上下文Context
//	@Param config 数据库配置
// 	@Return IDataBaseAdapter
func NewRedisClusterAdapter(ctx context.Context, conf *config.CacheConfig) config.IAdvanceCache {
	ret := &redisClusterAdapter{
		ctx:       ctx,
		config:    conf,
		hasOpen:   false,
		openError: errors.New("redis " + conf.Name + " open error"),
	}
	if conf.AutoConnect {
		ret.Open(ctx)
	}
	return ret
}

// WithContext
// 	@Description
// 	@Receiver r
//	@Param ctx
// 	@Return config.ICache
func (r *redisClusterAdapter) WithContext(ctx context.Context) config.ICache {
	if ctx == nil {
		return r
	}
	newR := r.clone()
	newR.ctx = ctx
	return newR
}

// WithAdvanceContext
// 	@Description
// 	@Receiver r redisAdapter
//	@Param ctx 上下文
// 	@Return *redisAdapter
func (r *redisClusterAdapter) WithAdvanceContext(ctx context.Context) config.IAdvanceCache {
	if ctx == nil {
		return r
	}
	newR := r.clone()
	newR.ctx = ctx
	return newR
}

func (r *redisClusterAdapter) clone() *redisClusterAdapter {
	copy := redisClusterAdapter{
		ctx:       r.ctx,
		config:    r.config,
		client:    r.client,
		openOnce:  sync.Once{},
		hasOpen:   r.hasOpen,
		openError: r.openError,
	}
	return &copy
}

// Get 从缓存获取string
// 	@Description 通过 `get key` 命令获取key字符串值
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
// 	@return string key 对应的字符串值
func (r *redisClusterAdapter) Get(key string) string {
	if !r.checkOpen() {
		return ""
	}
	var mes string
	strObj := r.client.Get(key)
	if err := strObj.Err(); err != nil {
		mes = ""
	} else {
		mes = strObj.Val()
	}
	return mes
}

// GetRaw
// 	@Description 通过 `get key` 命令获取key字节数组值
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
// 	@return []byte 字节数组
// 	@return error 错误
func (r *redisClusterAdapter) GetRaw(key string) ([]byte, error) {
	if !r.checkOpen() {
		return []byte{}, r.openError
	}
	c, err := r.client.Get(key).Bytes()
	if err != nil && err != redis.Nil {
		return []byte{}, err
	}
	return c, nil
}

// MGet
// 	@Description 通过 `get key` 命令获取key字节数组值
// 	@Receiver r redisClusterAdapter
//	@Param keys 键名字数组
// 	@return []string
// 	@return error
func (r *redisClusterAdapter) MGet(keys ...string) ([]string, error) {
	if !r.checkOpen() {
		return []string{}, r.openError
	}
	sliceObj := r.client.MGet(keys...)
	if err := sliceObj.Err(); err != nil && err != redis.Nil {
		return []string{}, err
	}
	tmp := sliceObj.Val()
	strSlice := make([]string, 0, len(tmp))
	for _, v := range tmp {
		if v != nil {
			strSlice = append(strSlice, v.(string))
		} else {
			strSlice = append(strSlice, "")
		}
	}
	return strSlice, nil
}

// MGets
// 	@Description 通过 `MGET key [key ...]` 获取指定键的值
// 	@Receiver r redisClusterAdapter
//	@Param keys 键名字数组
// 	@return []interface{} 键对应的值数组
// 	@return error 错误
func (r *redisClusterAdapter) MGets(keys []string) ([]interface{}, error) {
	if !r.checkOpen() {
		return []interface{}{}, r.openError
	}
	ret, err := r.client.MGet(keys...).Result()
	if err != nil && err != redis.Nil {
		return []interface{}{}, err
	}
	return ret, nil
}

// Set 设置redis的string
// 	@Description 通过 `SET key value [EX seconds]` 设置键对应的值
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param value 值
//	@Param expire 过期时间
// 	@return bool 是否设置成功布尔值
func (r *redisClusterAdapter) Set(key string, value interface{}, expire time.Duration) bool {
	if !r.checkOpen() {
		return false
	}
	err := r.client.Set(key, value, expire).Err()
	return err == nil
}

// SetWithErr
// 	@Description 通过 `SET key value [EX seconds]` 设置键对应的值，并返回错误值
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param value 键对应的值
//	@Param expire 过期时间
// 	@Return error
func (r *redisClusterAdapter) SetWithErr(key string, value interface{}, expire time.Duration) error {
	if !r.checkOpen() {
		return r.openError
	}
	err := r.client.Set(key, value, expire).Err()
	return err
}

// SetNx
// 	@Description 通过 `SETNX key value` 命令设置键不存在的值
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param value 键值
//	@Param expiration 过期时间
// 	@return bool 是否设置成功
func (r *redisClusterAdapter) SetNx(key string, value interface{}, expiration time.Duration) bool {
	if !r.checkOpen() {
		return false
	}
	result, err := r.client.SetNX(key, value, expiration).Result()
	if err != nil {
		return false
	}
	return result
}

// SetNxWithErr
// 	@Description 通过 `SETNX key value` 命令设置键不存在的值，并返回错误
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param value 键值
//	@Param expiration 过期时间
// 	@return bool 是否设置成功
// 	@return error 错误
func (r *redisClusterAdapter) SetNxWithErr(key string, value interface{}, expiration time.Duration) (bool, error) {
	if !r.checkOpen() {
		return false, r.openError
	}
	result, err := r.client.SetNX(key, value, expiration).Result()
	return result, err
}

// Incr redis自增
// 	@Description 通过 `INCR key` 自增键值
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
// 	@return bool
func (r *redisClusterAdapter) Incr(key string) bool {
	if !r.checkOpen() {
		return false
	}
	err := r.client.Incr(key).Err()
	return err == nil
}

// IncrWithErr ...
// IncrWithErr
// 	@Description 通过 `INCR key` 自增键值
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
// 	@Return int64 自增后的值
// 	@Return error
func (r *redisClusterAdapter) IncrWithErr(key string) (int64, error) {
	if !r.checkOpen() {
		return 0, r.openError
	}
	ret, err := r.client.Incr(key).Result()
	return ret, err
}

// IncrBy 将 key 所储存的值加上增量 increment
// IncrBy
// 	@Description 通过 `INCRBY key increment` 命令增加增量值
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param increment 增量值
// 	@Return int64 加上增量后的值
// 	@Return error 错误
func (r *redisClusterAdapter) IncrBy(key string, increment int64) (int64, error) {
	if !r.checkOpen() {
		return 0, r.openError
	}
	intObj := r.client.IncrBy(key, increment)
	if err := intObj.Err(); err != nil {
		return 0, err
	}
	return intObj.Val(), nil
}

// Decr
// 	@Description 通过 `DECR key` 自减键对应的值
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
// 	@Return bool
func (r *redisClusterAdapter) Decr(key string) bool {
	if !r.checkOpen() {
		return false
	}
	err := r.client.Decr(key).Err()
	return err == nil
}

// Del
// 	@Description 删除key
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
// 	@Return int64 删除key成功的个数
func (r *redisClusterAdapter) Del(key ...string) int64 {
	if !r.checkOpen() {
		return 0
	}
	result, err := r.client.Del(key...).Result()
	if err != nil {
		return 0
	}
	return result
}

// DelWithErr
// 	@Description 删除key，并返回错误
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
// 	@Return int64 删除key成功的个数
// 	@Return error 错误
func (r *redisClusterAdapter) DelWithErr(key string) (int64, error) {
	if !r.checkOpen() {
		return 0, r.openError
	}
	result, err := r.client.Del(key).Result()
	return result, err
}

// Exists
// 	@Description 返回key是否存在
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
// 	@Return bool 是否存在布尔值
func (r *redisClusterAdapter) Exists(key string) bool {
	if !r.checkOpen() {
		return false
	}
	result, err := r.client.Exists(key).Result()
	if err != nil {
		return false
	}
	return result == 1
}

// ExistsWithErr
// 	@Description 返回key是否存在，并返回错误
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
// 	@Return bool 是否存在布尔值
// 	@Return error 错误
func (r *redisClusterAdapter) ExistsWithErr(key string) (bool, error) {
	if !r.checkOpen() {
		return false, r.openError
	}
	result, err := r.client.Exists(key).Result()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

// Expire
// 	@Description 设置key的过期时间
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param expiration
// 	@Return bool 是否成功设置
// 	@Return error 错误
func (r *redisClusterAdapter) Expire(key string, expiration time.Duration) (bool, error) {
	if !r.checkOpen() {
		return false, r.openError
	}
	result, err := r.client.Expire(key, expiration).Result()
	if err != nil {
		return false, err
	}

	return result, err
}

// TTL 查询过期时间
// TTL
// 	@Description 通过 `TTL key` 查询对应键的过期时间
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
// 	@Return int64 键的过期时间
// 	@Return error 错误
func (r *redisClusterAdapter) TTL(key string) (int64, error) {
	if !r.checkOpen() {
		return 0, r.openError
	}
	if result, err := r.client.TTL(key).Result(); err != nil {
		return 0, err
	} else {
		return int64(result.Seconds()), nil
	}
}

// HGetAll
// 	@Description 通过 `HGETALL key` 命令获取 键下所有字段值
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
// 	@return map[string]string 键下所有字段对应的值
func (r *redisClusterAdapter) HGetAll(key string) map[string]string {
	if !r.checkOpen() {
		return nil
	}
	hashObj := r.client.HGetAll(key)
	hash := hashObj.Val()
	return hash
}

// HGet
// 	@Description 通过 `HGET key field` 命令获取值
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param fields 字段名字
// 	@Return string 键对应字段名字的值
// 	@Return error 错误
func (r *redisClusterAdapter) HGet(key string, fields string) (string, error) {
	if !r.checkOpen() {
		return "", r.openError
	}
	strObj := r.client.HGet(key, fields)
	err := strObj.Err()
	if err != nil && err != redis.Nil {
		return "", err
	}
	if err == redis.Nil {
		return "", nil
	}
	return strObj.Val(), nil
}

// HMGet
// 	@Description 批量获取hash值
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param fileds 给定字段数组
// 	@Return []string 给定字段数组对应的值
func (r *redisClusterAdapter) HMGet(key string, fileds []string) []string {
	if !r.checkOpen() {
		return []string{}
	}
	sliceObj := r.client.HMGet(key, fileds...)
	if err := sliceObj.Err(); err != nil && err != redis.Nil {
		return []string{}
	}
	tmp := sliceObj.Val()
	strSlice := make([]string, 0, len(tmp))
	for _, v := range tmp {
		if v != nil {
			strSlice = append(strSlice, v.(string))
		} else {
			strSlice = append(strSlice, "")
		}
	}
	return strSlice
}

// HMGetMap
// 	@Description 通过 `HGET key field` 批量获取hash值，返回map
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param fields 字段数组
// 	@Return map[string]string 键对应字段名字的值
func (r *redisClusterAdapter) HMGetMap(key string, fields []string) map[string]string {
	if !r.checkOpen() {
		return make(map[string]string)
	}
	if len(fields) == 0 {
		return make(map[string]string)
	}
	sliceObj := r.client.HMGet(key, fields...)
	if err := sliceObj.Err(); err != nil && err != redis.Nil {
		return make(map[string]string)
	}

	tmp := sliceObj.Val()
	hashRet := make(map[string]string, len(tmp))

	var tmpTagID string

	for k, v := range tmp {
		tmpTagID = fields[k]
		if v != nil {
			hashRet[tmpTagID] = v.(string)
		} else {
			hashRet[tmpTagID] = ""
		}
	}
	return hashRet
}

// HMSet
// 	@Description 通过 `HMSET key field value [field value ...]` 命令设置值
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param hash 要设置hash 值
//	@Param expire 过期时间
// 	@Return bool 返回设置成功布尔值
func (r *redisClusterAdapter) HMSet(key string, hash map[string]interface{}, expire time.Duration) bool {
	if !r.checkOpen() {
		return false
	}
	if len(hash) > 0 {
		err := r.client.HMSet(key, hash).Err()
		if err != nil {
			return false
		}
		if expire > 0 {
			r.client.Expire(key, expire)
		}
		return true
	}
	return false
}

// HSet
// 	@Description 通过 `HSET key field value` 命令 设置值
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param field 字段名字
//	@Param value 字段值
// 	@Return bool 返回是否成功设置布尔值
func (r *redisClusterAdapter) HSet(key string, field string, value interface{}) bool {
	if !r.checkOpen() {
		return false
	}
	err := r.client.HSet(key, field, value).Err()
	return err == nil
}

// HDel
// 	@Description 通过 `HDEL key field [field ...]` 删除键
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字数组
//	@Param field 字段数组
// 	@Return bool 返回删除是否成功布尔值
func (r *redisClusterAdapter) HDel(key string, field ...string) bool {
	if !r.checkOpen() {
		return false
	}
	IntObj := r.client.HDel(key, field...)
	err := IntObj.Err()
	return err == nil
}

// Scan
// 	@Description 通过 `SCAN cursor [MATCH pattern] [COUNT count] [TYPE type]` 命令迭代获取键
// 	@Receiver r redisClusterAdapter
//	@Param cursor 游标值
//	@Param match 模式串
//	@Param count 个数
// 	@Return []string 返回的键数组
// 	@Return error 错误
func (r *redisClusterAdapter) Scan(cursor uint64, match string, count int64) ([]string, error) {
	if !r.checkOpen() {
		return nil, r.openError
	}
	result, _, err := r.client.Scan(cursor, match, count).Result()
	return result, err
}

// Type
// 	@Description 通过 `TYPE key` 获取键的类型
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
// 	@Return string 键的类型
// 	@Return error 错误
func (r *redisClusterAdapter) Type(key string) (string, error) {
	if !r.checkOpen() {
		return "", r.openError
	}
	statusObj := r.client.Type(key)
	if err := statusObj.Err(); err != nil {
		return "", err
	}

	return statusObj.Val(), nil
}

// ZRevRange 倒序获取有序集合的部分数据
// ZRevRange
// 	@Description 通过 `ZREVRANGE key start stop ` 倒序获取有序集合的部分数据
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param start 开始索引
//	@Param stop 截止索引
// 	@Return []string 元素数组
// 	@Return error 错误
func (r *redisClusterAdapter) ZRevRange(key string, start, stop int64) ([]string, error) {
	if !r.checkOpen() {
		return []string{}, r.openError
	}
	strSliceObj := r.client.ZRevRange(key, start, stop)
	if err := strSliceObj.Err(); err != nil && err != redis.Nil {
		return []string{}, err
	}
	return strSliceObj.Val(), nil
}

// ZRevRangeWithScores
// 	@Description 通过 `ZREVRANGE key start stop WITHSCORES` 命令获取元素和分数
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param start 开始索引
//	@Param stop 截止索引
// 	@Return []Z 元素和分数的结构数组
// 	@Return error 错误
func (r *redisClusterAdapter) ZRevRangeWithScores(key string, start, stop int64) ([]config.Z, error) {
	if !r.checkOpen() {
		return []redis.Z{}, r.openError
	}
	zSliceObj := r.client.ZRevRangeWithScores(key, start, stop)
	if err := zSliceObj.Err(); err != nil && err != redis.Nil {
		return []redis.Z{}, err
	}
	return zSliceObj.Val(), nil
}

// ZRange
// 	@Description 通过 `ZRANGE key start stop ` 返回有序集合key 中指定范围的元素
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param start 开始索引
//	@Param stop 截止索引
// 	@Return []string 元素数组
// 	@Return error
func (r *redisClusterAdapter) ZRange(key string, start, stop int64) ([]string, error) {
	if !r.checkOpen() {
		return []string{}, r.openError
	}
	strSliceObj := r.client.ZRange(key, start, stop)
	if err := strSliceObj.Err(); err != nil && err != redis.Nil {
		return []string{}, err
	}
	return strSliceObj.Val(), nil
}

// ZRevRank
// 	@Description 通过 `ZREVRANK key member`，返回有序集key中成员member的排名，按照从大到小排名。
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param member 元素
// 	@Return int64 排名
// 	@Return error 错误
func (r *redisClusterAdapter) ZRevRank(key string, member string) (int64, error) {
	if !r.checkOpen() {
		return 0, r.openError
	}
	intObj := r.client.ZRevRank(key, member)
	if err := intObj.Err(); err != nil && err != redis.Nil {
		return 0, err
	}
	return intObj.Val(), nil
}

// ZRevRangeByScore
// 	@Description 通过 `ZRANGEBYSCORE key min max [WITHSCORES] [LIMIT offset count]` 返回在 min和max 之间的所有元素
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param opt min，max ，offset，count 可选参数
// 	@Return []string 元素数组
// 	@Return error 错误
func (r *redisClusterAdapter) ZRevRangeByScore(key string, opt config.ZRangeBy) ([]string, error) {
	if !r.checkOpen() {
		return []string{}, r.openError
	}
	res, err := r.client.ZRevRangeByScore(key, opt).Result()
	if err != nil && err != redis.Nil {
		return []string{}, err
	}

	return res, nil
}

// ZRevRangeByScoreWithScores
// 	@Description 通过 `ZRANGEBYSCORE key min max [WITHSCORES] [LIMIT offset count]` 返回在 min和max 之间的所有元素和分数
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param opt min，max ，offset，count 可选参数
// 	@Return []Z 元素和分数结构数组
// 	@Return error 错误
func (r *redisClusterAdapter) ZRevRangeByScoreWithScores(key string, opt config.ZRangeBy) ([]config.Z, error) {
	if !r.checkOpen() {
		return []redis.Z{}, r.openError
	}
	res, err := r.client.ZRevRangeByScoreWithScores(key, opt).Result()
	if err != nil && err != redis.Nil {
		return []redis.Z{}, err
	}

	return res, nil
}

// ZCard
// 	@Description 通过 `ZCARD key` 命令获取key的有序集元素个数
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
// 	@Return int64 元素个数
// 	@Return error 错误
func (r *redisClusterAdapter) ZCard(key string) (int64, error) {
	if !r.checkOpen() {
		return 0, r.openError
	}
	IntObj := r.client.ZCard(key)
	if err := IntObj.Err(); err != nil {
		return 0, err
	}
	return IntObj.Val(), nil
}

// ZScore
// 	@Description 通过 `ZSCORE key member` 命令 获取有序集key中，成员member的score值。
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param member 元素
// 	@Return float64 分数
// 	@Return error 错误
func (r *redisClusterAdapter) ZScore(key string, member string) (float64, error) {
	if !r.checkOpen() {
		return 0, r.openError
	}
	FloatObj := r.client.ZScore(key, member)
	err := FloatObj.Err()
	if err != nil && err != redis.Nil {
		return 0, err
	}

	return FloatObj.Val(), err
}

// ZAdd
// 	@Description 或多个 member 元素及其 score 值加入到有序集 key 当中
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param members 元素和分数结构数组
// 	@Return int64 添加到有序集合总的成员数量
// 	@Return error 错误
func (r *redisClusterAdapter) ZAdd(key string, members ...config.Z) (int64, error) {
	if !r.checkOpen() {
		return 0, r.openError
	}
	IntObj := r.client.ZAdd(key, members...)
	if err := IntObj.Err(); err != nil && err != redis.Nil {
		return 0, err
	}

	return IntObj.Val(), nil
}

// ZCount
// 	@Description 查询有序集 key 中， score 值在 min 和 max 之间(默认包括 score 值等于 min 或 max )的成员的数量。
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param min 最下值
//	@Param max 最大值
// 	@Return int64 指定分数范围的元素个数
// 	@Return error 错误
func (r *redisClusterAdapter) ZCount(key string, min, max string) (int64, error) {
	if !r.checkOpen() {
		return 0, r.openError
	}
	IntObj := r.client.ZCount(key, min, max)
	if err := IntObj.Err(); err != nil && err != redis.Nil {
		return 0, err
	}

	return IntObj.Val(), nil
}

// HIncrBy
// 	@Description 增加 key 指定的哈希集中指定字段的数值
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param field 键某个字段
//	@Param incr 增量
// 	@Return int64 执行后该字段额值
func (r *redisClusterAdapter) HIncrBy(key string, field string, incr int) int64 {
	if !r.checkOpen() {
		return 0
	}
	result, err := r.client.HIncrBy(key, field, int64(incr)).Result()
	if err != nil {
		return 0
	}
	return result
}

// HIncrByWithErr
// 	@Description 增加 key 指定的哈希集中指定字段的数值，并返回错误
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param field 键某个字段
//	@Param incr 增量
// 	@Return int64 执行后该字段额值
// 	@Return error 执行后的错误
func (r *redisClusterAdapter) HIncrByWithErr(key string, field string, incr int) (int64, error) {
	if !r.checkOpen() {
		return 0, r.openError
	}
	return r.client.HIncrBy(key, field, int64(incr)).Result()
}

// LPush
// 	@Description 将一个或多个值 value 插入到列表 key 的表头
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param values 值数组
// 	@return int64 执行后list 的长度
// 	@return error 错误
func (r *redisClusterAdapter) LPush(key string, values ...interface{}) (int64, error) {
	if !r.checkOpen() {
		return 0, r.openError
	}
	IntObj := r.client.LPush(key, values...)
	if err := IntObj.Err(); err != nil {
		return 0, err
	}

	return IntObj.Val(), nil
}

// RPush
// 	@Description 一个或多个值 value 插入到列表 key 的表尾(最右边)。
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param values 值数组
// 	@Return int64 执行后list 的长度
// 	@Return error 错误
func (r *redisClusterAdapter) RPush(key string, values ...interface{}) (int64, error) {
	if !r.checkOpen() {
		return 0, r.openError
	}
	IntObj := r.client.RPush(key, values...)
	if err := IntObj.Err(); err != nil {
		return 0, err
	}

	return IntObj.Val(), nil
}

// RPop
// 	@Description 移除并返回存于 key 的 list 的最后一个元素。
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
// 	@Return string 执行后最后一个元素的值
// 	@Return error 错误
func (r *redisClusterAdapter) RPop(key string) (string, error) {
	if !r.checkOpen() {
		return "", r.openError
	}
	strObj := r.client.RPop(key)
	if err := strObj.Err(); err != nil {
		return "", err
	}

	return strObj.Val(), nil
}

// LRange
// 	@Description 获取列表指定范围内的元素
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param start 开始下标
//	@Param stop 截止下标
// 	@Return []string 指定范围内的元素
// 	@Return error 错误
func (r *redisClusterAdapter) LRange(key string, start, stop int64) ([]string, error) {
	if !r.checkOpen() {
		return []string{}, r.openError
	}
	result, err := r.client.LRange(key, start, stop).Result()
	if err != nil {
		return []string{}, err
	}

	return result, nil
}

// LLen
// 	@Description 获取list 的长度
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
// 	@Return int64 list 长度
func (r *redisClusterAdapter) LLen(key string) int64 {
	if !r.checkOpen() {
		return 0
	}
	IntObj := r.client.LLen(key)
	if err := IntObj.Err(); err != nil {
		return 0
	}

	return IntObj.Val()
}

// LLenWithErr
// 	@Description 获取list 的长度，并返回错误
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
// 	@Return int64 list 长度
// 	@Return error 错误
func (r *redisClusterAdapter) LLenWithErr(key string) (int64, error) {
	if !r.checkOpen() {
		return 0, r.openError
	}
	ret, err := r.client.LLen(key).Result()
	return ret, err
}

// LRem
// 	@Description 存于 key 的列表里移除前 count 次出现的值为 value 的元素
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param count
//   count > 0: 从头往尾移除值为 value 的元素。
// 	 count < 0: 从尾往头移除值为 value 的元素。
//	 count = 0: 移除所有值为 value 的元素。
//	@Param value 要移除的值
// 	@Return int64 被移除的元素个数
func (r *redisClusterAdapter) LRem(key string, count int64, value interface{}) int64 {
	if !r.checkOpen() {
		return 0
	}
	IntObj := r.client.LRem(key, count, value)
	if err := IntObj.Err(); err != nil {
		return 0
	}

	return IntObj.Val()
}

// LIndex
// 	@Description 返回指定索引的键对应的元素
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param idx 索引
// 	@Return string 对应的元素数
// 	@Return error
func (r *redisClusterAdapter) LIndex(key string, idx int64) (string, error) {
	if !r.checkOpen() {
		return "", r.openError
	}
	ret, err := r.client.LIndex(key, idx).Result()
	return ret, err
}

// LTrim
// 	@Description 只包含指定范围的指定元素
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param start 开始下标
//	@Param stop 结束下标
// 	@Return string
// 	@Return error 错误
func (r *redisClusterAdapter) LTrim(key string, start, stop int64) (string, error) {
	if !r.checkOpen() {
		return "", r.openError
	}
	ret, err := r.client.LTrim(key, start, stop).Result()
	return ret, err
}

// ZRemRangeByRank
// 	@Description 移除有序集key中，指定排名(rank)区间内的所有成员
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param start 开始下标 这些索引也可是负数，表示位移从最高分处开始数
//	@Param stop 结束下标 这些索引也可是负数，表示位移从最高分处开始数
// 	@Return int64 被移除成员的数量
// 	@Return error 错误
func (r *redisClusterAdapter) ZRemRangeByRank(key string, start, stop int64) (int64, error) {
	if !r.checkOpen() {
		return 0, r.openError
	}
	result, err := r.client.ZRemRangeByRank(key, start, stop).Result()
	if err != nil {
		return 0, err
	}

	return result, nil
}

// ZRemRangeByScore
// 	@Description 移除有序集中，指定分数（score）区间内的所有成员。
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param min 分数最小值
//	@Param max 分数最大值
// 	@Return int64 删除的元素的个数
// 	@Return error 错误
func (r *redisClusterAdapter) ZRemRangeByScore(key string, min, max string) (int64, error) {
	if !r.checkOpen() {
		return 0, r.openError
	}
	result, err := r.client.ZRemRangeByScore(key, min, max).Result()
	if err != nil {
		return 0, err
	}

	return result, nil
}

// ZRem
// 	@Description 有序集合移除元素
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param members 元素数组
// 	@Return int64 删除的成员个数
// 	@Return error 错误
func (r *redisClusterAdapter) ZRem(key string, members ...interface{}) (int64, error) {
	if !r.checkOpen() {
		return 0, r.openError
	}
	result, err := r.client.ZRem(key, members...).Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}

// SAdd
// 	@Description 添加一个或多个指定的member元素到集合的 key中
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param member 元素数组
// 	@Return int64 成功添加到集合里元素的数量
// 	@Return error 错误
func (r *redisClusterAdapter) SAdd(key string, member ...interface{}) (int64, error) {
	if !r.checkOpen() {
		return 0, r.openError
	}
	intObj := r.client.SAdd(key, member...)
	if err := intObj.Err(); err != nil {
		return 0, err
	}
	return intObj.Val(), nil
}

// SMembers
// 	@Description 返回set的全部成员
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
// 	@Return []string 所有元素数组
// 	@Return error 错误
func (r *redisClusterAdapter) SMembers(key string) ([]string, error) {
	if !r.checkOpen() {
		return []string{}, r.openError
	}
	strSliceObj := r.client.SMembers(key)
	if err := strSliceObj.Err(); err != nil {
		return []string{}, err
	}
	return strSliceObj.Val(), nil
}

// SIsMember
// 	@Description 成员元素是否在集合中
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param member  成员元素
// 	@Return bool 是否存在布尔值
// 	@Return error
func (r *redisClusterAdapter) SIsMember(key string, member interface{}) (bool, error) {
	if !r.checkOpen() {
		return false, r.openError
	}
	boolObj := r.client.SIsMember(key, member)
	if err := boolObj.Err(); err != nil {
		return false, err
	}
	return boolObj.Val(), nil
}

// HKeys
// 	@Description 哈希集中所有字段的名字
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
// 	@Return []string 字段列表
func (r *redisClusterAdapter) HKeys(key string) []string {
	if !r.checkOpen() {
		return []string{}
	}
	strObj := r.client.HKeys(key)
	if err := strObj.Err(); err != nil && err != redis.Nil {
		return []string{}
	}
	return strObj.Val()
}

// HLen
// 	@Description 查询包含的字段的数量
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
// 	@Return int64 字段的数量
func (r *redisClusterAdapter) HLen(key string) int64 {
	if !r.checkOpen() {
		return 0
	}
	intObj := r.client.HLen(key)
	if err := intObj.Err(); err != nil && err != redis.Nil {
		return 0
	}
	return intObj.Val()
}

// GeoAdd
// 	@Description 通过 `GEOADD key [NX|XX] [CH] longitude latitude member ` 写入地理位置
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param location 地理位置
// 	@Return int64 添加到sorted set元素的数目
// 	@Return error 错误
func (r *redisClusterAdapter) GeoAdd(key string, location *config.GeoLocation) (int64, error) {
	if !r.checkOpen() {
		return 0, r.openError
	}
	res, err := r.client.GeoAdd(key, location).Result()
	if err != nil {
		return 0, err
	}

	return res, nil
}

// GeoRadius 根据经纬度查询列表
// GeoRadius
// 	@Description 通过 `GEORADIUS key longitude latitude radiu` 命令 获取位置信息
// 	@Receiver r redisClusterAdapter
//	@Param key 键名字
//	@Param longitude 经度
//	@Param latitude 维度
//	@Param query GeoRadiusQuery
// 	@Return []GeoLocation 地理位置数组
// 	@Return error 错误
func (r *redisClusterAdapter) GeoRadius(key string, longitude, latitude float64, query *config.GeoRadiusQuery) ([]config.GeoLocation, error) {
	if !r.checkOpen() {
		return []redis.GeoLocation{}, r.openError
	}
	res, err := r.client.GeoRadius(key, longitude, latitude, query).Result()
	if err != nil {
		return []redis.GeoLocation{}, err
	}

	return res, nil
}

// Open
// 	@Description 初始化redis 客户端
// 	@Receiver r redisClusterAdapter
//	@Param ctx 上下文
// 	@Return config.ICache 实例化后的redis 实例
func (r *redisClusterAdapter) Open(ctx context.Context) config.ICache {
	r.openOnce.Do(func() {
		client := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        r.config.Addrs,
			MaxRedirects: r.config.MaxRetries,
			ReadOnly:     r.config.ReadOnly,
			Password:     r.config.Password,
			MaxRetries:   r.config.MaxRetries,
			DialTimeout:  r.config.DialTimeout,
			ReadTimeout:  r.config.ReadTimeout,
			WriteTimeout: r.config.WriteTimeout,
			PoolSize:     r.config.PoolSize,
			MinIdleConns: r.config.MinIdleConns,
			IdleTimeout:  r.config.IdleTimeout,
		})
		if err := client.Ping().Err(); err != nil {
			switch r.config.OnDialError {
			case "panic":
				klog.KuaigoLogger.WithContext(ctx).Panicf("start cluster redis err:%v", err)
			default:
				klog.KuaigoLogger.WithContext(ctx).Errorf("start cluster redis err:%v", err)
			}
		}
		r.client = client
		r.hasOpen = true
	})
	return r
}

// GetClient
// 	@Description 获取缓存操作句柄
// 	@Receiver redisClusterAdapter
// 	@Return ICache 缓存操作句柄
func (r *redisClusterAdapter) GetClient() config.ICache {
	return r
}

// Close
// 	@Description 关闭连接
// 	@Receiver redisClusterAdapter
// 	@Return 关闭结果
func (r *redisClusterAdapter) Close() (err error) {
	if r.checkOpen() && r.client != nil {
		if c, ok := r.client.(*redis.ClusterClient); ok {
			return c.Close()
		}
	}
	return nil
}

// checkOpen
// 	@Description 检查连接是否打开
// 	@Receiver redisClusterAdapter
// 	@Return 打开结果
func (r *redisClusterAdapter) checkOpen() bool {
	if !r.hasOpen {
		r.Open(r.ctx)
	}
	return true
}
