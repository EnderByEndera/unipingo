package models

type User struct {
	BaseModel
	Name          string `json:"name" gorm:"column:name;type:text"`
	EMail         string `json:"email" gorm:"column:email;"`
	PasswordHash  string `json:"-" gorm:"column:password_hash"`
	WechatOpenID  string `json:"wechatOpenID" gorm:"column:wechat_openid"`
	WechatUnionID string `json:"wechatUnionID" gorm:"column:wechat_unionid"`
}

type LoginResponse struct {
	UserInfo User   `json:"user"`
	JWTToken string `json:"jwtToken"`
}
