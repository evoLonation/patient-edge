package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Conf struct {
	Edge   EdgeConf   `yaml:"edge"`
	Cloud  CloudConf  `yaml:"cloud"`
	Common CommonConf `yaml:"common"`
}
type EdgeServiceConf struct {
	EdgeDataSource  string   `yaml:"edgeDataSource"`
	CloudDataSource string   `yaml:"cloudDataSource"`
	Mqtt            MqttConf `yaml:"mqtt"`
}
type MqttConf struct {
	Broker   string `yaml:"broker"`
	ClientId string `yaml:"clientId"`
}

type CloudServiceConf struct {
	DataSource string `yaml:"dataSource"`
}

type EdgeConf struct {
	DataSource string `yaml:"dataSource"`
	MqttBroker string `yaml:"mqttBroker"`
	ClientId   string `yaml:"clientId"`
	Topic      struct {
		ReceiveTemperature string `yaml:"receiveTemperature"`
	} `yaml:"topic"`
}
type CloudConf struct {
	DataSource     string `yaml:"dataSource"`
	MqttBroker     string `yaml:"mqttBroker"`
	ClientId       string `yaml:"clientId"`
	Address        string `yaml:"address"`
	RpcPort        string `yaml:"rpcPort"`
	HttpServerPort string `yaml:"httpServerPort"`
	Topic          struct {
		Notice string `yaml:"notice"`
	} `yaml:"topic"`
}
type CommonConf struct {
	Topic struct {
		Sync string `yaml:"sync"`
	} `yaml:"topic"`
}

var Config Conf

var configFile string = "./etc/config.yaml"

func init() {
	dirs, err := os.ReadDir("./etc")
	if err != nil {
		log.Fatal(errors.Wrap(err, "read directory error"))
	}
	var dirInfo string
	for _, dir := range dirs {
		dirInfo += dir.Name() + ", "
	}
	log.Printf("files: %s\n", dirInfo)

	log.Println("start read config file")
	content, err := os.ReadFile(configFile)
	// log.Print(string(content))
	if err != nil {
		log.Fatal(errors.Wrap(err, "read config file error"))
	}
	if err := yaml.Unmarshal(content, &Config); err != nil {
		log.Fatal(errors.Wrap(err, "unmarshal config file error"))
	}
	conf, _ := json.Marshal(&Config)
	log.Print(string(conf))
}
