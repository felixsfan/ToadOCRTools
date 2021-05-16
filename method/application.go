package method

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"suvvm.work/ToadOCRTools/dal/db"
	"suvvm.work/ToadOCRTools/model"
)

func VerifySecret(appID, basicToken, cntLen string) error {
	idStr, err := strconv.Atoi(appID)
	if err != nil {
		return fmt.Errorf("appID not int %v", err)
	}
	appInfo := &model.AppInfo{}
	appInfo.ID = idStr
	appInfo, err = db.GetAppInfo(appInfo)
	if err != nil {
		return fmt.Errorf("appID not exists %v", err)
	}
	hasher := md5.New()
	hasher.Write([]byte(appInfo.Secret + cntLen))
	md5Token := hex.EncodeToString(hasher.Sum(nil))
	fmt.Printf("md5Token:%v", md5Token)
	if  md5Token != basicToken {
		return fmt.Errorf("basic token incompatible")
	}
	return nil
}

func DoAddApplication(req *model.AppInfoReq) *model.AppInfoResp {
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
	reply.AppInfo = appInfo
	return reply
}

func DoDelApplication(req *model.AppInfoReq) *model.AppInfoResp {
	reply := &model.AppInfoResp{}
	reply.Code = 0
	reply.Msg = "success"
	reply.AppInfo = req.ToAppInfo()
	if !req.Verify() {
		reply.Code = 1
		reply.Msg = "code is not same"
		return reply
	}
	err :=  db.DelAppInfo(req.ToAppInfo())
	if err != nil {
		log.Printf("db.AddAppInfo err:%v", err)
		reply.Code = 1
		reply.Msg = "add fail, please check the connection between the server and db"
		return reply
	}
	return reply
}

func DoGetApplication(req *model.AppInfoReq) *model.AppInfoResp {
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
	return reply
}
