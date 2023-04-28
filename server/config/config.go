package config

import (
	"io/ioutil"
	"log"
	"path"
	"runtime"

	"gopkg.in/yaml.v2"
)

func GetRootFilepath() string {
	_, filename, _, _ := runtime.Caller(0)
	// TODO：先暂时这样获取项目根目录，之后有更好的办法再说
	// 		 一旦项目整体路径发生更改该方法就有可能失效
	return path.Dir(path.Dir(path.Dir(filename)))
}

type Config struct {
	INFRASTRUCTURE_USER struct {
		NAME     string `yaml:"name"`
		PASSWORD string `yaml:"password"`
	} `yaml:"infra_user"`
	ADDRESSES struct {
		OSS_ADDR      string `yaml:"oss_addr"`
		OSS_PORT      string `yaml:"oss_port"`
		PGSQL_ADDR    string `yaml:"pgsql_addr"`
		PGSQL_PORT    string `yaml:"pgsql_port"`
		PGSQL_DB_NAME string `yaml:"pgsql_db_name"`
	} `yaml:"addresses"`
	WECHAT struct {
		APPID            string `yaml:"appid"`
		SECRET           string `yaml:"secret"`
		MCHID            string `yaml:"mchid"`
		MCHCERTSERIALNUM string `yaml:"mch_certificate_serial_number"`
		MCHAPIV3KEY      string `yaml:"mch_apiv3_key"`
	} `yaml:"wechat"`
}

var config *Config = nil

func ReadConfig() *Config {
	var f []byte
	rootPath := GetRootFilepath()
	f, err := ioutil.ReadFile(path.Join(rootPath, "config.yaml"))
	if err != nil {
		f, err = ioutil.ReadFile(path.Join(rootPath, "config/config.yaml"))
		if err != nil {
			log.Fatalln("no configure file 'config.yaml' or 'config/config.yaml' " + err.Error())
		}
	}

	err = yaml.Unmarshal(f, &config)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return config
}

func GetConfig() Config {
	if config == nil {
		config = ReadConfig()
	}
	return *config
}
