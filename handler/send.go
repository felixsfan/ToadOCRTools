package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"suvvm.work/ToadOCRTools/method"
	"suvvm.work/ToadOCRTools/model"
)

func Email(ctx *gin.Context) {

	//email, emailOk := ctx.GetPostForm("email")
	//code, codeOk := ctx.GetPostForm("code")
	//log.Printf("email:%v, code:%v", email, code)
	//if !codeOk || !emailOk {
	//	ctx.JSON(http.StatusBadRequest, gin.H{
	//		"code": 1,
	//		"message": "missing required parameters",
	//	})
	//	return
	//}
	var req model.EmailRequest
	if err := ctx.BindJSON(&req); err != nil {
		log.Printf("%v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 1,
			"message": "missing required parameters",
		})
		return
	}
	err := method.SendEmail(req.Email, req.Code)
	if err != nil {
		log.Printf("%v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 1,
			"message": "server mail sdk failure, please try again later",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
	})
}

func Sms(ctx *gin.Context) {
	//pNum, pNumOk := ctx.GetPostForm("p_num")
	//code, codeOk := ctx.GetPostForm("code")
	//log.Printf("p_num:%v, code:%v", pNum, code)
	//if !codeOk || !pNumOk{
	//	ctx.JSON(http.StatusBadRequest, gin.H{
	//		"code": 1,
	//		"message": "missing required parameters",
	//	})
	//	return
	//}
	var req model.SmsRequest
	if err := ctx.BindJSON(&req); err != nil {
		log.Printf("%v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 1,
			"message": "missing required parameters",
		})
		return
	}
	err := method.SendSms(req.PNum, req.Code)
	if err != nil {
		log.Printf("%v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 1,
			"message": "server sms sdk failure, please try again later",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
	})
}
