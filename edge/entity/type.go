package entity

import "time"

type Temperature struct {
	TemperatureId int32     `db:"temperature_id"`
	PatientId     string    `db:"patient_id"`
	Value         float64   `db:"value"`
	Timestamp     time.Time `db:"timestamp"`
}

type Patient struct {
	PatientId string `db:"patient_id"`
}
