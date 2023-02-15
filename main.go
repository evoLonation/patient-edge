package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"strings"
	"time"

	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	IP             = "edge-mqtt:1883"
	ClientId       = "PatientChecker-01"
	SubscribeTopic = "$patient/sensor/+/temperature"
	DBResource     = "root:2002116yy@tcp(edge-mysql:3306)/patient_edge?parseTime=true"
)

var DB *sqlx.DB

type TemperatureSensor struct {
	Temperature float64 `json:"temperature"`
}

type TemperatureDB struct {
	TemperatureId int32     `db:"temperature_id"`
	PatientId     string    `db:"patient_id"`
	Value         float64   `db:"value"`
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
	log.Printf("from topic %s\n", message.Topic())
	// 解析温度数据
	temperatureSensor := &TemperatureSensor{}
	if err := json.Unmarshal(message.Payload(), temperatureSensor); err != nil {
		log.Fatal(errors.Wrap(err, "unmarshal payload error"))
	}
	// 从topic解析出病人id
	patientId := strings.Split(message.Topic(), "/")[2]
	patient := &Patient{}
	if err := DB.Get(patient, "select * from patient where patient_id = ?", patientId); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("the patient id %s from does not exists\n", patientId)
			return
		} else {
			log.Fatal(errors.Wrap(err, "get patient error"))
		}
	}
	temperature := &TemperatureDB{
		PatientId: patientId,
		Value:     temperatureSensor.Temperature,
		Timestamp: time.Now(),
	}
	if _, err := DB.NamedExec("insert into temperature (patient_id, value, timestamp) values (:patient_id, :value, :timestamp)", temperature); err != nil {
		log.Fatal(errors.Wrap(err, "insert temperature data error"))
	}
	// todo 直接调用耦合性太高，是否应该改为订阅\发布调用？
	checkPatient(patientId)
}

func checkPatient(patientId string) {

}

func main() {
	var err error
	DB, err = sqlx.Open("mysql", DBResource)
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
