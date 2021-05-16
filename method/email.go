package method

import (
	"github.com/aliyun-sdk/mail-go/smtp"
	"suvvm.work/ToadOCRTools/config"
)

var (
	client *smtp.Client
)

func InitEmail() {
	client = smtp.New("smtpdm.aliyun.com:80",
		config.AppConfig.SdkConfig.Email, config.AppConfig.SdkConfig.SmtpPsw)
}

// SendEmail 发送验证邮件
func SendEmail(email, code string ) error {
	smtp.SendTo()
	err := client.Send(
		smtp.From("ToadOCRDevTeam"),
		smtp.Subject("ToadOCRMessage"),
		smtp.SendTo(email),
		smtp.Content(smtp.Plain, "Hello," +
			email + ".We have received your application for " +
			"ToadOCR, here is your verification code " +
			code +
			"."),
	)
	if err != nil {
		return err
	}
	return nil
}
