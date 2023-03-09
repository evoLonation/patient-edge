package operation

import (
	"database/sql"
	"fmt"
	"log"
	"patient-edge/common"
	"patient-edge/config"
	"patient-edge/entity"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type UploadTemperatureService struct {
	// config is private
	edgeDB     *sqlx.DB
	cloudDB    *sqlx.DB
	redis      string // todo
	mqttClient mqtt.Client
	// TempProperty is public
	CurrentPatient  *entity.Patient
	CurrentAbnormal *entity.Abnormal
}

func NewUploadTemperatureService(conf config.EdgeServiceConf) *UploadTemperatureService {
	return &UploadTemperatureService{
		edgeDB:     common.GetMysqlDB(conf.EdgeDataSource),
		cloudDB:    common.GetMysqlDB(conf.CloudDataSource),
		mqttClient: common.GetMqttClient(conf.Mqtt),
	}
}

func (p *UploadTemperatureService) CreateTemperature(value float64, patientId string) error {

	log.Println("operation CreateTemperature start")
	defer log.Println("operation ReceiveTemperature done")

	wg := sync.WaitGroup{}

	log.Printf("definition start")
	patient := &entity.Patient{}
	var patientUndefined bool

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := p.edgeDB.Get(patient, "select * from patient where patient_id = ?", patientId); err != nil {
			if err == sql.ErrNoRows {
				log.Printf("the patient id %s does not exists\n", patientId)
				patientUndefined = true
				return
			} else {
				log.Fatal(errors.Wrap(err, "get patient error"))
			}
		}
		patientUndefined = false
	}()
	wg.Wait()

	log.Printf("precondition start")
	if !(patientUndefined == false) {
		return errors.New("precondition not satisfied")
	}

	log.Println("post condition start")
	temperature := &entity.Temperature{
		PatientId: patientId,
		Value:     value,
		Timestamp: time.Now(),
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if _, err := p.edgeDB.NamedExec("insert into temperature (patient_id, value, timestamp) values (:patient_id, :value, :timestamp)", temperature); err != nil {
			log.Fatal(errors.Wrap(err, "insert temperature data error"))
		}
	}()
	wg.Wait()
	p.CurrentPatient = patient
	return nil
}

func (p *UploadTemperatureService) ComputeAbnormal() error {
	log.Println("operation CreateTemperature start")
	defer log.Println("operation ReceiveTemperature done")

	wg := sync.WaitGroup{}

	log.Printf("definition start")

	var temperatureSet []entity.Temperature
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := p.edgeDB.Select(&temperatureSet, "select * from temperature where patient_id = ? order by timestamp desc limit 5"); err != nil {
			log.Fatal(errors.Wrap(err, "get temperature data error"))
		}
	}()
	wg.Wait()

	log.Printf("precondition start")

	log.Println("post condition start")

	var valuesStr string
	var avg float64
	for _, temperature := range temperatureSet {
		avg += temperature.Value
		valuesStr += fmt.Sprintf("%f, ", temperature.Value)
	}
	log.Printf("temperature data: %s\n", valuesStr)
	avg /= float64(len(temperatureSet))
	log.Printf("average temperature is %f\n", avg)
	if avg >= 27 {
		abnormal := &entity.Abnormal{
			Value:     avg,
			Timestamp: time.Now(),
		}
		p.CurrentAbnormal = abnormal
		return nil
	} else {
		return nil
	}
}
