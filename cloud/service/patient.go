package service

type Patient struct {
	*context
}

func (p *Patient) AddPatient(patientId string, doctorId string) (bool, error) {
	return true, nil
}
