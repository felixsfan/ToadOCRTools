package rpc

import (
	"ToadOCRTools/dal/cluster"
	"ToadOCRTools/dal/db"
	"ToadOCRTools/model"
	pb "ToadOCRTools/rpc/idl"
	"context"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
	"log"
	"strconv"
	"time"
)

var (
	successCode               = flag.Int("success code", 0, "rpc reply code")
	serv                      = flag.String("service", "toad_ocr_preprocessor", "service name")
	reg                       = flag.String("reg", "http://localhost:2379", "register etcd address")
	toadOCRPreprocessorClient pb.ToadOcrPreprocessorClient
)

func init() {
	flag.Parse()
	r := NewResolver(*reg, *serv)
	resolver.Register(r)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	conn, err := grpc.DialContext(ctx, r.Scheme()+"://authority/"+*serv,
		grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	toadOCRPreprocessorClient = pb.NewToadOcrPreprocessorClient(*conn)
}

func Process(ctx context.Context, netFlag, appID string, image []byte) ([]string, error) {
	req := &pb.ProcessRequest{
		AppId:   appID,
		NetFlag: netFlag,
		Image:   image,
	}
	token, err := GetBasicToken(ctx, req.AppId, req.NetFlag+strconv.Itoa(len(req.Image)))
	if err != nil {
		return nil, err
	}
	req.BasicToken = token
	resp, err := toadOCRPreprocessorClient.Process(context.Background(), req)
	if err != nil {
		return nil, err
	}
	if resp.Code != int32(*successCode) {
		err = fmt.Errorf("resp code not success code:%v message:%v", resp.Code, resp.Message)
		return nil, err
	}
	return resp.Labels, nil
}

func GetBasicToken(ctx context.Context, appID, verifyStr string) (string, error) {
	hasher := md5.New()
	secret, err := cluster.GetKV(ctx, appID)
	if err != nil {
		idInt, errInner := strconv.Atoi(appID)
		if errInner != nil {
			log.Printf("AppID:%v not int", appID)
			return "", err
		}
		appInfo, errInner := db.GetAppInfo(&model.AppInfo{ID: idInt})
		if errInner != nil {
			log.Printf("db get app err:%v", errInner)
			return "", err
		}
		secret = appInfo.Secret
		cluster.PutKV(ctx, appID, secret)
	}
	//imglen := strconv.Itoa(len(req.Image))
	//hasher.Write([]byte(secret + req.NetFlag + imglen))
	hasher.Write([]byte(secret + verifyStr))
	md5Token := hex.EncodeToString(hasher.Sum(nil))
	return md5Token, nil
}
