package rpc

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
	pb "suvvm.work/ToadOCRTools/rpc/idl"
	"time"
)

var (
	successCode = flag.Int("success code", 0, "rpc reply code")
	serv = flag.String("service", "toad_ocr_preprocessor", "service name")
	reg  = flag.String("reg", "http://localhost:2379", "register etcd address")
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

func Process(netFlag string, image []byte) ([]string, error) {
	resp, err := toadOCRPreprocessorClient.Process(context.Background(),
		&pb.ProcessRequest{NetFlag: netFlag, Image: image})
	if err != nil {
		return nil, err
	}
	if resp.Code != int32(*successCode) {
		err = fmt.Errorf("resp code not success code:%v message:%v", resp.Code, resp.Message)
		return nil, err
	}
	return resp.Labels, nil
}
