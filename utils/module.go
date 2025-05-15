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
			NewRedisClient,         // 回傳 *redis.Client
			NewGormDB,              // 回傳 *gorm.DB
		),
		fx.Invoke(
			InitNacos, // 取 config 並註冊服務
		),
	)
}

// InitNacos 取 config 並註冊服務
func InitNacos(cfg NacosConfig, configClient interface{}, namingClient interface{}) error {
	cc, _ := configClient.(config_client.IConfigClient)
	nc, _ := namingClient.(naming_client.INamingClient)
	if err := PrintConfig(cc, cfg); err != nil {
		return err
	}
	if err := RegisterService(nc, cfg); err != nil {
		return err
	}
	return nil
}
