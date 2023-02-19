package models

type User struct {
	BaseModel
	Name         string `json:"name" gorm:"column:name;type:text"`
	EMail        string `json:"email" gorm:"column:email;"`
	PasswordHash string `gorm:"column:password_hash"`
}
