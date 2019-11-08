package yaml

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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

type NodeYaml struct {
	Nodes struct {
		IP []string `yaml:"IP"`
	}
}

func ReadConfigYaml(File string) (conf ConfigYaml, err error) {
	conf = ConfigYaml{}
	yamlFile, err := ioutil.ReadFile(File)
	if err != nil {
		fmt.Printf("Reading config file error:%v ", err)
		return ConfigYaml{}, err
	}
	if err = yaml.Unmarshal(yamlFile, &conf); err != nil {
		fmt.Printf("Parsing config file error: %v", err)
		return ConfigYaml{}, err
	}
	return conf, nil
}

func ReadNodeYaml(File string) (conf NodeYaml, err error) {
	conf = NodeYaml{}
	yamlFile, err := ioutil.ReadFile(File)
	if err != nil {
		fmt.Printf("Reading config file error:%v ", err)
		return NodeYaml{}, err
	}
	if err = yaml.Unmarshal(yamlFile, &conf); err != nil {
		fmt.Printf("Parsing config file error: %v", err)
		return NodeYaml{}, err
	}
	return conf, nil
}
