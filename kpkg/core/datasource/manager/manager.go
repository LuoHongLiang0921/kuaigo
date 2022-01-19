// @description
// @author yixia
// Copyright 2021 sndks.com. All rights reserved.
// @datetime 2021/1/14 5:21 下午
// @lastmodify 2021/1/14 5:21 下午

package manager

import (
	"errors"
	"github.com/LuoHongLiang0921/kuaigo/kpkg/conf"
	"net/url"
)

var (
	//ErrConfigAddr not config
	ErrConfigAddr = errors.New("no config... ")
	// ErrInvalidDataSource defines an error that the scheme has been registered
	ErrInvalidDataSource = errors.New("invalid data source, please make sure the scheme has been registered")
	registry             map[string]DataSourceCreatorFunc
	//DefaultScheme ..
	DefaultScheme string
)

// DataSourceCreatorFunc represents a dataSource creator function
type DataSourceCreatorFunc func() conf.DataSource

func init() {
	registry = make(map[string]DataSourceCreatorFunc)
}

// Register registers a dataSource creator function to the registry
func Register(scheme string, creator DataSourceCreatorFunc) {
	registry[scheme] = creator
}

// CreateDataSource creates a dataSource witch has been registered
// func CreateDataSource(scheme string) (conf.DataSource, error) {
// 	creatorFunc, exist := registry[scheme]
// 	if !exist {
// 		return nil, ErrInvalidDataSource
// 	}
// 	return creatorFunc(), nil
// }

//NewDataSource ..
func NewDataSource(configAddr string) (conf.DataSource, error) {
	if configAddr == "" {
		return nil, ErrConfigAddr
	}
	urlObj, err := url.Parse(configAddr)
	if err == nil && len(urlObj.Scheme) > 1 {
		DefaultScheme = urlObj.Scheme
	}

	creatorFunc, exist := registry[DefaultScheme]
	if !exist {
		return nil, ErrInvalidDataSource
	}
	return creatorFunc(), nil
}
