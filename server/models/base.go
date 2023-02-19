package models

import "time"

type ModelID struct {
	ID int `json:"id" gorm:"column:id;unique;primaryKey;autoIncrement"`
}

type BaseModel struct {
	ModelID
	CreatedAt time.Time `json:"-" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"-" gorm:"column:updated_at"`
}
