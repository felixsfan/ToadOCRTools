package model

type AppInfo struct {
	ID int	`gorm:"column:id" json:"id"`
	Secret string `gorm:"column:secret" json:"secret"`
	Email string `gorm:"column:email" json:"email"`
	PNum string `gorm:"column:p_num" json:"p_num"`
}

type AppInfoReq struct {
	PNum string `json:"p_num"`
	Email string `json:"email"`
	UserVerifyCode string `json:"user_verify_code"`
	ClientVerifyCode string `json:"client_verify_code"`
}

type AppInfoResp struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	AppInfo *AppInfo `json:"app_info"`
}

func (req *AppInfoReq) Verify() bool {
	if req.UserVerifyCode == req.ClientVerifyCode {
		return true
	}
	return false
}

func (req *AppInfoReq) ToAppInfo() *AppInfo {
	return &AppInfo{
		PNum: req.PNum,
		Email: req.Email,
	}
}
