package main

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	IP             = "localhost:1883"
	ClientId       = "PatientChecker-01"
	SubscribeTopic = "$patient/sensor/+/temperature"
	DB             = "root:2002116yy@tcp(edge-mysql:3306)/patient_edge?parseTime=true"
)

var Db *sqlx.DB

type TemperatureSensor struct {
	Temperature float64 `json:"temperature"`
}

type TemperatureDB struct {
	TemperatureId int32     `db:"temperature_id"`
	PatientId     string    `db:"patient_id"`
	Value         float64   `db:"temperature"`
	Timestamp     time.Time `db:"timestamp"`
}
type Patient struct {
	PatientId string `db:"patient_id"`
}

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

func onTemperatureMessage(client mqtt.Client, message mqtt.Message) {
	log.Printf("receive message %s", message.Payload())
	// 解析温度数据
	temperatureSensor := &TemperatureSensor{}
	if err := json.Unmarshal(message.Payload(), temperatureSensor); err != nil {
		log.Fatal(err)
	}
	// 从topic解析出病人id
	patientId := strings.Split(message.Topic(), "/")[2]
	var patient *Patient
	if err := Db.Get(patient, "select * from patient where patient_id = "+patientId); err != nil {
		log.Fatal(err)
	}
	if patient == nil {
		return
	}
	temperature := &TemperatureDB{
		PatientId: patientId,
		Value:     temperatureSensor.Temperature,
		Timestamp: time.Now(),
	}
	if _, err := Db.NamedExec("insert into temperature (patient_id, value, timestamp) values (:patient_id, :value, :timestamp)", temperature); err != nil {
		log.Fatal(err)
	}
}

func checkPatient(patientId string) {

}

func main() {
	var err error
	Db, err = sqlx.Open("mysql", DB)
	if err != nil {
		log.Fatal(err)
	}

	client, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	if tc := client.Subscribe(SubscribeTopic, 0, onTemperatureMessage); tc.Wait() && tc.Error() != nil {
		log.Fatal(tc.Error())
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
