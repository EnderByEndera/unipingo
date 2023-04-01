package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

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
}

var config *Config = nil

func ReadConfig() *Config {
	var f []byte
	f, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		f, err = ioutil.ReadFile("config/config.yaml")
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
