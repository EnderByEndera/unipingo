package models

type Tag struct {
	BaseModel
	Name        string `json:"name" gorm:"column:name;unique;type:text"`
	ABBR        string `json:"abbr" gorm:"column:abbr;type:text"`
	Description string `json:"description" gorm:"column:description;type:text"`
}
