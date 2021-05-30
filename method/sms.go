package method

import (
	"ToadOCRTools/config"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"log"
)

// SendSms 发送短信验证码
func SendSms(pNum, code string) error {
	client, err := dysmsapi.NewClientWithAccessKey("cn-qingdao",
		config.AppConfig.SdkConfig.AppID, config.AppConfig.SdkConfig.AppSecret)
	if err != nil {
		log.Print(err.Error())
		return err
	}
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = pNum
	request.SignName = "suvvm"
	request.TemplateCode = "SMS_216838466"
	request.TemplateParam = "{\"code\":\"" +
		code +
		"\"}"
	response, err := client.SendSms(request)
	if err != nil {
		log.Print(err.Error())
		return err
	}
	log.Printf("response is %#v\n", response)
	return nil
}
