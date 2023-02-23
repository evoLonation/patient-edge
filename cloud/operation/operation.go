package operation

import (
	"database/sql"
	"log"
	"patient-edge/cloud/db"
	"patient-edge/cloud/mqtt"
	"time"

	"github.com/pkg/errors"
)

func ReceiveAbnormal(value float64, patientId string) bool {
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
	abnormal := &db.Abnormal{
		PatientId: patientId,
		Value:     value,
		Timestamp: time.Now(),
	}
	if _, err := DB.NamedExec("insert into abnormal (patient_id, value, timestamp) values (:patient_id, :value, :timestamp)", abnormal); err != nil {
		log.Fatal(errors.Wrap(err, "insert temperature data error"))
	}
	mqtt.Notice(value, patientId, patient.DoctorId)
	return true
}

func AddPatient(patientId string, doctorId string, edgeId string) bool {
	DB := db.DB
	patient := db.Patient{}
	if err := DB.Get(&patient, "select * from patient where patient_id = ?", patientId); err == nil || err != sql.ErrNoRows {
		if err == nil {
			log.Printf("the patient id %s already exists\n", patientId)
			return false
		}
		log.Fatal(errors.Wrap(err, "get patient error"))
	}
	patient = db.Patient{
		DoctorId:  doctorId,
		PatientId: patientId,
		EdgeId:    edgeId,
	}
	if _, err := DB.NamedExec("insert into patient (patient_id, edge_id, doctor_id) values  (:patient_id, :edge_id, :doctor_id)", &patient); err != nil {
		log.Fatal(errors.Wrap(err, "insert new patient error"))
	}
	return true
}

func AddDoctor(doctorId string) bool {
	DB := db.DB
	doctor := db.Doctor{}
	if err := DB.Get(&doctor, "select * from doctor where doctor_id = ?", doctor); err == nil || err != sql.ErrNoRows {
		if err == nil {
			log.Printf("the doctor id %s already exists\n", doctorId)
			return false
		}
		log.Fatal(errors.Wrap(err, "get doctor error"))
	}
	doctor = db.Doctor{
		DoctorId: doctorId,
	}
	if _, err := DB.NamedExec("insert into doctor (doctor_id) values  (:doctor_id)", &doctor); err != nil {
		log.Fatal(errors.Wrap(err, "insert new doctor error"))
	}
	return true
}
