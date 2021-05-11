package users

import (
	"fmt"
	"time"

	"github.com/marianodsr/nura-api/authentication"
	"github.com/marianodsr/nura-api/companies"
	"github.com/marianodsr/nura-api/storage"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//User struct
type User struct {
	gorm.Model
	Name           string              `json:"name"`
	LastName       string              `json:"last_name"`
	Email          string              `json:"email" gorm:"default:0"`
	Password       string              `json:"password"`
	LoginCount     int                 `json:"login-count"`
	Companies      []companies.Company `json:"companies"`
	CurrentCompnay uint                `json:"current_company"`
	LastLogin      time.Time
}

func attemptLogin(email string, password string) ([]string, error) {
	user, err := GetUserByEmail(email)
	if err != nil || user == nil {
		return nil, fmt.Errorf("Invalid email or password")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("Invalid email or password")
	}
	pair, err := authentication.GenerateTokenPair(user.ID)
	if err != nil {
		return nil, fmt.Errorf("Invalid email or password")
	}
	user.LastLogin = time.Now()
	user.LoginCount = user.LoginCount + 1
	UpdateUser(user)
	return pair, err
}

func signUp(userAttempt *User) (*User, error) {
	user, err := CreateUser(userAttempt)
	user.CurrentCompnay = user.Companies[0].ID
	UpdateUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

//CreateUser func
func CreateUser(userAttempt *User) (*User, error) {
	db := storage.GetDbConnection()
	hashedPassword, err := hashPassword(userAttempt.Password)
	if err != nil {
		return nil, err
	}
	userAttempt.Password = string(hashedPassword)
	db.Create(userAttempt)
	return userAttempt, nil
}

func UpdateUser(user *User) error {
	db := storage.GetDbConnection()
	res := db.Save(user)
	if res.Error != nil {
		return fmt.Errorf("error updating user")
	}
	return nil
}

func hashPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func GetUserByID(userID uint) (*User, error) {
	db := storage.GetDbConnection()
	user := &User{}
	db.Where(&User{}).Preload("Companies").First(user, userID)
	if user.Email == "" {
		return nil, fmt.Errorf("error finding user")
	}
	fmt.Printf("\nUser companies: %+v\n", user.Companies)
	return user, nil
}

//GetUserByEmail func
func GetUserByEmail(email string) (*User, error) {
	db := storage.GetDbConnection()
	user := &User{}
	db.Where(&User{Email: email}).First(user)
	if user.Email == "" {
		return nil, db.Error
	}
	return user, nil
}
