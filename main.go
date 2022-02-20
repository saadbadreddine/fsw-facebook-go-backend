package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//Config to maintain DB configuration properties
type Config struct {
	ServerName string
	User       string
	Password   string
	DB         string
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	Status      string `json:"status"`
	AccessToken string `json:"token"`
}

type User struct {
	ID         int
	First_Name string
	Last_Name  string
	Dob        string
	Email      string
	Password   string
	Timestamp  string
	Address_ID int
}

func main() {

	//Connect creates MySQL connection
	config :=
		Config{
			ServerName: "localhost:3306",
			User:       "debian-sys-maint",
			Password:   "7LRTlMIJFQQH3tSc",
			DB:         "facebookdb",
		}

	connectionString := getConnectionString(config)
	err := Connect(connectionString)
	if err != nil {
		panic(err.Error())
	}

	// The router is now formed by calling the `newRouter` constructor function
	// that we defined above. The rest of the code stays the same
	r := newRouter()
	http.ListenAndServe(":8080", r)
}

// The new router function creates the router and
// returns it to us. We can now use this function
// to instantiate and test the router outside of the main function
func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/signin", SignIn).Methods("POST")

	// Declare the static file directory and point it to the
	// directory we just made
	staticFileDirectory := http.Dir("./assets/")
	// Declare the handler, that routes requests to their respective filename.
	// The fileserver is wrapped in the `stripPrefix` method, because we want to
	// remove the "/assets/" prefix when looking for files.
	// For example, if we type "/assets/index.html" in our browser, the file server
	// will look for only "index.html" inside the directory declared above.
	// If we did not strip the prefix, the file server would look for
	// "./assets/assets/index.html", and yield an error
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
	// The "PathPrefix" method acts as a matcher, and matches all routes starting
	// with "/assets/", instead of the absolute route itself
	r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")
	return r
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

var getConnectionString = func(config Config) string {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&multiStatements=true", config.User, config.Password, config.ServerName, config.DB)
	return connectionString
}

//var mySigningKey = os.Get("MY_JWT_TOKEN")
var mySigningKey = []byte("charizard010")

func GenerateJWT(id int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user_id"] = id
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something went wrong: %s", err.Error())
		return "", err
	}

	return tokenString, err
}

func SignIn(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var str = string(body)
	var creds Credentials
	json.Unmarshal([]byte(str), &creds)

	fmt.Println(creds.Password)
	password_bytes := []byte(creds.Password)
	hash := sha256.New()
	hash.Write(password_bytes)
	hash.Sum(nil)
	//hashed_password := sha256.Sum256(password_bytes)
	//fmt.Println(hashed_password)
	str_hashed_pass := hex.EncodeToString(hash.Sum(nil))
	fmt.Println(str_hashed_pass)
	json.Unmarshal([]byte(str), &creds)
	var user User

	Connector.Table("users").Where("email = ? AND password = ?", creds.Email, str_hashed_pass).Select("id").Scan(&user)

	var loginUserResponse LoginUserResponse

	if user.ID != 0 {

		loginUserResponse.AccessToken, err = GenerateJWT(user.ID)
		if err != nil {
			fmt.Println(err)
			return
		}
		loginUserResponse.Status = "Logged in"

	} else {
		loginUserResponse.Status = "Incorrect combination"
	}

	json_response, err := json.Marshal(loginUserResponse)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write([]byte(json_response))
}
