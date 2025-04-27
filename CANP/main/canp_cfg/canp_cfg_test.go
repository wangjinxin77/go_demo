package canp_cfg

import (
	"io/ioutil"
	"os"
	"testing"

	"canp.server/common"
)

func TestReadConfig(t *testing.T) {
	tests := []struct {
		name     string
		fileData string
		expected Config
	}{
		{
			name: "完整配置",
			fileData: `
edge_info:
  server_port: 39000
  send_publish_period: 30
carrier_info:
  server_ip: 10.27.33.2
  server_port: 8081
`,
			expected: Config{
				EdgeInfo: EdgeInfo{
					ServerPort:        39000,
					SendPublishPeriod: 30,
				},
				CarrierInfo: CarrierInfo{
					ServerIP:   "10.27.33.2",
					ServerPort: 8081,
				},
			},
		},
		{
			name: "缺失字段配置",
			fileData: `
edge_info:
  server_port: 39000
carrier_info:
  server_ip: 10.27.33.2
`,
			expected: Config{
				EdgeInfo: EdgeInfo{
					ServerPort:        39000,
					SendPublishPeriod: 10, // 默认值
				},
				CarrierInfo: CarrierInfo{
					ServerIP:   "10.27.33.2",
					ServerPort: common.DefaultCarrierPort, // 默认值
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := ioutil.TempFile("", "config-*.yaml")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpFile.Name())

			if _, err := tmpFile.Write([]byte(tt.fileData)); err != nil {
				t.Fatal(err)
			}
			tmpFile.Close()

			cfg, err := readConfig(tmpFile.Name())
			if err != nil {
				t.Errorf("ReadConfig() error = %v", err)
			}

			if cfg.EdgeInfo.ServerPort != tt.expected.EdgeInfo.ServerPort {
				t.Errorf("EdgeInfo.ServerPort = %v, want %v", cfg.EdgeInfo.ServerPort, tt.expected.EdgeInfo.ServerPort)
			}
			if cfg.EdgeInfo.SendPublishPeriod != tt.expected.EdgeInfo.SendPublishPeriod {
				t.Errorf("SendPublishPeriod = %v, want %v", cfg.EdgeInfo.SendPublishPeriod, tt.expected.EdgeInfo.SendPublishPeriod)
			}
			if cfg.CarrierInfo.ServerIP != tt.expected.CarrierInfo.ServerIP {
				t.Errorf("ServerIP = %v, want %v", cfg.CarrierInfo.ServerIP, tt.expected.CarrierInfo.ServerIP)
			}
			if cfg.CarrierInfo.ServerPort != tt.expected.CarrierInfo.ServerPort {
				t.Errorf("ServerPort = %v, want %v", cfg.CarrierInfo.ServerPort, tt.expected.CarrierInfo.ServerPort)
			}
		})
	}
}
