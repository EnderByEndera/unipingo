package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Question struct {
	Title       string `yaml:"title" json:"title"`
	Description string `yaml:"description" json:"description"`
}

type Topic struct {
	Name        string     `yaml:"name" json:"name"`
	Description string     `yaml:"description" json:"description"`
	Questions   []Question `yaml:"questions" json:"questions"`
}

type Topics []Topic

var topics Topics = nil

func ReadTopics() Topics {
	var f []byte
	f, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		f, err = ioutil.ReadFile("config/topics.yaml")
		if err != nil {
			log.Fatalln("no configure file 'topics.yaml' or 'config/topics.yaml' " + err.Error())
		}
	}
	out, _ := yaml.Marshal(Topics{Topic{Name: "aaa", Questions: []Question{Question{Title: "aaaaa"}}}})
	println(string(out))
	err = yaml.Unmarshal(f, &topics)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return topics
}

func GetTopics() Topics {
	if topics == nil {
		topics = make(Topics, 0)
		topics = ReadTopics()
	}
	return topics
}
