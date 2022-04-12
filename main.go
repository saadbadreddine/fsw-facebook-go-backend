package main

import (
	"net/http"

	"github.com/saadbadreddine/fsw-facebook-go-backend/database"
	"github.com/saadbadreddine/fsw-facebook-go-backend/router"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {

	config :=
		database.Config{
			ServerName: "localhost:3306",
			User:       "debian-sys-maint",
			Password:   "7LRTlMIJFQQH3tSc",
			DB:         "facebookdb",
		}

	//Connect creates MySQL connection
	connectionString := database.GetConnectionString(config)
	err := database.Connect(connectionString)
	if err != nil {
		panic(err.Error())
	}

	r := router.NewRouter()
	http.ListenAndServe(":8080", r)
}
