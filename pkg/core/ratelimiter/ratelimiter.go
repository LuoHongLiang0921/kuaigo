// @Description 基础限流库

package ratelimiter

import (
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/server/kgin"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/storage/redis"
	"github.com/LuoHongLiang0921/kuaigo/pkg/ecode"
	"github.com/LuoHongLiang0921/kuaigo/pkg/errs"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"time"
)

// RateLimiter is a antispam instance.
type RateLimiter struct {
	cfg   *Config
	conf  *RLConfig
	Redis *redis.Redis
}

// RLConfig antispam config.
type RLConfig struct {
	IP      string // ip
	Path    string // 接口地址
	UDID    string // 设备id
	On      bool   // switch on/off 开关
	Minutes int    // every N Minutes allow N requests.
	N       int    // one unit allow N requests.
	Hour    int    // every N hour allow M requests.
	M       int    // one winodw allow M requests.

	UniqueIDs string // unique ids: _udid,_ip
}

const (
	defMinutes = 1
	defN       = 60 // 每分钟上限不超过60
	defHour    = 1
	defM       = 100 // 每小时下限不超过100
	defUDid    = "_udid"
	defIp      = "_ip"
)
const (
	prefixMinuteKey = "m_%s_%s_%d"
	prefixHourKey   = "h_%s_%s_%d"
)

// Execute
//  @Description: 设置执行限流操作
//  @Receiver s
//  @Param name
//  @Param did
//  @Param resource
//  @Param ip
//  @Return *RateLimiter
func (s *RateLimiter) Execute(c *kgin.TContext, name, udid, resource, ip string) *RateLimiter {
	cfg := s.cfg.Rule[name]
	if cfg == nil {
		klog.Error("RateLimiter config nil", klog.Any("name", name))
	}
	cfg.UDID = udid
	cfg.Path = resource
	cfg.IP = ip
	s.conf = cfg
	cfg.setParams()
	return s
}

// MinuteKey
//  @Description: 生成限流redis key
//  @Param uniqueID
//  @Param path
//  @Param burst
//  @Return string
func MinuteKey(uniqueID string, path string, burst int) string {
	return fmt.Sprintf(prefixMinuteKey, uniqueID, path, burst)
}

func hourKey(uniqueID string, path string, burst int) string {
	return fmt.Sprintf(prefixHourKey, uniqueID, path, burst)
}

func (s *RateLimiter) WithRedis(r *redis.Redis) *RateLimiter {
	s.Redis = r
	return s
}

func (s *RateLimiter) Find(resource string) []string {
	return s.cfg.resourceAns[resource]
}

// setParams
//  @Description:设置限流参数
//  @Receiver c
func (c *RLConfig) setParams() {
	if c != nil {
		if c.Minutes < defMinutes {
			c.Minutes = defMinutes
		}
		// 设置分钟与小时级别上限
		// limit most times per Minutes
		if c.N/c.Minutes > defN {
			c.N = defN * c.Minutes
		}
		if c.Hour < defHour {
			c.Hour = defHour
		}
		// 设置分钟与小时级别下限
		if c.M/c.Hour < defM {
			c.M = defM * c.Hour
		}
		// fix uids
		if c.UniqueIDs == "" {
			c.UniqueIDs = defUDid
		}
	}
}

// RateLimiter
//  @Description: 限流操作
//  @Receiver s 超过限流阈值 true
//  @Return bool
func (s *RateLimiter) RateLimiter() bool {
	if s.conf.On {
		uStr := ""
		if s.conf.UniqueIDs == defUDid {
			uStr = s.conf.UDID
		} else if s.conf.UniqueIDs == defIp {
			uStr = s.conf.IP
		}
		if uStr == "" {
			return false
		}
		if err := s.HourTotal(uStr, s.conf.Path, s.conf.Hour, s.conf.M); err != nil {
			return true
		}
		if err := s.MinuteRate(uStr, s.conf.Path, s.conf.Minutes, s.conf.N); err != nil {
			return true
		}
	}

	return false
}

// MinuteRate
//  @Description: 分钟级限流
//  @Receiver s
//  @Param uniqueID 唯一值
//  @Param path 限流路径
//  @Param minutes 配置频率（分钟级）
//  @Param count 每minutes单位限流上限
//  @Return err
func (s *RateLimiter) MinuteRate(uniqueID, path string, minutes, count int) (err error) {
	curMinutes := int(time.Now().Unix() / 60)
	burst := curMinutes - curMinutes%minutes
	key := MinuteKey(uniqueID, path, burst)
	second := int64(minutes * 60)
	return s.doRateLimit(key, uniqueID, second, count)
}

// HourTotal
//  @Description: 小时级限流
//  @Receiver s
//  @Param uniqueID 唯一值
//  @Param path 限流路径
//  @Param hour 配置频率（小时级）
//  @Param count 每hour单位限流上限
//  @Return err
func (s *RateLimiter) HourTotal(uniqueID, path string, hour, count int) (err error) {
	curHour := int(time.Now().Unix() / 3600)
	burst := curHour - curHour%hour
	key := hourKey(uniqueID, path, burst)
	second := int64((curHour+hour)*3600) - time.Now().Unix()
	return s.doRateLimit(key, uniqueID, second, count)
}

// antispam
//  @Description: 限流redis操作
//  @Receiver s
//  @Param key 限流 redis key
//  @Param uniqueID 唯一值
//  @Param interval 过期时间 单位秒
//  @Param count 计数
//  @Return err
func (s *RateLimiter) doRateLimit(key, uniqueID string, expireTime int64, count int) (err error) {
	curInt, err := s.Redis.IncrWithErr(key)
	if err != nil {
		klog.Errorf("doRateLimte redis incr key:%s err:%v", key, err)
	}
	if (curInt - 1) >= int64(count) {
		klog.Infof("The key: %s has been current limited", key)
		return errs.NewError(ecode.CodeRateLimitError)
	}
	resExpire, err := s.Redis.Expire(key, time.Duration(expireTime)*time.Second)
	if !resExpire || err != nil {
		klog.Errorf("doRateLimte redis expire key:%s time:%v err:%v", key, expireTime, err)
	}
	return nil
}
