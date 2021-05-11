package scripts

import (
	"encoding/json"
	"log"
	"os"

	"github.com/marianodsr/nura-api/companies"
	"github.com/marianodsr/nura-api/storage"
	"github.com/marianodsr/nura-api/users"
)

func FillDb() {
	db := storage.GetDbConnection()
	db.AutoMigrate(&users.User{}, &companies.Company{})

	type jsonData struct {
		Users     []users.User        `json:"users"`
		Companies []companies.Company `json:"companies"`
	}

	file, err := os.Open("scripts/dbMock.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	data := &jsonData{}
	json.NewDecoder(file).Decode(data)

	for _, user := range data.Users {
		users.CreateUser(&user)
	}
	for _, company := range data.Companies {
		db.Create(&company)
	}
}
