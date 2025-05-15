package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type NacosConfig struct {
	Host        string
	Port        uint64
	Namespace   string
	User        string
	Password    string
	DataId      string
	Group       string
	Service     string
	IP          string
	ServicePort uint64
}

type AppConfig struct {
	Port          string `json:"PORT"`
	DBHost        string `json:"DB_HOST"`
	DBPort        int    `json:"DB_PORT"`
	DBName        string `json:"DB_NAME"`
	DBUser        string `json:"DB_USER"`
	DBPassword    string `json:"DB_PASSWORD"`
	RedisHost     string `json:"REDIS_HOST"`
	RedisPort     int    `json:"REDIS_PORT"`
	RedisUser     string `json:"REDIS_USERNAME"`
	RedisPassword string `json:"REDIS_PASSWORD"`
	RedisDB       int    `json:"REDIS_DB"`
	// ... 其他欄位
}

// 讀取 .env 並回傳 NacosConfig
func LoadNacosConfigFromEnv() (NacosConfig, error) {
	_ = godotenv.Load()
	port, _ := strconv.ParseUint(os.Getenv("NACOS_PORT"), 10, 64)
	servicePort, _ := strconv.ParseUint(os.Getenv("NACOS_SERVICE_PORT"), 10, 64)
	cfg := NacosConfig{
		Host:        os.Getenv("NACOS_HOST"),
		Port:        port,
		Namespace:   os.Getenv("NACOS_NAMESPACE"),
		User:        os.Getenv("NACOS_USERNAME"),
		Password:    os.Getenv("NACOS_PASSWORD"),
		DataId:      os.Getenv("NACOS_DATAID"),
		Group:       os.Getenv("NACOS_GROUP"),
		Service:     os.Getenv("NACOS_SERVICE_NAME"),
		IP:          os.Getenv("NACOS_IP"),
		ServicePort: servicePort,
	}
	fmt.Printf("NacosConfig: %+v\n", cfg)
	return cfg, nil
}

// 建立 ConfigClient 與 NamingClient
func NewNacosClients(cfg NacosConfig) (config_client.IConfigClient, naming_client.INamingClient, error) {
	sc := []constant.ServerConfig{{
		IpAddr: cfg.Host,
		Port:   cfg.Port,
	}}
	cc := &constant.ClientConfig{
		NamespaceId:         cfg.Namespace,
		Username:            cfg.User,
		Password:            cfg.Password,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "./nacos/log",
		CacheDir:            "./nacos/cache",
		LogLevel:            "info",
	}
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		return nil, nil, err
	}
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		return nil, nil, err
	}
	return configClient, namingClient, nil
}

// 取得並印出設定檔案內容
func PrintConfig(cfgClient config_client.IConfigClient, cfg NacosConfig) error {
	content, err := cfgClient.GetConfig(vo.ConfigParam{
		DataId: cfg.DataId,
		Group:  cfg.Group,
	})
	if err != nil {
		return err
	}
	fmt.Println("Nacos 設定檔內容:")
	fmt.Println(content)
	return nil
}

// 註冊服務
func RegisterService(namingClient naming_client.INamingClient, cfg NacosConfig) error {
	_, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          cfg.IP,
		Port:        cfg.ServicePort,
		ServiceName: cfg.Service,
		Weight:      1.0,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})
	if err != nil {
		return err
	}
	fmt.Println("服務註冊成功:", cfg.Service)
	return nil
}

// ParseConfig 取得並解析 Nacos 設定檔內容
func ParseConfig(cfgClient config_client.IConfigClient, cfg NacosConfig) (AppConfig, error) {
	content, err := cfgClient.GetConfig(vo.ConfigParam{
		DataId: cfg.DataId,
		Group:  cfg.Group,
	})
	if err != nil {
		return AppConfig{}, err
	}
	var appCfg AppConfig
	if err := json.Unmarshal([]byte(content), &appCfg); err != nil {
		return AppConfig{}, fmt.Errorf("Nacos 設定檔 JSON 解析失敗: %w", err)
	}
	return appCfg, nil
}
