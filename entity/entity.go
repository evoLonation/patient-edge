package entity

import "time"

/**
	for edge entity, need a attribute called edgeid to point to what edge node it store in.
	for entity with no id in requirement model, need a attribute called GenerateId(generate_id),
because sql need at least one primary key.
*/

type Temperature struct {
	// id
	GenerateId int32 `db:"generate_id"`
	// attribute
	Value     float64   `db:"value"`
	Timestamp time.Time `db:"timestamp"`
	// association
	PatientId string `db:"patient_id"`
	// edge
	EdgeId string `db:"edge_id"`
}

type Abnormal struct {
	// id
	GenerateId int32 `db:"generate_id"`
	// attribute
	Value     float64   `db:"value"`
	Timestamp time.Time `db:"timestamp"`
	// association
	PatientId string `db:"patient_id"`
}

type Patient struct {
	// id
	PatientId string `db:"patient_id"`
	// attribute

	// association
	DoctorId string `db:"doctor_id"`
}

type Doctor struct {
	// id
	DoctorId string `db:"doctor_id"`
}
