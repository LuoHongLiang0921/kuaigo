package test

import (
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/conf"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/configsource/apollo"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/configsource/file"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/configsource/manager"
	"os"

	"gopkg.in/yaml.v3"
)

// InitTestForFile 必选先设置环境变量 CONFIG_FILE_ADDR= ""
func InitTestForFile() error {
	file.RegisterConfigHandler()
	configAddr := os.Getenv("CONFIG_FILE_ADDR")
	err := doLoadConfig(configAddr)
	if err != nil {
		return err
	}
	fmt.Println("v2", "tabby-test")
	return nil
}

// InitTestForApollo 必选先设置apollo 变量 APOLLO_CONFIG_FILE=""
func InitTestForApollo() error {
	apollo.RegisterConfigHandler()
	configAddr := os.Getenv("APOLLO_CONFIG_FILE")

	err := doLoadConfig(configAddr)
	if err != nil {
		return err
	}
	fmt.Println("v2", "tabby-test")
	return nil
}

func doLoadConfig(schema string) error {
	provider, err := manager.NewConfigSource(schema)
	if err == manager.ErrConfigAddr {
		return err
	}
	if err := conf.LoadFromConfigSource(provider, yaml.Unmarshal); err != nil {
		return err
	}
	return nil
}
