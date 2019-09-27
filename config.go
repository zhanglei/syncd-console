package main

import (
	"github.com/astaxie/beego/config"
)

const (
	configFileName = "cfg.ini"
)

type AccessConfig struct {
	Schema     string
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

func InitConfig() {
	logger.Println("init config")
	syncdCfg.Load()
	//token, err := GetToken()
	//if err != nil {
	//	logger.Println("token is empty")
	//	request := NewRequest(syncdCfg.access)
	//	request.Login()
	//
	//}else{
	//	logger.Println("token exists:", token)
	//	//check valid
	//}

}

func (c *SyncdConfig) Load() {
	cfg, err := config.NewConfig("ini", configFileName)
	if err != nil {
		panic("配置文件加载失败")
	}

	syncdCfg.cfg = cfg

	syncdCfg.access.Schema = cfg.String("schema")
	syncdCfg.access.Host = cfg.String("host")
	syncdCfg.access.Username = cfg.String("username")
	syncdCfg.access.Password = cfg.String("password")
}

func (c *SyncdConfig) Save() {
}
