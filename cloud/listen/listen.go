package listen

// 1、用于处理请求、相应式消息（web服务器）
// 2、用于处理订阅、发布式（mqttClient）

import (
	"patient-edge/cloud/service"
	"patient-edge/common"
	"patient-edge/config"
	"patient-edge/entity"

	"github.com/gin-gonic/gin"
)

type addPatientReq struct {
	PatientId string `json:"patient_id"`
	DoctorId  string `json:"doctor_id"`
}

type addDoctorReq struct {
	DoctorId string `json:"doctor_id"`
}

type getTemperatureReq struct {
	PatientId string `json:"patient_id"`
	Number    int    `json:"number"`
}
type getTemperatureRes struct {
	Temperature []entity.Temperature `json:"temperatures"`
}

type res struct {
	Result bool `json:"result"`
}

func Start(config config.ListenConf, patientSvc *service.Patient, doctorSvc *service.Doctor) {
	router := gin.Default()
	router.POST("patient/add-patient", common.GetGinHandler(func(req *addPatientReq) (*res, error) {
		ret, err := patientSvc.AddPatient(req.PatientId, req.DoctorId)
		return &res{Result: ret}, err
	}))
	router.POST("doctor/add-doctor", common.GetGinHandler(func(req *addDoctorReq) (*res, error) {
		ret, err := doctorSvc.AddDoctor(req.DoctorId)
		return &res{Result: ret}, err
	}))
	router.POST("doctor/get-temperature", common.GetGinHandler(func(req *getTemperatureReq) (*getTemperatureRes, error) {
		ret, err := doctorSvc.GetTemperature(req.PatientId, req.Number)
		return &getTemperatureRes{Temperature: ret}, err
	}))
	go func() {
		router.Run(":" + config.Http.Port)
	}()
}
