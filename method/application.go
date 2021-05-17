package method

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"suvvm.work/ToadOCRTools/dal/cluster"
	"suvvm.work/ToadOCRTools/dal/db"
	"suvvm.work/ToadOCRTools/model"
)

func VerifySecret(ctx context.Context, appID, basicToken, cntLen string) error {
	idStr, err := strconv.Atoi(appID)
	if err != nil {
		return fmt.Errorf("appID not int %v", err)
	}
	appInfo := &model.AppInfo{}
	appInfo.ID = idStr
	appSecret, err := cluster.GetKV(ctx, strconv.Itoa(appInfo.ID))
	if err != nil {
		appInfo, err = db.GetAppInfo(appInfo)
		if err != nil {
			return fmt.Errorf("appID not exists %v", err)
		}
		cluster.PutKV(ctx, strconv.Itoa(appInfo.ID), appSecret)
		appSecret = appInfo.Secret
	}
	hasher := md5.New()
	hasher.Write([]byte(appSecret + cntLen))
	md5Token := hex.EncodeToString(hasher.Sum(nil))
	fmt.Printf("md5Token:%v", md5Token)
	if  md5Token != basicToken {
		return fmt.Errorf("basic token incompatible")
	}
	return nil
}

func DoAddApplication(ctx context.Context, req *model.AppInfoReq) *model.AppInfoResp {
	reply := &model.AppInfoResp{}
	reply.Code = 0
	reply.Msg = "success"
	reply.AppInfo = req.ToAppInfo()
	if !req.Verify() {
		reply.Code = 1
		reply.Msg = "code is not same"
		return reply
	}
	appInfo, err := db.GetAppInfo(req.ToAppInfo())	// appInfo是否已经存在
	if err == nil {
		reply.Code = 1
		reply.Msg = "app info has already existed"
		reply.AppInfo = appInfo
		return reply
	}
	appInfo, err =  db.AddAppInfo(req.ToAppInfo())
	if err != nil {
		log.Printf("db.AddAppInfo err:%v", err)
		reply.Code = 1
		reply.Msg = "delete fail, please check the connection between the server and db"
		return reply
	}
	if err = cluster.PutKV(ctx, strconv.Itoa(appInfo.ID), appInfo.Secret); err != nil {
		log.Printf("cluster.PutKV err:%v", err)
		reply.Code = 2
		reply.Msg = "cluster cache kv failed. " +
			"this problem may cause slow server processing " +
			"but does not affect normal use."
		reply.AppInfo = appInfo
		return reply
	}
	reply.AppInfo = appInfo
	return reply
}

func DoDelApplication(ctx context.Context, req *model.AppInfoReq) *model.AppInfoResp {
	reply := &model.AppInfoResp{}
	reply.Code = 0
	reply.Msg = "success"
	reply.AppInfo = req.ToAppInfo()
	if !req.Verify() {
		reply.Code = 1
		reply.Msg = "code is not same"
		return reply
	}
	appInfo, err := db.GetAppInfo(req.ToAppInfo())
	reply.AppInfo = appInfo
	if err != nil {
		log.Printf("db.DelAppInfo err:%v", err)
		reply.Code = 1
		reply.Msg = "get fail, please check the param and the connection between the server and db"
		return reply
	}
	if err = cluster.DelKV(ctx, strconv.Itoa(appInfo.ID)); err != nil {
		log.Printf("cluster.DelKV err:%v", err)
		reply.Code = 1
		reply.Msg = "remove cluster cache kv failed."
		return reply
	}
	err =  db.DelAppInfo(req.ToAppInfo())
	if err != nil {
		log.Printf("db.AddAppInfo err:%v", err)
		reply.Code = 1
		reply.Msg = "add fail, please check the connection between the server and db"
		return reply
	}
	return reply
}

func DoGetApplication(ctx context.Context, req *model.AppInfoReq) *model.AppInfoResp {
	reply := &model.AppInfoResp{}
	reply.Code = 0
	reply.Msg = "success"
	reply.AppInfo = req.ToAppInfo()
	if !req.Verify() {
		reply.Code = 1
		reply.Msg = "code is not same"
		return reply
	}
	appInfo, err :=  db.GetAppInfo(req.ToAppInfo())
	if err != nil {
		log.Printf("db.AddAppInfo err:%v", err)
		reply.Code = 1
		reply.Msg = "get fail, please check the param and the connection between the server and db"
		return reply
	}
	reply.AppInfo = appInfo
	cluster.PutKV(ctx, strconv.Itoa(appInfo.ID), appInfo.Secret)
	return reply
}
