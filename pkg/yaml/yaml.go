package yaml

import (
"fmt"
yaml "gopkg.in/yaml.v2"
"io/ioutil"
)

type Yaml struct {
	Pump2 struct {
		ServerIP string `yaml:"serverip"`
		ServerPort string `yaml:"serverport"`
		TLS struct {
			TLSKey string `yaml:"tlskey"`
			TLSCrt string `yaml:"tlscrt"`
		}
	}
}

func ReadYaml(File string) (conf Yaml, err error){
	conf = Yaml{}
	yamlFile, err := ioutil.ReadFile(File)
	if err != nil {
		fmt.Printf("Reading config file error:%v ", err)
		return Yaml{}, err
	}
	if err = yaml.Unmarshal(yamlFile, &conf); err != nil {
		fmt.Printf("Parsing config file error: %v", err)
		return Yaml{}, err
	}
	return conf, nil
}