package config

import (
	"io/ioutil"
	"log"
	"path"
	"runtime"
	"time"

	"gopkg.in/yaml.v2"
)

func getRootFilePath() string {
	_, filename, _, _ := runtime.Caller(0)
	// TODO：先暂时这样获取项目根目录，之后有更好的办法再说
	// 		 一旦项目整体路径发生更改该方法就有可能失效
	return path.Dir(path.Dir(path.Dir(filename)))
}

type Config struct {
	Server struct {
		Timeout  time.Duration `yaml:"timeout"`
		CertFile string        `yaml:"certfile"`
		KeyFile  string        `yaml:"keyfile"`
	} `yaml:"server"`
	MongoDB struct {
		Timeout time.Duration `yaml:"timeout"`
		URI     string        `yaml:"uri"`
	} `yaml:"mongodb"`
	InfrastructureUser struct {
		NAME     string `yaml:"name"`
		PASSWORD string `yaml:"password"`
	} `yaml:"infra_user"`
	Addresses struct {
		OssAddr     string `yaml:"oss_addr"`
		OssPort     string `yaml:"oss_port"`
		PgsqlAddr   string `yaml:"pgsql_addr"`
		PgsqlPort   string `yaml:"pgsql_port"`
		PgsqlDbName string `yaml:"pgsql_db_name"`
	} `yaml:"addresses"`
	Wechat struct {
		AppID            string `yaml:"appid"`
		Secret           string `yaml:"secret"`
		MchID            string `yaml:"mchid"`
		MchCertSerialNum string `yaml:"mch_certificate_serial_number"`
		MchAPIV3Key      string `yaml:"mch_apiv3_key"`
	} `yaml:"wechat"`
}

var config *Config = nil

func ReadConfig() *Config {
	var f []byte
	rootPath := getRootFilePath()
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
