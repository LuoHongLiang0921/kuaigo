package server

import (
	"context"
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg"
	"github.com/LuoHongLiang0921/kuaigo/pkg/constant"
)

type Option func(c *ServiceInfo)

// Server
type Server interface {
	Serve() error
	Stop() error
	GracefulStop(ctx context.Context) error
	Info() *ServiceInfo
}
type ServiceInfo struct {
	Name     string               `json:"name"`
	AppID    string               `json:"appId"`
	Scheme   string               `json:"scheme"`
	Address  string               `json:"address"`
	Weight   float64              `json:"weight"`
	Enable   bool                 `json:"enable"`
	Healthy  bool                 `json:"healthy"`
	Metadata map[string]string    `json:"metadata"`
	Region   string               `json:"region"`
	Zone     string               `json:"zone"`
	Kind     constant.ServiceKind `json:"kind"`
	// Deployment 部署组: 不同组的流量隔离
	// 比如某些服务给内部调用和第三方调用，可以配置不同的deployment,进行流量隔离
	Deployment string `json:"deployment"`
	// Group 流量组: 流量在Group之间进行负载均衡
	Group    string              `json:"group"`
	Services map[string]*Service `json:"services" toml:"services"`
}

// Service ...
type Service struct {
	Namespace string            `json:"namespace" toml:"namespace"`
	Name      string            `json:"name" toml:"name"`
	Labels    map[string]string `json:"labels" toml:"labels"`
	Methods   []string          `json:"methods" toml:"methods"`
}

// Label
//  @Description  启动服务的scheme与adress
//  @Receiver si ServiceInfo
//  @Return string
func (si ServiceInfo) Label() string {
	return fmt.Sprintf("%s://%s", si.Scheme, si.Address)
}

func WithMetaData(key, value string) Option {
	return func(c *ServiceInfo) {
		c.Metadata[key] = value
	}
}

func WithScheme(scheme string) Option {
	return func(c *ServiceInfo) {
		c.Scheme = scheme
	}
}

func WithAddress(address string) Option {
	return func(c *ServiceInfo) {
		c.Address = address
	}
}

func WithKind(kind constant.ServiceKind) Option {
	return func(c *ServiceInfo) {
		c.Kind = kind
	}
}

// ApplyOptions
//  @Description  组装配置项
//  @Param options
//  @Return ServiceInfo
func ApplyOptions(options ...Option) ServiceInfo {
	info := defaultServiceInfo()
	for _, option := range options {
		option(&info)
	}
	return info
}

// defaultServiceInfo
//  @Description  默认服务信息
//  @Return ServiceInfo
func defaultServiceInfo() ServiceInfo {
	si := ServiceInfo{
		Name:       pkg.GetAppName(),
		AppID:      pkg.GetAppID(),
		Weight:     100,
		Enable:     true,
		Healthy:    true,
		Metadata:   make(map[string]string),
		Region:     pkg.GetAppRegion(),
		Zone:       pkg.GetAppZone(),
		Kind:       0,
		Deployment: "",
		Group:      "",
	}
	si.Metadata["appMode"] = pkg.GetAppMode()
	si.Metadata["appHost"] = pkg.GetAppHost()
	si.Metadata["startTime"] = pkg.GetStartTime()
	si.Metadata["buildTime"] = pkg.GetBuildTime()
	si.Metadata["appVersion"] = pkg.GetAppVersion()
	si.Metadata["tabbyVersion"] = pkg.GetTabbyVersion()
	return si
}
