package conf

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Env string

// Cfg 配置文件. *Cfg后字段的值可以不用指针了. 可以添加required tag字段，设为false则跳过空值检查.
type Cfg struct {
	AppName     string           `yaml:"appName"`
	Environment string           `yaml:"environment"` //test online 默认online
	HTTPServer  []HttpServerConf `yaml:"httpServer"`
	Redis       RedisConf        `yaml:"redis"`
	PprofPort   uint16           `yaml:"pprofPort"`
	Pulsar      PulsarConf       `yaml:"pulsar"`
	Database    DBConf           `yaml:"database"`
	LogDir      string           `yaml:"logDir"`
	ConstantKey ConstantKey      `yaml:"constantKey"`
}

func NewCfg(env Env) (*Cfg, error) {
	b, err := os.ReadFile(string(env))
	if err != nil {
		return nil, fmt.Errorf("配置文件错误. err: %w", err)
	}
	var cfg Cfg
	err = yaml.Unmarshal(b, &cfg)
	if err != nil {
		return nil, fmt.Errorf("配置文件错误. err: %w", err)
	}
	return &cfg, nil
}

// HttpServerConf 有关时间的整数设置，均为毫秒.
type HttpServerConf struct {
	CertFile string   `yaml:"certFile"`
	KeyFile  string   `yaml:"keyFile"`
	Port     int      `yaml:"port"`
	Reverses []string `yaml:"reverses"` // 反向代理的ip地址
}

// PulsarConf .
type PulsarConf struct {
	ServiceURL string `yaml:"serviceUrl"`
}

type DBConf struct {
	Dsn        string `yaml:"dsn"` // todo 日志不输出Dsn
	DriverName string `yaml:"driverName"`
}

type ConstantKey struct {
	ResourceDir string `yaml:"resourceDir,omitempty"`
	ResourceURI string `yaml:"resourceURI,omitempty"`
}

type RedisConf struct {
	Addr     string `yaml:"addr"`
	UserName string `yaml:"userName"`
	Password string `yaml:"password" json:"-"`
	DB       int    `yaml:"db"`
}
