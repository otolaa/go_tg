package config

import (
	"bufio"
	"log"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	VERSION string = "0.0.1"
)

var TOKEN string
var URL_BOT string
var DB_DSN string
var DB *gorm.DB

type User struct {
	ID      uint   `gorm:"primarykey"`
	Tid     int64  `gorm:"unique_index"`
	Name    string `gorm:"size:255"`
	Active  bool   `gorm:"type:bool"`
	Sending bool   `gorm:"type:bool"`
}

// get data from .env
func init() {
	file, err := os.Open(".env")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		if strings.Contains(s, "TOKEN=") {
			TOKEN = strings.ReplaceAll(s, "TOKEN=", "")
		}

		if strings.Contains(s, "URL_BOT=") {
			URL_BOT = strings.ReplaceAll(s, "URL_BOT=", "")
		}

		if strings.Contains(s, "DB_DSN=") {
			DB_DSN = strings.ReplaceAll(s, "DB_DSN=", "")
		}
	}

	db, err := gorm.Open(postgres.Open(DB_DSN), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	db.AutoMigrate(&User{})
	DB = db
}

func SetUser(db *gorm.DB, tid int64, userName string) User {
	var us User
	us.Tid = tid
	us.Name = userName
	us.Sending = true
	us.Active = true
	db.FirstOrCreate(&us, User{Tid: tid})

	return us
}

func GetListUsers(db *gorm.DB) []User {
	var users []User
	db.Select("id", "tid", "name").Where("active = ?", true).Where("sending = ?", true).Find(&users)

	return users
}

func SetUserSending(db *gorm.DB, uid uint, sending bool) User {
	var us User
	us.ID = uid
	db.Model(&us).Update("sending", sending)

	return us
}
