package service

import "patient-edge/entity"

type Doctor struct {
	*context
}

func (p *Doctor) AddDoctor(doctorId string) (bool, error) {
	return true, nil
}

func (p *Doctor) GetTemperature(patientId string, number int) ([]entity.Temperature, error) {
	return nil, nil
}
