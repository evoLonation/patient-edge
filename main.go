package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	IP = "edge-mqtt:1883"
	// IP               = "localhost:1883"
	ClientId         = "PatientChecker-01"
	TemperatureTopic = "$patient/sensor/sensor-temperature-01/temperature"
	CheckTopic       = "$patient/edge/sensor-temperature-01/check"
	DBResource       = "root:2002116yy@tcp(edge-mysql:3306)/patient_edge?parseTime=true"
	CloudDBResource  = "root:2002116yy@tcp(cloud-mysql:3306)/patient_cloud?parseTime=true"
)

var DB *sqlx.DB
var CDB *sqlx.DB

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

type Abnormal struct {
	AbnormalId int32     `db:"abnormal_id"`
	PatientId  string    `db:"patient_id"`
	Value      float64   `db:"value"`
	Timestamp  time.Time `db:"timestamp"`
}

var Client mqtt.Client

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
	patient := Patient{}
	if err := DB.Get(&patient, "select * from patient where patient_id = ?", patientId); err != nil {
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
	go func() {
		if tc := Client.Publish(CheckTopic, 0, false, []byte("")); tc.Wait() && tc.Error() != nil {
			log.Fatal(errors.Wrap(tc.Error(), "publish check topic error"))
		}
	}()
}

func onCheckMessage(client mqtt.Client, message mqtt.Message) {
	log.Printf("receive message %s", message.Payload())
	log.Printf("from topic %s\n", message.Topic())
	patientId := strings.Split(message.Topic(), "/")[2]
	// 获取数据库的历史数据
	var temperatures []TemperatureDB
	if err := DB.Select(&temperatures, "select * from temperature where patient_id = ? order by timestamp desc limit 5", patientId); err != nil {
		log.Fatal(errors.Wrap(err, "get temperature data error"))
	}
	var valuesStr string
	var avg float64
	for _, temperature := range temperatures {
		avg += temperature.Value
		valuesStr += fmt.Sprintf("%f, ", temperature.Value)
	}
	log.Printf("temperature data: %s\n", valuesStr)
	avg /= float64(len(temperatures))
	log.Printf("average temperature is %f\n", avg)
	if avg >= 28 {
		log.Printf("patient %s is abnormal, average temperature is %f\n", patientId, avg)
		abnormal := &Abnormal{
			PatientId: patientId,
			Value:     avg,
			Timestamp: time.Now(),
		}
		if _, err := CDB.NamedExec("insert into abnormal (patient_id, value, timestamp) values (:patient_id, :value, :timestamp)", abnormal); err != nil {
			log.Fatal(errors.Wrap(err, "insert abnormal data error"))
		}
	}
}

func publishTicker() {
	ch := time.Tick(5 * time.Second)
	i := 0.0
	for range ch {
		log.Println("publish a temperature data")

		if tc := Client.Publish(TemperatureTopic, 0, false, fmt.Sprintf(`{"temperature": %f}`, 27.5+i)); tc.Wait() && tc.Error() != nil {
			log.Fatal(tc.Error())
		}
		i += 0.1
	}
}

func main() {
	var err error
	DB, err = sqlx.Open("mysql", DBResource)
	if err != nil {
		log.Fatal(err)
	}
	CDB, err = sqlx.Open("mysql", CloudDBResource)
	if err != nil {
		log.Fatal(err)
	}

	Client, err = connect()
	if err != nil {
		log.Fatal(err)
	}
	if tc := Client.Subscribe(TemperatureTopic, 0, onTemperatureMessage); tc.Wait() && tc.Error() != nil {
		log.Fatal(tc.Error())
	}
	if tc := Client.Subscribe(CheckTopic, 0, onCheckMessage); tc.Wait() && tc.Error() != nil {
		log.Fatal(tc.Error())
	}
	log.Println("init done, start work")
	// publishTicker()
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
