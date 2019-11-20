package yaml

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type ConfigYaml struct {
	Pump2 struct {
		ServerIP   string `yaml:"serverip"`
		ServerPort string `yaml:"serverport"`
		TLS        struct {
			TLSKey string `yaml:"tlskey"`
			TLSCrt string `yaml:"tlscrt"`
		}
	}
}

type Node struct {
	IP string `yaml:"ip"`
}

type NodeYaml struct {
	Nodes []Node `yaml:"nodes"`
}

func ReadConfigYaml(File string) (conf ConfigYaml, err error) {
	conf = ConfigYaml{}
	yamlFile, err := ioutil.ReadFile(File)
	if err != nil {
		log.Printf("Reading config file error:%v ", err)
		return ConfigYaml{}, err
	}
	if err = yaml.Unmarshal(yamlFile, &conf); err != nil {
		log.Printf("Parsing config file error: %v", err)
		return ConfigYaml{}, err
	}
	return conf, nil
}

func ReadNodeYaml(File string) (conf NodeYaml, err error) {
	conf = NodeYaml{}
	yamlFile, err := ioutil.ReadFile(File)
	if err != nil {
		log.Printf("Reading config file error:%v ", err)
		return NodeYaml{}, err
	}
	if err = yaml.Unmarshal(yamlFile, &conf); err != nil {
		log.Printf("Parsing config file error: %v", err)
		return NodeYaml{}, err
	}
	return conf, nil
}
