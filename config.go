package main

import (
	"github.com/astaxie/beego/config"
	z "github.com/nutzam/zgo"
)

const (
	configFileName = "syncd-console.ini"
)

type AccessConfig struct {
	Schema   string
	Host     string
	Username string
	Password string
}

type SyncdConfig struct {
	configPath string
	cfg        config.Configer
	access     AccessConfig
	loaded     bool
}

var syncdCfg SyncdConfig

func InitConfig() AccessConfig {
	logger.Println("loading config")
	syncdCfg.Load()

	return syncdCfg.access
}

func (c *SyncdConfig) Load() {
	cfg, err := config.NewConfig("ini", configFileName)
	if err != nil {
		panic("配置文件加载失败")
	}

	syncdCfg.cfg = cfg

	syncdCfg.access.Schema = z.Trim(cfg.String("schema"))
	syncdCfg.access.Host = z.Trim(cfg.String("host"))
	syncdCfg.access.Username = z.Trim(cfg.String("username"))
	syncdCfg.access.Password = z.Trim(cfg.String("password"))
	if z.IsBlank(syncdCfg.access.Username)||z.IsBlank(syncdCfg.access.Host)||z.IsBlank(syncdCfg.access.Username)||z.IsBlank(syncdCfg.access.Password){
		panic("请先设置配置文件 syncd-console.ini 的参数")
	}
}

func (c *SyncdConfig) Save() {
}
