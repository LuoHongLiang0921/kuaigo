package kconfig

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/core/net/xhttp"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/klog"
	"net/http"
	"sync"
	"time"
)

var appIDMapper sync.Map

const defaultInterval = 2 * time.Second

// 配置服务
type AppConfig struct {
	//ServiceName 服务名字
	ServiceName string `json:"serviceName"`
	//bussinessID string `json:"bussinessId"`
	//AppIDByPkgURL 据包名查询appId
	AppIDByPkgURL string `yaml:"appIdByPkgUrl"`
	//GetConfigByKey 根据keys查询配置信息
	ConfigByKeyUrl string `yaml:"configByKeyUrl"`
	//GetKongCretUrl 根据包名查询kong网关证书
	KongCretUrl string `yaml:"kongCretUrl"`
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

type GetConfigByKeyResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data *ConfigData `json:"data"`
}

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

func Build(ctx context.Context, key string) *AppConfig {
	var config AppConfig
	err := conf.UnmarshalKey(key, &config)
	if err != nil {
		klog.Panic(ctx, fmt.Sprintf("load key %v config centre err:%v", key, err), klog.FieldMod("biz"))
		return nil
	}
	return &config
}

// 默认启动
// key 为配置服务
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
			klog.Debug(ctx, "start load config key")
			data, err := c.GetConfigByKey(ctx, key)
			if err != nil {
				klog.Error(ctx, fmt.Sprintf("key %v load config from centre server err:%v", key, err))
				return err
			}
			err = json.Unmarshal([]byte(data.Value), &resp)
			if err != nil {
				klog.Error(ctx, fmt.Sprintf("unmarshal key %v,err:%v ", key, err))
				//return err
			}
		}
	}
	return nil
}

// GetConfigByKey 根据包名查询konga网关证书
func (c *AppConfig) GetKongCret(ctx context.Context, pkg string) (*ConfigData, error) {
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	req := &GetKongCretReq{
		Pkg: pkg,
	}
	resp := &GetConfigByKeyResp{}
	err := xhttp.PostWithUnmarshal(ctx, nil, c.KongCretUrl, header, req, resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, err
}

// GetConfigByKey 根据keys查询配置信息
func (c *AppConfig) GetConfigByKey(ctx context.Context, key string) (*ConfigData, error) {
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	req := &GetByKeyReq{
		Key: key,
	}
	resp := &GetConfigByKeyResp{}
	err := xhttp.PostWithUnmarshal(ctx, nil, c.ConfigByKeyUrl, header, req, resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, err
}

// GetAppID 同过获取应用id
func (c *AppConfig) GetAppID(ctx context.Context, pkg string) string {
	if v, ok := appIDMapper.Load(pkg); ok {
		return v.(string)
	}
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	req := &GetAppIDReq{
		Pkg: pkg,
	}
	var resp GetAppIDResp
	err := xhttp.PostWithUnmarshal(ctx, nil, c.AppIDByPkgURL, header, req, &resp)
	if err != nil {
		return "0"
	}
	return resp.Data
}
