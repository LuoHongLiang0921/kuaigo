// @Description redis 常用命令

package redis

import (
	"time"

	"github.com/go-redis/redis"
)

// Get 从redis获取string
// 	@Description 通过 `get key` 命令获取key字符串值
// 	@Receiver r Redis
//	@Param key 键名字
// 	@return string key 对应的字符串值
func (r *Redis) Get(key string) string {
	var mes string
	strObj := r.Client.Get(key)
	if err := strObj.Err(); err != nil {
		mes = ""
	} else {
		mes = strObj.Val()
	}
	return mes
}

// GetRaw
// 	@Description 通过 `get key` 命令获取key字节数组值
// 	@Receiver r Redis
//	@Param key 键名字
// 	@return []byte 字节数组
// 	@return error 错误
func (r *Redis) GetRaw(key string) ([]byte, error) {
	c, err := r.Client.Get(key).Bytes()
	if err != nil && err != redis.Nil {
		return []byte{}, err
	}
	return c, nil
}

// MGet
// 	@Description 通过 `get key` 命令获取key字节数组值
// 	@Receiver r Redis
//	@Param keys 键名字数组
// 	@return []string
// 	@return error
func (r *Redis) MGet(keys ...string) ([]string, error) {
	sliceObj := r.Client.MGet(keys...)
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
// 	@Receiver r Redis
//	@Param keys 键名字数组
// 	@return []interface{} 键对应的值数组
// 	@return error 错误
func (r *Redis) MGets(keys []string) ([]interface{}, error) {
	ret, err := r.Client.MGet(keys...).Result()
	if err != nil && err != redis.Nil {
		return []interface{}{}, err
	}
	return ret, nil
}

// Set 设置redis的string
// 	@Description 通过 `SET key value [EX seconds]` 设置键对应的值
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param value 值
//	@Param expire 过期时间
// 	@return bool 是否设置成功布尔值
func (r *Redis) Set(key string, value interface{}, expire time.Duration) bool {
	err := r.Client.Set(key, value, expire).Err()
	return err == nil
}

// HGetAll
// 	@Description 通过 `HGETALL key` 命令获取 键下所有字段值
// 	@Receiver r Redis
//	@Param key 键名字
// 	@return map[string]string 键下所有字段对应的值
func (r *Redis) HGetAll(key string) map[string]string {
	hashObj := r.Client.HGetAll(key)
	hash := hashObj.Val()
	return hash
}

// HGet
// 	@Description 通过 `HGET key field` 命令获取值
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param fields 字段名字
// 	@Return string 键对应字段名字的值
// 	@Return error 错误
func (r *Redis) HGet(key string, fields string) (string, error) {
	strObj := r.Client.HGet(key, fields)
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
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param fileds 给定字段数组
// 	@Return []string 给定字段数组对应的值
func (r *Redis) HMGet(key string, fileds []string) []string {
	sliceObj := r.Client.HMGet(key, fileds...)
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
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param fields 字段数组
// 	@Return map[string]string 键对应字段名字的值
func (r *Redis) HMGetMap(key string, fields []string) map[string]string {
	if len(fields) == 0 {
		return make(map[string]string)
	}
	sliceObj := r.Client.HMGet(key, fields...)
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
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param hash 要设置hash 值
//	@Param expire 过期时间
// 	@Return bool 返回设置成功布尔值
func (r *Redis) HMSet(key string, hash map[string]interface{}, expire time.Duration) bool {
	if len(hash) > 0 {
		err := r.Client.HMSet(key, hash).Err()
		if err != nil {
			return false
		}
		if expire > 0 {
			r.Client.Expire(key, expire)
		}
		return true
	}
	return false
}

// HSet
// 	@Description 通过 `HSET key field value` 命令 设置值
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param field 字段名字
//	@Param value 字段值
// 	@Return bool 返回是否成功设置布尔值
func (r *Redis) HSet(key string, field string, value interface{}) bool {
	err := r.Client.HSet(key, field, value).Err()
	return err == nil
}

// HDel
// 	@Description 通过 `HDEL key field [field ...]` 删除键
// 	@Receiver r Redis
//	@Param key 键名字数组
//	@Param field 字段数组
// 	@Return bool 返回删除是否成功布尔值
func (r *Redis) HDel(key string, field ...string) bool {
	IntObj := r.Client.HDel(key, field...)
	err := IntObj.Err()
	return err == nil
}

// SetWithErr
// 	@Description 通过 `SET key value [EX seconds]` 设置键对应的值，并返回错误值
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param value 键对应的值
//	@Param expire 过期时间
// 	@Return error
func (r *Redis) SetWithErr(key string, value interface{}, expire time.Duration) error {
	err := r.Client.Set(key, value, expire).Err()
	return err
}

// SetNx
// 	@Description 通过 `SETNX key value` 命令设置键不存在的值
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param value 键值
//	@Param expiration 过期时间
// 	@return bool 是否设置成功
func (r *Redis) SetNx(key string, value interface{}, expiration time.Duration) bool {

	result, err := r.Client.SetNX(key, value, expiration).Result()

	if err != nil {
		return false
	}

	return result
}

// SetNxWithErr
// 	@Description 通过 `SETNX key value` 命令设置键不存在的值，并返回错误
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param value 键值
//	@Param expiration 过期时间
// 	@return bool 是否设置成功
// 	@return error 错误
func (r *Redis) SetNxWithErr(key string, value interface{}, expiration time.Duration) (bool, error) {
	result, err := r.Client.SetNX(key, value, expiration).Result()
	return result, err
}

// Incr redis自增
// 	@Description 通过 `INCR key` 自增键值
// 	@Receiver r Redis
//	@Param key 键名字
// 	@return bool
func (r *Redis) Incr(key string) bool {
	err := r.Client.Incr(key).Err()
	return err == nil
}

// IncrWithErr ...
// IncrWithErr
// 	@Description 通过 `INCR key` 自增键值
// 	@Receiver r Redis
//	@Param key 键名字
// 	@Return int64 自增后的值
// 	@Return error
func (r *Redis) IncrWithErr(key string) (int64, error) {
	ret, err := r.Client.Incr(key).Result()
	return ret, err
}

// IncrBy 将 key 所储存的值加上增量 increment
// IncrBy
// 	@Description 通过 `INCRBY key increment` 命令增加增量值
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param increment 增量值
// 	@Return int64 加上增量后的值
// 	@Return error 错误
func (r *Redis) IncrBy(key string, increment int64) (int64, error) {
	intObj := r.Client.IncrBy(key, increment)
	if err := intObj.Err(); err != nil {
		return 0, err
	}
	return intObj.Val(), nil
}

// Decr
// 	@Description 通过 `DECR key` 自减键对应的值
// 	@Receiver r Redis
//	@Param key 键名字
// 	@Return bool
func (r *Redis) Decr(key string) bool {
	err := r.Client.Decr(key).Err()
	return err == nil
}

// Scan
// 	@Description 通过 `SCAN cursor [MATCH pattern] [COUNT count] [TYPE type]` 命令迭代获取键
// 	@Receiver r Redis
//	@Param cursor 游标值
//	@Param match 模式串
//	@Param count 个数
// 	@Return []string 返回的键数组
// 	@Return error 错误
func (r *Redis) Scan(cursor uint64, match string, count int64) ([]string, error) {
	result, _, err := r.Client.Scan(cursor, match, count).Result()
	return result, err
}

// Type
// 	@Description 通过 `TYPE key` 获取键的类型
// 	@Receiver r Redis
//	@Param key 键名字
// 	@Return string 键的类型
// 	@Return error 错误
func (r *Redis) Type(key string) (string, error) {
	statusObj := r.Client.Type(key)
	if err := statusObj.Err(); err != nil {
		return "", err
	}

	return statusObj.Val(), nil
}

// ZRevRange 倒序获取有序集合的部分数据
// ZRevRange
// 	@Description 通过 `ZREVRANGE key start stop ` 倒序获取有序集合的部分数据
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param start 开始索引
//	@Param stop 截止索引
// 	@Return []string 元素数组
// 	@Return error 错误
func (r *Redis) ZRevRange(key string, start, stop int64) ([]string, error) {
	strSliceObj := r.Client.ZRevRange(key, start, stop)
	if err := strSliceObj.Err(); err != nil && err != redis.Nil {
		return []string{}, err
	}
	return strSliceObj.Val(), nil
}

// ZRevRangeWithScores
// 	@Description 通过 `ZREVRANGE key start stop WITHSCORES` 命令获取元素和分数
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param start 开始索引
//	@Param stop 截止索引
// 	@Return []Z 元素和分数的结构数组
// 	@Return error 错误
func (r *Redis) ZRevRangeWithScores(key string, start, stop int64) ([]Z, error) {
	zSliceObj := r.Client.ZRevRangeWithScores(key, start, stop)
	if err := zSliceObj.Err(); err != nil && err != redis.Nil {
		return []redis.Z{}, err
	}
	return zSliceObj.Val(), nil
}

// ZRange
// 	@Description 通过 `ZRANGE key start stop ` 返回有序集合key 中指定范围的元素
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param start 开始索引
//	@Param stop 截止索引
// 	@Return []string 元素数组
// 	@Return error
func (r *Redis) ZRange(key string, start, stop int64) ([]string, error) {
	strSliceObj := r.Client.ZRange(key, start, stop)
	if err := strSliceObj.Err(); err != nil && err != redis.Nil {
		return []string{}, err
	}
	return strSliceObj.Val(), nil
}

// ZRevRank
// 	@Description 通过 `ZREVRANK key member`，返回有序集key中成员member的排名，按照从大到小排名。
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param member 元素
// 	@Return int64 排名
// 	@Return error 错误
func (r *Redis) ZRevRank(key string, member string) (int64, error) {
	intObj := r.Client.ZRevRank(key, member)
	if err := intObj.Err(); err != nil && err != redis.Nil {
		return 0, err
	}
	return intObj.Val(), nil
}

// ZRevRangeByScore
// 	@Description 通过 `ZRANGEBYSCORE key min max [WITHSCORES] [LIMIT offset count]` 返回在 min和max 之间的所有元素
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param opt min，max ，offset，count 可选参数
// 	@Return []string 元素数组
// 	@Return error 错误
func (r *Redis) ZRevRangeByScore(key string, opt ZRangeBy) ([]string, error) {
	res, err := r.Client.ZRevRangeByScore(key, opt).Result()
	if err != nil && err != redis.Nil {
		return []string{}, err
	}

	return res, nil
}

// ZRevRangeByScoreWithScores
// 	@Description 通过 `ZRANGEBYSCORE key min max [WITHSCORES] [LIMIT offset count]` 返回在 min和max 之间的所有元素和分数
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param opt min，max ，offset，count 可选参数
// 	@Return []Z 元素和分数结构数组
// 	@Return error 错误
func (r *Redis) ZRevRangeByScoreWithScores(key string, opt ZRangeBy) ([]Z, error) {
	res, err := r.Client.ZRevRangeByScoreWithScores(key, opt).Result()
	if err != nil && err != redis.Nil {
		return []redis.Z{}, err
	}

	return res, nil
}

// ZCard
// 	@Description 通过 `ZCARD key` 命令获取key的有序集元素个数
// 	@Receiver r Redis
//	@Param key 键名字
// 	@Return int64 元素个数
// 	@Return error 错误
func (r *Redis) ZCard(key string) (int64, error) {
	IntObj := r.Client.ZCard(key)
	if err := IntObj.Err(); err != nil {
		return 0, err
	}
	return IntObj.Val(), nil
}

// ZScore
// 	@Description 通过 `ZSCORE key member` 命令 获取有序集key中，成员member的score值。
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param member 元素
// 	@Return float64 分数
// 	@Return error 错误
func (r *Redis) ZScore(key string, member string) (float64, error) {
	FloatObj := r.Client.ZScore(key, member)
	err := FloatObj.Err()
	if err != nil && err != redis.Nil {
		return 0, err
	}

	return FloatObj.Val(), err
}

// ZAdd
// 	@Description 或多个 member 元素及其 score 值加入到有序集 key 当中
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param members 元素和分数结构数组
// 	@Return int64 添加到有序集合总的成员数量
// 	@Return error 错误
func (r *Redis) ZAdd(key string, members ...Z) (int64, error) {
	IntObj := r.Client.ZAdd(key, members...)
	if err := IntObj.Err(); err != nil && err != redis.Nil {
		return 0, err
	}

	return IntObj.Val(), nil
}

// ZCount
// 	@Description 查询有序集 key 中， score 值在 min 和 max 之间(默认包括 score 值等于 min 或 max )的成员的数量。
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param min 最下值
//	@Param max 最大值
// 	@Return int64 指定分数范围的元素个数
// 	@Return error 错误
func (r *Redis) ZCount(key string, min, max string) (int64, error) {
	IntObj := r.Client.ZCount(key, min, max)
	if err := IntObj.Err(); err != nil && err != redis.Nil {
		return 0, err
	}

	return IntObj.Val(), nil
}

// Del
// 	@Description 删除key
// 	@Receiver r Redis
//	@Param key 键名字
// 	@Return int64 删除key成功的个数
func (r *Redis) Del(key string) int64 {
	result, err := r.Client.Del(key).Result()
	if err != nil {
		return 0
	}
	return result
}

// DelWithErr
// 	@Description 删除key，并返回错误
// 	@Receiver r Redis
//	@Param key 键名字
// 	@Return int64 删除key成功的个数
// 	@Return error 错误
func (r *Redis) DelWithErr(key string) (int64, error) {
	result, err := r.Client.Del(key).Result()
	return result, err
}

// HIncrBy
// 	@Description 增加 key 指定的哈希集中指定字段的数值
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param field 键某个字段
//	@Param incr 增量
// 	@Return int64 执行后该字段额值
func (r *Redis) HIncrBy(key string, field string, incr int) int64 {
	result, err := r.Client.HIncrBy(key, field, int64(incr)).Result()
	if err != nil {
		return 0
	}
	return result
}

// HIncrByWithErr
// 	@Description 增加 key 指定的哈希集中指定字段的数值，并返回错误
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param field 键某个字段
//	@Param incr 增量
// 	@Return int64 执行后该字段额值
// 	@Return error 执行后的错误
func (r *Redis) HIncrByWithErr(key string, field string, incr int) (int64, error) {
	return r.Client.HIncrBy(key, field, int64(incr)).Result()
}

// Exists
// 	@Description 返回key是否存在
// 	@Receiver r Redis
//	@Param key 键名字
// 	@Return bool 是否存在布尔值
func (r *Redis) Exists(key string) bool {
	result, err := r.Client.Exists(key).Result()
	if err != nil {
		return false
	}
	return result == 1
}

// ExistsWithErr
// 	@Description 返回key是否存在，并返回错误
// 	@Receiver r Redis
//	@Param key 键名字
// 	@Return bool 是否存在布尔值
// 	@Return error 错误
func (r *Redis) ExistsWithErr(key string) (bool, error) {
	result, err := r.Client.Exists(key).Result()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

// LPush
// 	@Description 将一个或多个值 value 插入到列表 key 的表头
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param values 值数组
// 	@return int64 执行后list 的长度
// 	@return error 错误
func (r *Redis) LPush(key string, values ...interface{}) (int64, error) {
	IntObj := r.Client.LPush(key, values...)
	if err := IntObj.Err(); err != nil {
		return 0, err
	}

	return IntObj.Val(), nil
}

// RPush
// 	@Description 一个或多个值 value 插入到列表 key 的表尾(最右边)。
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param values 值数组
// 	@Return int64 执行后list 的长度
// 	@Return error 错误
func (r *Redis) RPush(key string, values ...interface{}) (int64, error) {
	IntObj := r.Client.RPush(key, values...)
	if err := IntObj.Err(); err != nil {
		return 0, err
	}

	return IntObj.Val(), nil
}

// RPop
// 	@Description 移除并返回存于 key 的 list 的最后一个元素。
// 	@Receiver r Redis
//	@Param key 键名字
// 	@Return string 执行后最后一个元素的值
// 	@Return error 错误
func (r *Redis) RPop(key string) (string, error) {
	strObj := r.Client.RPop(key)
	if err := strObj.Err(); err != nil {
		return "", err
	}

	return strObj.Val(), nil
}

// LRange
// 	@Description 获取列表指定范围内的元素
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param start 开始下标
//	@Param stop 截止下标
// 	@Return []string 指定范围内的元素
// 	@Return error 错误
func (r *Redis) LRange(key string, start, stop int64) ([]string, error) {
	result, err := r.Client.LRange(key, start, stop).Result()
	if err != nil {
		return []string{}, err
	}

	return result, nil
}

// LLen
// 	@Description 获取list 的长度
// 	@Receiver r Redis
//	@Param key 键名字
// 	@Return int64 list 长度
func (r *Redis) LLen(key string) int64 {
	IntObj := r.Client.LLen(key)
	if err := IntObj.Err(); err != nil {
		return 0
	}

	return IntObj.Val()
}

// LLenWithErr
// 	@Description 获取list 的长度，并返回错误
// 	@Receiver r Redis
//	@Param key 键名字
// 	@Return int64 list 长度
// 	@Return error 错误
func (r *Redis) LLenWithErr(key string) (int64, error) {
	ret, err := r.Client.LLen(key).Result()
	return ret, err
}

// LRem
// 	@Description 存于 key 的列表里移除前 count 次出现的值为 value 的元素
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param count
//   count > 0: 从头往尾移除值为 value 的元素。
// 	 count < 0: 从尾往头移除值为 value 的元素。
//	 count = 0: 移除所有值为 value 的元素。
//	@Param value 要移除的值
// 	@Return int64 被移除的元素个数
func (r *Redis) LRem(key string, count int64, value interface{}) int64 {
	IntObj := r.Client.LRem(key, count, value)
	if err := IntObj.Err(); err != nil {
		return 0
	}

	return IntObj.Val()
}

// LIndex
// 	@Description 返回指定索引的键对应的元素
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param idx 索引
// 	@Return string 对应的元素数
// 	@Return error
func (r *Redis) LIndex(key string, idx int64) (string, error) {
	ret, err := r.Client.LIndex(key, idx).Result()
	return ret, err
}

// LTrim
// 	@Description 只包含指定范围的指定元素
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param start 开始下标
//	@Param stop 结束下标
// 	@Return string
// 	@Return error 错误
func (r *Redis) LTrim(key string, start, stop int64) (string, error) {
	ret, err := r.Client.LTrim(key, start, stop).Result()
	return ret, err
}

// ZRemRangeByRank
// 	@Description 移除有序集key中，指定排名(rank)区间内的所有成员
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param start 开始下标 这些索引也可是负数，表示位移从最高分处开始数
//	@Param stop 结束下标 这些索引也可是负数，表示位移从最高分处开始数
// 	@Return int64 被移除成员的数量
// 	@Return error 错误
func (r *Redis) ZRemRangeByRank(key string, start, stop int64) (int64, error) {
	result, err := r.Client.ZRemRangeByRank(key, start, stop).Result()
	if err != nil {
		return 0, err
	}

	return result, nil
}

// ZRemRangeByScore
// 	@Description 移除有序集中，指定分数（score）区间内的所有成员。
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param min 分数最小值
//	@Param max 分数最大值
// 	@Return int64 删除的元素的个数
// 	@Return error 错误
func (r *Redis) ZRemRangeByScore(key string, min, max string) (int64, error) {
	result, err := r.Client.ZRemRangeByScore(key, min, max).Result()
	if err != nil {
		return 0, err
	}

	return result, nil
}

// Expire
// 	@Description 设置key的过期时间
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param expiration
// 	@Return bool 是否成功设置
// 	@Return error 错误
func (r *Redis) Expire(key string, expiration time.Duration) (bool, error) {
	result, err := r.Client.Expire(key, expiration).Result()
	if err != nil {
		return false, err
	}

	return result, err
}

// ZRem
// 	@Description 有序集合移除元素
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param members 元素数组
// 	@Return int64 删除的成员个数
// 	@Return error 错误
func (r *Redis) ZRem(key string, members ...interface{}) (int64, error) {
	result, err := r.Client.ZRem(key, members...).Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}

// SAdd
// 	@Description 添加一个或多个指定的member元素到集合的 key中
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param member 元素数组
// 	@Return int64 成功添加到集合里元素的数量
// 	@Return error 错误
func (r *Redis) SAdd(key string, member ...interface{}) (int64, error) {
	intObj := r.Client.SAdd(key, member...)
	if err := intObj.Err(); err != nil {
		return 0, err
	}
	return intObj.Val(), nil
}

// SMembers
// 	@Description 返回set的全部成员
// 	@Receiver r Redis
//	@Param key 键名字
// 	@Return []string 所有元素数组
// 	@Return error 错误
func (r *Redis) SMembers(key string) ([]string, error) {
	strSliceObj := r.Client.SMembers(key)
	if err := strSliceObj.Err(); err != nil {
		return []string{}, err
	}
	return strSliceObj.Val(), nil
}

// SIsMember
// 	@Description 成员元素是否在集合中
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param member  成员元素
// 	@Return bool 是否存在布尔值
// 	@Return error
func (r *Redis) SIsMember(key string, member interface{}) (bool, error) {
	boolObj := r.Client.SIsMember(key, member)
	if err := boolObj.Err(); err != nil {
		return false, err
	}
	return boolObj.Val(), nil
}

// HKeys
// 	@Description 哈希集中所有字段的名字
// 	@Receiver r Redis
//	@Param key 键名字
// 	@Return []string 字段列表
func (r *Redis) HKeys(key string) []string {
	strObj := r.Client.HKeys(key)
	if err := strObj.Err(); err != nil && err != redis.Nil {
		return []string{}
	}
	return strObj.Val()
}

// HLen
// 	@Description 查询包含的字段的数量
// 	@Receiver r Redis
//	@Param key 键名字
// 	@Return int64 字段的数量
func (r *Redis) HLen(key string) int64 {
	intObj := r.Client.HLen(key)
	if err := intObj.Err(); err != nil && err != redis.Nil {
		return 0
	}
	return intObj.Val()
}

// GeoAdd
// 	@Description 通过 `GEOADD key [NX|XX] [CH] longitude latitude member ` 写入地理位置
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param location 地理位置
// 	@Return int64 添加到sorted set元素的数目
// 	@Return error 错误
func (r *Redis) GeoAdd(key string, location *GeoLocation) (int64, error) {
	res, err := r.Client.GeoAdd(key, location).Result()
	if err != nil {
		return 0, err
	}

	return res, nil
}

// GeoRadius 根据经纬度查询列表
// GeoRadius
// 	@Description 通过 `GEORADIUS key longitude latitude radiu` 命令 获取位置信息
// 	@Receiver r Redis
//	@Param key 键名字
//	@Param longitude 经度
//	@Param latitude 维度
//	@Param query GeoRadiusQuery
// 	@Return []GeoLocation 地理位置数组
// 	@Return error 错误
func (r *Redis) GeoRadius(key string, longitude, latitude float64, query *GeoRadiusQuery) ([]GeoLocation, error) {
	res, err := r.Client.GeoRadius(key, longitude, latitude, query).Result()
	if err != nil {
		return []redis.GeoLocation{}, err
	}

	return res, nil
}

// TTL 查询过期时间
// TTL
// 	@Description 通过 `TTL key` 查询对应键的过期时间
// 	@Receiver r Redis
//	@Param key 键名字
// 	@Return int64 键的过期时间
// 	@Return error 错误
func (r *Redis) TTL(key string) (int64, error) {
	if result, err := r.Client.TTL(key).Result(); err != nil {
		return 0, err
	} else {
		return int64(result.Seconds()), nil
	}
}

// Close
// 	@Description 关闭redis，释放连接资源
// 	@Receiver r Redis
// 	@Return err 错误
func (r *Redis) Close() (err error) {
	err = nil
	if r.Client != nil {
		if r.Cluster() != nil {
			err = r.Cluster().Close()
		}

		if r.Stub() != nil {
			err = r.Stub().Close()
		}
	}
	return err
}
