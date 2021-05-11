package storage

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

//InitDB func
func InitDB() {
	dsn := "host=localhost user=postgres password=auth159123456 dbname=exo-auth port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

}

//GetDbConnection func
func GetDbConnection() *gorm.DB {
	return db
}
