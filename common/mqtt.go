package common

import (
	"encoding/json"
	"log"
	"patient-edge/config"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
)

func NewMqttClient(config config.MqttConf) mqtt.Client {
	// mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)
	client := mqtt.NewClient(mqtt.NewClientOptions().AddBroker(config.Broker).SetClientID(config.ClientId))
	if tc := client.Connect(); tc.Wait() && tc.Error() != nil {
		log.Fatal(errors.Wrap(tc.Error(), "connect to broker error"))
	}
	return client
}

func MqttSubscribe[T any](client mqtt.Client, topic string, handler func(req T)) {
	if tc := client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("receive message %s", msg.Payload())
		log.Printf("from topic %s\n", msg.Topic())
		var req T
		if err := json.Unmarshal(msg.Payload(), &req); err != nil {
			log.Fatal(errors.Wrap(err, "unmarshal payload error"))
		}
		handler(req)
		return
	}); tc.Wait() && tc.Error() != nil {
		log.Fatal(errors.Wrap(tc.Error(), "subscribe topic "+topic+"error"))
	}
}
func MqttPublish(client mqtt.Client, topic string, paramNames []string, params []any) {
	var err error
	mp := make(map[string]any)
	for i, name := range paramNames {
		mp[name] = params[i]
	}
	payload, err := json.Marshal(mp)
	if err != nil {
		log.Fatal(errors.Wrap(err, "marshal params error"))
	}
	log.Printf("publish to topic: %s \n", topic)
	if tc := client.Publish(topic, 0, false, payload); tc.Wait() && tc.Error() != nil {
		log.Fatal(errors.Wrap(tc.Error(), "publish sync error"))
	}
}
