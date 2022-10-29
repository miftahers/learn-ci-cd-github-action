package config

import (
	"fmt"
	"os"
	"praktikum/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}

	err = MigrateDB(db)
	if err != nil {
		panic(err)
	}
	return db
}

var (
	DB_Address = os.Getenv("DB_ADDRESS")
	DB_Name    = os.Getenv("DB_NAME")
)

func ConnectDB() (*gorm.DB, error) {
	connectionString := fmt.Sprintf("root:root@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		DB_Address, DB_Name)

	return gorm.Open(mysql.Open(connectionString), &gorm.Config{})
}

func MigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(
		model.User{},
	)
}
