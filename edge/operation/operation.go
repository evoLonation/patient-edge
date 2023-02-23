package operation

import (
	"database/sql"
	"fmt"
	"log"
	"patient-edge/edge/db"
	"time"

	"github.com/pkg/errors"
)

func ReceiveTemperature(value float64, patientId string) bool {
	DB := db.DB
	patient := db.Patient{}
	if err := DB.Get(&patient, "select * from patient where patient_id = ?", patientId); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("the patient id %s does not exists\n", patientId)
			return false
		} else {
			log.Fatal(errors.Wrap(err, "get patient error"))
		}
	}
	temperature := &db.Temperature{
		PatientId: patientId,
		Value:     value,
		Timestamp: time.Now(),
	}
	if _, err := DB.NamedExec("insert into temperature (patient_id, value, timestamp) values (:patient_id, :value, :timestamp)", temperature); err != nil {
		log.Fatal(errors.Wrap(err, "insert temperature data error"))
	}
	return true
}

func IsAbnormal(patientId string) (bool, float64) {
	DB := db.DB
	// 获取数据库的历史数据
	var temperatures []db.Temperature
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
	return avg >= 27, avg
}
