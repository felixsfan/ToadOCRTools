package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"suvvm.work/ToadOCRTools/common"
	"suvvm.work/ToadOCRTools/dal/cluster"
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
	ctx.JSON(http.StatusOK, method.DoAddApplication(ctx, req))
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
	ctx.JSON(http.StatusOK, method.DoDelApplication(ctx, req))
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
	ctx.JSON(http.StatusOK, method.DoGetApplication(ctx, req))
}

func ApplicationCache(ctx *gin.Context) {
	reply := &model.AppInfoResp{}
	appID := ctx.DefaultQuery("app_id","")
	if appID == "" {
		reply.Code = common.HandlerReadBodyErr
		reply.Msg = common.HandlerReadBodyErrMsg
		ctx.JSON(http.StatusBadRequest, reply)
		return
	}
	value, err := cluster.GetKV(ctx, appID)
	if err != nil {
		reply.Code = 1
		reply.Msg = "no cash found"
		ctx.JSON(http.StatusBadRequest, reply)
		return
	}
	idInt, err  := strconv.Atoi(appID)
	if err != nil {
		reply.Code = 1
		reply.Msg = "app id not int"
		ctx.JSON(http.StatusBadRequest, reply)
		return
	}
	reply.Code = 0
	reply.Msg = "success"
	reply.AppInfo = &model.AppInfo{
		ID: idInt,
		Secret: value,
	}
	ctx.JSON(http.StatusOK, reply)
}
