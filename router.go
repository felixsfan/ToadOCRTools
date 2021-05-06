package main

import (
	"github.com/gin-gonic/gin"
	"suvvm.work/ToadOCRTools/handler"
)

func register(r *gin.Engine) {
	r.Any("/toad_ocr/ping", handler.Pong)
	r.Any("/toad_ocr/process", handler.Process)
}
