package service

import (
	"database/sql"
	"fmt"
	"log"
	"patient-edge/cloud/rpc/cloudclient"
	"patient-edge/entity"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/pkg/errors"
)

type UploadTemperature struct {
	*context
	// caller
	mqttClient           mqtt.Client //does not need
	uploadTemperatureRpc cloudclient.UploadTemperature
	// TempProperty
	currentPatient  *entity.Patient
	currentAbnormal *entity.Abnormal
}

func (p *UploadTemperature) createTemperature(value float64, patientId string) (bool, error) {

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
		return false, errors.New("precondition not satisfied")
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
	p.currentPatient = patient
	return true, nil
}

func (p *UploadTemperature) computeAbnormal() (bool, error) {
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
		p.currentAbnormal = abnormal
		return true, nil
	} else {
		return true, nil
	}
}

func (p *UploadTemperature) UploadTemperature(temperature float64, patientId string) error {
	if _, err := p.createTemperature(temperature, patientId); err != nil {
		return errors.Wrap(err, "createTemperature error")
	}
	if _, err := p.computeAbnormal(); err != nil {
		return errors.Wrap(err, "computeAbnormal error")
	}
	if p.currentAbnormal != nil {
		if _, err := p.uploadTemperatureRpc.UploadAbnormal(p.currentAbnormal); err != nil {
			return errors.Wrap(err, "uploadAbnormal error")
		}
	}
	return nil
}
