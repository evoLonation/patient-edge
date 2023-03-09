package common

import (
	"log"
	"patient-edge/config"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
)

func GetMqttClient(config config.MqttConf) mqtt.Client {
	// mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)
	client := mqtt.NewClient(mqtt.NewClientOptions().AddBroker(config.Broker).SetClientID(config.ClientId))
	if tc := client.Connect(); tc.Wait() && tc.Error() != nil {
		log.Fatal(errors.Wrap(tc.Error(), "connect to broker error"))
	}
	return client
}
