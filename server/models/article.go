package models

import (
	"github.com/jinzhu/gorm/dialects/postgres"
)

type Article struct {
	BaseModel
	Name                string         `json:"name" gorm:"column:name;type:text"`
	Abstract            string         `json:"abstract" gorm:"column:abstract;type:text;default:''"`
	Description         string         `json:"description" gorm:"column:description;type:text;default:''"`
	Rate                float64        `json:"rate" gorm:"type:double precision"`
	Tags                []Tag          `json:"tags" gorm:"many2many:article_tags;"`
	Authors             postgres.Jsonb `json:"authors" gorm:"type:jsonb;default:'[]'"`
	Links               postgres.Jsonb `json:"links" gorm:"type:jsonb;default:'[]'"`
	IntroducerReference int
	Introducer          User `json:"introducer" gorm:"foreignKey:IntroducerReference"` // 发布这篇文章的人
}
