// package operation

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"patient-edge/common/config"
// 	"patient-edge/edge/entity"
// 	"sync"
// 	"time"

// 	"github.com/jmoiron/sqlx"
// 	"github.com/pkg/errors"
// )

// var Op Operation

// type Operation struct {
// 	db              *sqlx.DB
// 	CurrentPatient  *entity.Patient
// 	CurrentAbnormal float64
// }

// func init() {
// 	conf := config.Config.Edge.Operation
// 	var err error
// 	Op.db, err = sqlx.Open("mysql", conf.DataSource)
// 	if err != nil {
// 		log.Fatal(errors.Wrap(err, "open db error"))
// 	}
// }

// /*
// precondition: 存在目标patientId的patient
// postcondition: 增加patient关联的温度数据
// */
// func (p *Operation) ReceiveTemperature(value float64, patientId string) error {
// 	log.Println("operation ReceiveTemperature start")
// 	defer log.Println("operation ReceiveTemperature done")
// 	wg := sync.WaitGroup{}
// 	log.Printf("definition start")
// 	var patient entity.Patient
// 	var patientUndefined bool

// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		if err := p.db.Get(&patient, "select * from patient where patient_id = ?", patientId); err != nil {
// 			if err == sql.ErrNoRows {
// 				log.Printf("the patient id %s does not exists\n", patientId)
// 				patientUndefined = true
// 				return
// 			} else {
// 				log.Fatal(errors.Wrap(err, "get patient error"))
// 			}
// 		}
// 		patientUndefined = false
// 	}()
// 	wg.Wait()
// 	log.Printf("definition done")
// 	if !(patientUndefined == false) {
// 		return errors.New("precondition not satisfied")
// 	}
// 	log.Println("pre condition satisfied")
// 	log.Println("post condition start")
// 	temperature := entity.Temperature{
// 		PatientId: patientId,
// 		Value:     value,
// 		Timestamp: time.Now(),
// 	}
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		if _, err := p.db.NamedExec("insert into temperature (patient_id, value, timestamp) values (:patient_id, :value, :timestamp)", &temperature); err != nil {
// 			log.Fatal(errors.Wrap(err, "insert temperature data error"))
// 		}
// 	}()
// 	wg.Wait()
// 	p.CurrentPatient = &patient
// 	return nil
// }

// /*
// precondition: currentPatient存在
// postcondition: 根据patient的温度数据，计算得到currentabnormal（如果有的话）
// */
// func (p *Operation) IsAbnormal() (bool, error) {
// 	log.Println("operation IsAbnormal start")
// 	defer log.Println("operation IsAbnormal done")
// 	wg := sync.WaitGroup{}

// 	if p.CurrentPatient == nil {
// 		return false, errors.New("precondition not satisfied")
// 	}
// 	log.Println("pre condition satisfied")
// 	log.Printf("definition start")
// 	var temperatureSet []entity.Temperature
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		if err := p.db.Select(&temperatureSet, "select * from temperature where patient_id = ? order by timestamp desc limit 5"); err != nil {
// 			log.Fatal(errors.Wrap(err, "get temperature data error"))
// 		}
// 	}()
// 	wg.Wait()
// 	log.Printf("definition done")
// 	log.Println("post condition start")
// 	var valuesStr string
// 	var avg float64
// 	for _, temperature := range temperatureSet {
// 		avg += temperature.Value
// 		valuesStr += fmt.Sprintf("%f, ", temperature.Value)
// 	}
// 	log.Printf("temperature data: %s\n", valuesStr)
// 	avg /= float64(len(temperatureSet))
// 	log.Printf("average temperature is %f\n", avg)
// 	if avg >= 27 {
// 		p.CurrentAbnormal = avg
// 		return true, nil
// 	} else {
// 		return false, nil
// 	}
// }
