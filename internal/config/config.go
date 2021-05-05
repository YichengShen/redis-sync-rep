package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var (
	Conf Config
)

func init() {
	Conf.GetConf("config.yaml")
}

type Config struct {
	MasterIp string `yaml:"MasterIp"`
	MasterPort string `yaml:"MasterPort"`
	ReplicaIp string `yaml:"ReplicaIp"`
	ReplicaPort string `yaml:"ReplicaPort"`

	LogDir string `yaml:"LogDir"`

	NClients int `yaml:"NClients"`
	ClientTimeout int `yaml:"ClientTimeout"`
	NClientRequests int `yaml:"NClientRequests"`
	ClientBatchSize int `yaml:"ClientBatchSize"`

	KeyLen int `yaml:"KeyLen"`
	ValLen int `yaml:"ValLen"`
}

// GetConf reads the yaml configuration into Config struct
func (c *Config) GetConf(configPath string) *Config {
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return c
}