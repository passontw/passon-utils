package utils

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"go.uber.org/fx"
)

// Module 封裝所有服務連線，供 fx 使用
func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			LoadNacosConfigFromEnv, // 回傳 NacosConfig
			NewNacosClients,        // 回傳 configClient, namingClient
			ParseConfig,            // 回傳 AppConfig
			NewRedisClient,         // 依賴 AppConfig
			NewGormDB,              // 依賴 AppConfig
		),
		fx.Invoke(
			InitNacos, // 取 config 並註冊服務
		),
	)
}

// InitNacos 取 config 並註冊服務
func InitNacos(cfg NacosConfig, configClient config_client.IConfigClient, namingClient naming_client.INamingClient) error {
	if err := PrintConfig(configClient, cfg); err != nil {
		return err
	}
	if err := RegisterService(namingClient, cfg); err != nil {
		return err
	}
	return nil
}
