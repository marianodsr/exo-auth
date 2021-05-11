package companies

import (
	"fmt"

	"github.com/marianodsr/nura-api/storage"
	"gorm.io/gorm"
)

//Company struct
type Company struct {
	gorm.Model
	Name   string `json:"name"`
	Logo   string `json:"logo"`
	UserID uint
}

func GetCompanyByID(id uint) *Company {
	company := &Company{}
	db := storage.GetDbConnection()
	db.First(company, id)
	return company
}

func UpdateCompany(company *Company) error {
	db := storage.GetDbConnection()
	res := db.Save(company)
	if res.Error != nil {
		return fmt.Errorf("Error updating company")
	}
	return nil
}

func CreateCompany(company *Company) error {
	db := storage.GetDbConnection()
	res := db.Create(company)
	if res.Error != nil {
		fmt.Printf("Db error creating user: %s", res.Error.Error())
		return fmt.Errorf("Error creating user...")
	}
	return nil
}
