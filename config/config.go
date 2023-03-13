package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Conf struct {
	Edge  EdgeConf  `yaml:"edge"`
	Cloud CloudConf `yaml:"cloud"`
}
type EdgeConf struct {
	Service   EdgeServiceConf
	Listen    ListenConf
	RpcClient RpcClientConf
}
type CloudConf struct {
	Service   CloudServiceConf
	Listen    ListenConf
	RpcServer RpcServerConf
}
type CloudServiceConf struct {
	CloudDataSource string            `yaml:"dataSource"`
	Edges           []CloudToEdgeConf `yaml:"edges"`
	Mqtt            MqttConf          `yaml:"mqtt"`
}
type CloudToEdgeConf struct {
	EdgeId     string `yaml:"edgeId"`
	DataSource string `yaml:"dataSource"`
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
type ListenConf struct {
	Mqtt MqttConf       `yaml:"mqtt"`
	Http HttpServerConf `yaml:"http"`
}

type HttpServerConf struct {
	Port string `yaml:"port"`
}
type RpcServerConf struct {
	Port string `yaml:"port"`
}
type RpcClientConf struct {
	Address string `yaml:"address"`
	Port    string `yaml:"port"`
}

var configFile string = "./etc/config.yaml"

func ParseConfig() (config Conf) {
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
	if err := yaml.Unmarshal(content, &config); err != nil {
		log.Fatal(errors.Wrap(err, "unmarshal config file error"))
	}
	conf, _ := json.Marshal(&config)
	log.Print(string(conf))
	return
}
