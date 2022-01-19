// @description
// @author yixia
// Copyright 2021 sndks.com. All rights reserved.
// @datetime 2021/1/14 5:21 下午
// @lastmodify 2021/1/14 5:21 下午

package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/klog"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/kutils/kgo"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

/*
基于http的配置轮询的配置获取
*/
type yaseeDataSource struct {
	lastRevision int64
	enableWatch  bool
	client       *resty.Client
	addr         string
	changed      chan struct{}
	data         string
}

// default client resp struct
type yaseeRes struct {
	Code int        `json:"code"`
	Msg  string     `json:"msg"`
	Data ConfigData `json:"data"`
}

// ConfigData ...
type ConfigData struct {
	Content      string `json:"content"`
	LastRevision int64  `json:"last_revision"`
}

// NewDataSource ...
func NewDataSource(addr string, enableWatch bool) *yaseeDataSource {
	yasee := &yaseeDataSource{
		client:      resty.New(),
		addr:        addr,
		changed:     make(chan struct{}),
		enableWatch: enableWatch,
	}
	if enableWatch {
		kgo.Go(yasee.watch)
	}
	return yasee
}

// ReadConfig ...
func (y *yaseeDataSource) ReadConfig() ([]byte, error) {
	// 检查watch 如果watch为真，走长轮询逻辑
	switch y.enableWatch {
	case true:
		return []byte(y.data), nil
	default:
		content, err := y.getConfigInner(y.addr, y.enableWatch)
		return []byte(content), err
	}
}

// IsConfigChanged ...
func (y *yaseeDataSource) IsConfigChanged() <-chan struct{} {
	return y.changed
}

// Close ...
func (y *yaseeDataSource) Close() error {
	close(y.changed)
	return nil
}

func (y *yaseeDataSource) getContext() context.Context {
	return context.TODO()
}

func (y *yaseeDataSource) watch() {
	for {
		resp, err := y.client.R().SetQueryParam("watch", strconv.FormatBool(y.enableWatch)).Get(y.addr)
		// client get err
		if err != nil {
			time.Sleep(time.Second * 1)
			klog.Error(y.getContext(), "yaseeDataSource", klog.String("listenConfig curl err", err.Error()))
			continue
		}
		if resp.StatusCode() != 200 {
			time.Sleep(time.Second * 1)
			klog.Error(y.getContext(), "yaseeDataSource", klog.String("listenConfig status err", resp.Status()))
		}
		var yaseeRes yaseeRes
		if err := json.Unmarshal(resp.Body(), &yaseeRes); err != nil {
			time.Sleep(time.Second * 1)
			klog.Error(y.getContext(), "yaseeDataSource", klog.String("unmarshal err", err.Error()))
			continue
		}
		// default code != 200 means not change
		if yaseeRes.Code != 200 {
			time.Sleep(time.Second * 1)
			klog.Info(y.getContext(), "yaseeDataSource", klog.Int64("code", int64(yaseeRes.Code)))
			continue
		}
		select {
		case y.changed <- struct{}{}:
			// record the config change data
			y.data = yaseeRes.Data.Content
			y.lastRevision = yaseeRes.Data.LastRevision
			klog.Info(y.getContext(), "yaseeDataSource", klog.String("change", yaseeRes.Data.Content))
		default:
		}
	}
}

func (y *yaseeDataSource) getConfigInner(addr string, enableWatch bool) (string, error) {
	var content string
	//todo: 从基础服务获取数据
	return content, nil
}

func (y *yaseeDataSource) getConfig(addr string, enableWatch bool) (string, error) {
	resp, err := y.client.SetDebug(true).R().SetQueryParam("watch", strconv.FormatBool(enableWatch)).Get(addr)
	if err != nil {
		return "", errors.New("get config err")
	}
	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("get config reply err code:%v", resp.Status())
	}
	configRes := yaseeRes{}
	if err := json.Unmarshal(resp.Body(), &configRes); err != nil {
		return "", fmt.Errorf("unmarshal config err:%v", err.Error())
	}
	if configRes.Code != 200 {
		return "", fmt.Errorf("get config reply err code:%v", resp.Status())
	}
	return configRes.Data.Content, nil
}
