package main

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
)

const (
	IP             = "localhost:1883"
	ClientId       = "PatientChecker-01"
	SubscribeTopic = "$patient/sensor/+/temperature"
)

// connect connect to the Mqtt server.
func connect() (client mqtt.Client, err error) {
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "mqtt connect error")
		}
	}()
	opts := mqtt.NewClientOptions().AddBroker(IP).SetClientID(ClientId).SetCleanSession(true)
	client = mqtt.NewClient(opts)
	if tc := client.Connect(); tc.Wait() && tc.Error() != nil {
		return nil, tc.Error()
	}
	return client, nil
}

func onMessage(client mqtt.Client, message mqtt.Message) {
	log.Printf("receive message %s", message.Payload())
	message.Topic()
}

func main() {
	var err error
	client, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	if tc := client.Subscribe(SubscribeTopic, 0, onMessage); tc.Wait() && tc.Error() != nil {
		log.Fatal(tc.Error())
	}
}
