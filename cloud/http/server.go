package http

import (
	"net/http"
	"patient-edge/cloud/operation"
	"patient-edge/config"

	"github.com/gin-gonic/gin"
)

type addPatientReq struct {
	PatientId string `json:"patient_id"`
	DoctorId  string `json:"doctor_id"`
	EdgeId    string `json:"edge_id"`
}
type addDoctorReq struct {
	DoctorId string `json:"doctor_id"`
}
type Res struct {
	Result bool `json:"result"`
}

// 该函数返回一个gin.H，gin.H是一个map，存储着键值对，将要返回给请求者
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func Start() {
	router := gin.Default()
	router.POST("/add-patient", addPatient)
	router.POST("/add-doctor", addDoctor)
	go func() {
		router.Run(":" + config.Config.Cloud.HttpServerPort)
	}()
}

func addPatient(ctx *gin.Context) {
	var req addPatientReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		//证明请求对于该结构体并不有效
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ret := operation.AddPatient(req.PatientId, req.DoctorId, req.EdgeId)
	ctx.JSON(http.StatusOK, &Res{Result: ret})
}

func addDoctor(ctx *gin.Context) {
	var req addDoctorReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		//证明请求对于该结构体并不有效
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ret := operation.AddDoctor(req.DoctorId)
	ctx.JSON(http.StatusOK, &Res{Result: ret})
}
