package db

import "time"

type Abnormal struct {
	AbnormalId int32     `db:"abnormal_id"`
	PatientId  string    `db:"patient_id"`
	Value      float64   `db:"value"`
	Timestamp  time.Time `db:"timestamp"`
}

type Patient struct {
	PatientId string `db:"patient_id"`
	EdgeId    string `db:"edge_id"`
	DoctorId  string `db:"doctor_id"`
}

type Doctor struct {
	DoctorId string `db:"doctor_id"`
}
type EdgeNode struct {
	EdgeId string `db:"edge_id"`
}
