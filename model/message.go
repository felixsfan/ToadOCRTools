package model

type EmailRequest struct {
	Email string `json:"email"`
	Code string `json:"code"`
}

type SmsRequest struct {
	PNum string `json:"p_num"`
	Code string `json:"code"`
}
