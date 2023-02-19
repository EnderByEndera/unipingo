package db

import (
	"log"
	"melodie-site/server/auth"
	"melodie-site/server/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

type Users struct {
	Id        int       `json:"id" gorm:"column:id;unique;primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"column:name"`
	Age       int       `json:"age" gorm:"column:age"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

// 定义表名
func (Users) TableName() string {
	return "users"
}

var database *gorm.DB

func GetDBConn() *gorm.DB {
	return database
}

func InitDB() {
	dsn := "host=127.0.0.1 port=5432 user=melodie dbname=postgres password=melodie-test sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	database = db
	if err != nil {
		log.Println(err)
	}
	db.AutoMigrate(&Users{})
	db.AutoMigrate(&models.Article{})
	db.AutoMigrate(&models.Tag{})
	db.AutoMigrate(&models.User{})
	user := &models.User{Name: "admin"}
	err = db.Model(&models.User{}).First(user).Error
	if err != nil {
		db.Create(&models.User{Name: "admin", EMail: "1295752786@qq.com", PasswordHash: auth.EncryptPassword("123456")})
	}
}
