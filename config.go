package main

import (
	"github.com/astaxie/beego/config"
	"io/ioutil"
)

const (
	configFileName = "cfg.ini"
	tokenFile      = ".token"
)

type AccessConfig struct {
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

func SetToken(token string) {
	err := ioutil.WriteFile(tokenFile, []byte(token), 0644)
	if err != nil {
		panic(err)
	}
}

func GetToken() (string, error) {
	tokenByte, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		return "", err
	}

	return string(tokenByte), nil
}

func InitConfig() {
	logger.Println("init config")
	token, err := GetToken()
	if err != nil {
		logger.Println("token is empty")
		syncdCfg.Load()
		request := NewRequest(syncdCfg.access)
		loginToken, err := request.Login()
		if err != nil {
			panic(err)
		}
		logger.Println("set new token:", loginToken)
		token = loginToken
		SetToken(token)
	}else{
		logger.Println("token exists:", token)
		//check valid
	}

}

func (c *SyncdConfig) Load() {
	cfg, err := config.NewConfig("ini", configFileName)
	if err != nil {
		panic("配置文件加载失败")
	}

	syncdCfg.cfg = cfg

	syncdCfg.access.Host = cfg.String("host")
	syncdCfg.access.Username = cfg.String("username")
	syncdCfg.access.Password = cfg.String("password")
}

func (c *SyncdConfig) Save() {
}
