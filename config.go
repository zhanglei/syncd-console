package main

import (
	"github.com/astaxie/beego/config"
)

var configFileName = "syncdcfg.ini"
var loaded = false

type Configs struct {
	Secure   bool
	Host     string
	Username string
	Password string
}

//type Config struct {
//	FilePath string
//	Cfg      config.Configer
//	Items    ConfigItems
//}

var syncdCfg *Configs

func init()  {
	InitSyncdConfig()
}

func InitSyncdConfig() {
	if !loaded {
		syncdCfg = &Configs{}
		cfg, err := config.NewConfig("ini", configFileName)
		if err != nil {
			panic("配置文件加载失败")
		}

		secure, err := cfg.Bool("secure")
		if err != nil {
			panic("无法加载配置secure")
		}
		syncdCfg.Secure = secure
		syncdCfg.Host = cfg.String("host")
		syncdCfg.Username = cfg.String("username")
		syncdCfg.Password = cfg.String("password")
	}
}
