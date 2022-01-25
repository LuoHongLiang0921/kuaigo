package kconfig

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/net/xhttp"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"sync"
	"time"

)

var appIDMapper sync.Map

const defaultInterval = 2 * time.Second

// 配置服务
type AppConfig struct {
	//ServiceName 服务名字
	ServiceName string `json:"serviceName"`
	//AppIDByPkgURL 据包名查询appId
	AppIDByPkgURL string `yaml:"appIdByPkgUrl"`
	//GetConfigByKey 根据keys查询配置信息
	ConfigByKeyUrl string `yaml:"configByKeyUrl"`
	//GetKongCertUrl 根据包名查询kong网关证书
	KongCertUrl string `yaml:"kongCertUrl"`
	// 获取配置间隔,单位为秒
	PollInterval int `yaml:"longPollInterval"`
}

type GetAppIDReq struct {
	Pkg string `json:"pkg"`
}

type GetAppIDResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

type GetByKeyReq struct {
	Key string `json:"key"`
}

type GetKongCretReq struct {
	Pkg string `json:"pkg"`
}

// GetConfigByKeyResp 返回内容
type GetConfigByKeyResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data *ConfigData `json:"data"`
}

// ConfigData ...
type ConfigData struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	AppId      int    `json:"appId"`
	BusinessId int    `json:"businessId"`
	Key        string `json:"key"`
	Value      string `json:"value"`
	Summary    string `json:"summary"`
	CreateTime int64  `json:"createTime"`
	CreateBy   int    `json:"createBy"`
	UpdateTime int64  `json:"updateTime"`
	UpdateBy   int    `json:"updateBy"`
}

// Build
//  @Description
//  @Param ctx
//  @Param key
//  @Return *AppConfig
func Build(ctx context.Context, key string) *AppConfig {
	var config AppConfig
	err := conf.UnmarshalKey(key, &config)
	if err != nil {
		klog.TabbyLogger.Panic(fmt.Sprintf("load key %v config centre err:%v", key, err), klog.FieldMod("biz"))
		return nil
	}
	return &config
}

// Start
//  @Description  默认启动
//  @Receiver c
//  @Param ctx
//  @Param key 配置服务
//  @Param resp
//  @Return error
func (c *AppConfig) Start(ctx context.Context, key string, resp interface{}) error {
	var interval time.Duration
	if c.PollInterval != 0 {
		interval = time.Duration(c.PollInterval) * time.Second
	} else {
		interval = defaultInterval
	}
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			klog.WithContext(ctx).Debug("start load config key")
			data, err := c.GetConfigByKey(ctx, key)
			if err != nil {
				klog.WithContext(ctx).Errorf("key %v load config from centre server err:%v", key, err)
				return err
			}
			err = json.Unmarshal([]byte(data.Value), &resp)
			if err != nil {
				klog.WithContext(ctx).Errorf("unmarshal key %v,err:%v ", key, err)
			}
		}
	}
	return nil
}

// GetKongCert
//  @Description  根据包名查询kong网关证书
//  @Receiver c
//  @Param ctx
//  @Param pkg 包名
//  @Return *ConfigData
//  @Return error
func (c *AppConfig) GetKongCert(ctx context.Context, pkg string) (*ConfigData, error) {
	req := &GetKongCretReq{
		Pkg: pkg,
	}
	configResp := &GetConfigByKeyResp{}
	resp,err := khttp.Post(ctx, c.KongCertUrl, req)
	if err != nil {
		return nil, err
	}
	resp.Json(configResp)
	return configResp.Data, err
}

//  @Deprecated 请使用 GetKongCert 替代 [废弃]
func (c *AppConfig) GetKongCret(ctx context.Context, pkg string) (*ConfigData, error) {
	return c.GetKongCert(ctx, pkg)
}

// GetConfigByKey
//  @Description  根据keys查询配置信息
//  @Receiver c
//  @Param ctx
//  @Param key 配置key
//  @Return *ConfigData
//  @Return error
func (c *AppConfig) GetConfigByKey(ctx context.Context, key string) (*ConfigData, error) {
	req := &GetByKeyReq{
		Key: key,
	}
	cfgResp := &GetConfigByKeyResp{}
	resp,err := khttp.PostJson(ctx, c.ConfigByKeyUrl, req)
	if err != nil {
		return nil, err
	}
	resp.Json(cfgResp)
	return cfgResp.Data, err
}

// GetAppID
//  @Description  通过包名获取APPID
//  @Receiver c
//  @Param ctx
//  @Param pkg
//  @Return string
func (c *AppConfig) GetAppID(ctx context.Context, pkg string) string {
	if v, ok := appIDMapper.Load(pkg); ok {
		return v.(string)
	}
	req := &GetAppIDReq{
		Pkg: pkg,
	}
	var appIdResp GetAppIDResp
	resp,err := khttp.PostJson(ctx, c.AppIDByPkgURL, req)
	if err != nil {
		return "0"
	}
	resp.Json(&appIdResp)
	return appIdResp.Data
}
