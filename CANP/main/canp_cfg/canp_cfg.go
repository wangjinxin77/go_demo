package canp_cfg

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"canp.server/common"
	"gopkg.in/yaml.v2"
)

// 定义配置结构体
type EdgeInfo struct {
	ServerPort        int `yaml:"server_port"`         // 默认值 38000
	SendPublishPeriod int `yaml:"send_publish_period"` // 默认值 10
}

type CarrierInfo struct {
	ServerIP   string `yaml:"server_ip"`   // 默认值 127.0.0.1
	ServerPort int    `yaml:"server_port"` // 默认值 58000
}

type Config struct {
	EdgeInfo    EdgeInfo    `yaml:"edge_info"`
	CarrierInfo CarrierInfo `yaml:"carrier_info"`
}

// 构造函数初始化默认值
func NewConfig() *Config {
	return &Config{
		EdgeInfo: EdgeInfo{
			ServerPort:        common.DefaultPort,
			SendPublishPeriod: 10,
		},
		CarrierInfo: CarrierInfo{
			ServerIP:   common.DefaultCarrierIP,
			ServerPort: common.DefaultCarrierPort,
		},
	}
}

func readConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	cfg := NewConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %v", err)
	}

	return cfg, nil
}

func ReadConfigCurrentDir(path string) (*Config, error) {
	// 获取可执行文件真实路径
	exePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %v", err)
	}

	// 获取可执行文件所在目录
	exeDir := filepath.Dir(exePath)

	// 拼接配置文件绝对路径
	configPath := filepath.Join(exeDir, path)
	return readConfig(configPath)
}

// 打印配置信息
func (c *Config) Print() {
	fmt.Printf("Edge Info: %+v\n", c.EdgeInfo)
	fmt.Printf("Carrier Info: %+v\n", c.CarrierInfo)
}
