package main

import (
	"ToadOCRTools/handler"
	"github.com/gin-gonic/gin"
)

func register(r *gin.Engine) {
	r.Any("/toad_ocr/ping", handler.Pong)

	r.POST("/toad_ocr/process", handler.Process)
	r.POST("/toad_ocr/process/v2", handler.ProcessV2)

	r.POST("/toad_ocr/send/sms", handler.Sms)
	r.POST("/toad_ocr/send/email", handler.Email)

	r.POST("/toad_ocr/application", handler.ApplicationAdd)
	// r.PUT("/toad_ocr/application", handler.Process)
	r.DELETE("/toad_ocr/application", handler.ApplicationDel)
	r.GET("/toad_ocr/application", handler.ApplicationGet)

	r.GET("/toad_ocr/application/cache", handler.ApplicationCache)
}
