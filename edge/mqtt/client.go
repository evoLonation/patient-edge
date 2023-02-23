package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"patient-edge/config"
	"patient-edge/edge/db"
	"patient-edge/edge/operation"
	"patient-edge/edge/rpc"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
)

var Client mqtt.Client

func Start() {
	edgeConfig := &config.Config.Edge
	Client = mqtt.NewClient(mqtt.NewClientOptions().AddBroker(edgeConfig.MqttBroker).SetClientID(edgeConfig.ClientId))
	if tc := Client.Connect(); tc.Wait() && tc.Error() != nil {
		log.Fatal(errors.Wrap(tc.Error(), "connect to broker error"))
	}
	// 数据同步相关
	if tc := Client.Subscribe(config.Config.Common.Topic.Sync, 0, sync); tc.Wait() && tc.Error() != nil {
		log.Fatal(errors.Wrap(tc.Error(), "subscribe error"))
	}
	// 服务相关
	if tc := Client.Subscribe(edgeConfig.Topic.ReceiveTemperature, 0, receiveTemperature); tc.Wait() && tc.Error() != nil {
		log.Fatal(errors.Wrap(tc.Error(), "subscribe error"))
	}
}

type receiveTemperatureReq struct {
	Temperature float64 `json:"temperature"`
	PatientId   string  `json:"patient_id"`
}

func receiveTemperature(client mqtt.Client, msg mqtt.Message) {
	log.Printf("receive message %s", msg.Payload())
	log.Printf("from topic %s\n", msg.Topic())
	req := &receiveTemperatureReq{}
	if err := json.Unmarshal(msg.Payload(), req); err != nil {
		log.Fatal(errors.Wrap(err, "unmarshal payload error"))
	}

	if !operation.ReceiveTemperature(req.Temperature, req.PatientId) {
		return
	}

	success, abnormal := operation.IsAbnormal(req.PatientId)
	if !success {
		return
	}

	rpc.AddAbnormal(abnormal, req.PatientId)
}

func sync(client mqtt.Client, msg mqtt.Message) {
	log.Printf("receive message %s", msg.Payload())
	log.Printf("from topic %s\n", msg.Topic())
	if _, err := db.DB.Exec(string(msg.Payload())); err != nil {
		log.Fatal(errors.Wrap(err, "sync data error"))
	}
}

func PublishTicker() {
	ch := time.Tick(5 * time.Second)
	i := 0.0
	for range ch {
		log.Println("publish a temperature data")

		if tc := Client.Publish(config.Config.Edge.Topic.ReceiveTemperature, 0, false, fmt.Sprintf(`{"temperature": %f, "": "123"}`, 27.5+i)); tc.Wait() && tc.Error() != nil {
			log.Fatal(tc.Error())
		}
		i += 0.1
	}
}
