package database

import (
	"log"

	"github.com/jinzhu/gorm"
)

type Config struct {
	ServerName string
	User       string
	Password   string
	DB         string
}

//Connector variable used for CRUD operation's
var Connector *gorm.DB

//Connect creates MySQL connection
func Connect(connectionString string) error {
	var err error
	Connector, err = gorm.Open("mysql", connectionString)
	if err != nil {
		return err
	}
	log.Println("Connection was successful!!")
	return nil
}
