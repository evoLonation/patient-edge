package mqtt

import (
	"encoding/json"
	"log"
	"patient-edge/config"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
)

var Client mqtt.Client

func init() {
	log.Println("init mqtt client")
	cloudConfig := &config.Config.Cloud
	Client = mqtt.NewClient(mqtt.NewClientOptions().AddBroker(cloudConfig.MqttBroker).SetClientID(cloudConfig.ClientId))
	if tc := Client.Connect(); tc.Wait() && tc.Error() != nil {
		log.Fatal(errors.Wrap(tc.Error(), "connect to broker error"))
	}
	log.Println("init mqtt done")
}

type NoticeMsg struct {
	PatientId string  `json:"patient_id"`
	Abnormal  float64 `json:"abnormal"`
}

func Notice(abnormal float64, patientId string, doctorId string) {
	msg := NoticeMsg{
		PatientId: patientId,
		Abnormal:  abnormal,
	}
	topic := strings.Replace(config.Config.Cloud.Topic.Notice, "+", doctorId, 1)
	msgb, err := json.Marshal(msg)
	if err != nil {
		log.Fatal(errors.Wrap(err, "marshal msg error"))
	}
	if tc := Client.Publish(topic, 0, false, msgb); tc.Wait() && tc.Error() != nil {
		log.Fatal(errors.Wrap(tc.Error(), "publish error"))
	}
}

func Sync(query string, edgeId string) {
	topic := strings.Replace(config.Config.Common.Topic.Sync, "+", edgeId, 1)
	log.Printf("publish to topic: %s \n", topic)
	if tc := Client.Publish(topic, 0, false, query); tc.Wait() && tc.Error() != nil {
		log.Fatal(errors.Wrap(tc.Error(), "publish sync error"))
	}
}
