package core

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Copilot CopilotConfig `yaml:"copilot"`
	Secret  string        `yaml:"secret"`
}

type ServerConfig struct {
	Domain   string `yaml:"domain"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	CertPath string `yaml:"certPath"`
	KeyPath  string `yaml:"keyPath"`
}

type CopilotConfig struct {
	GithubApiUrl string   `yaml:"github_api_url"`
	Tokens       []string `yaml:"tokens"`
}

// 初始化配置文件
func initConfig() Config {
	// 读取配置文件
	configFile, err := os.Open("./config/config.yaml")
	if err != nil {
		panic("file \"./config/config.yaml\" not found")
	}
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			panic("failed to close file \"./config/config.yaml\"")
		}
	}(configFile)
	decoder := yaml.NewDecoder(configFile)
	conf := Config{}
	err = decoder.Decode(&conf)
	if err != nil {
		panic("invalid config format")
	}
	return conf
}

func init() {
	config = initConfig()
}
