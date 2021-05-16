package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"suvvm.work/ToadOCRTools/common"
	"suvvm.work/ToadOCRTools/method"
	"suvvm.work/ToadOCRTools/model"
)

func ApplicationAdd(ctx *gin.Context) {
	req := &model.AppInfoReq{}
	reply := &model.AppInfoResp{}
	if err := ctx.BindJSON(req); err != nil {
		reply.Code = common.HandlerReadBodyErr
		reply.Msg = common.HandlerReadBodyErrMsg
		reply.AppInfo = req.ToAppInfo()
		ctx.JSON(http.StatusForbidden, reply)
		return
	}
	ctx.JSON(http.StatusOK, method.DoAddApplication(req))
}

func ApplicationDel(ctx *gin.Context) {
	req := &model.AppInfoReq{}
	reply := &model.AppInfoResp{}
	if err := ctx.BindJSON(req); err != nil {
		reply.Code = common.HandlerReadBodyErr
		reply.Msg = common.HandlerReadBodyErrMsg
		reply.AppInfo = req.ToAppInfo()
		ctx.JSON(http.StatusForbidden, reply)
		return
	}
	ctx.JSON(http.StatusOK, method.DoDelApplication(req))
}

func ApplicationGet(ctx *gin.Context) {
	req := &model.AppInfoReq{}
	reply := &model.AppInfoResp{}
	req.PNum = ctx.DefaultQuery("p_num","")
	req.Email = ctx.DefaultQuery("email","")
	//req.UserVerifyCode = ctx.DefaultQuery("","")
	//req.ClientVerifyCode = ctx.DefaultQuery("","")
	if req.Email == "" && req.PNum == "" {
		reply.Code = common.HandlerReadBodyErr
		reply.Msg = common.HandlerReadBodyErrMsg
		reply.AppInfo = req.ToAppInfo()
		ctx.JSON(http.StatusForbidden, reply)
		return
	}
	ctx.JSON(http.StatusOK, method.DoGetApplication(req))
}
