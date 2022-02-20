package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/saadbadreddine/fsw-facebook-go-backend/database"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	Status      string `json:"status"`
	AccessToken string `json:"token"`
}

type AuthorizationToken struct {
	Token string `json:"token"`
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

var mySecretKey = []byte(os.Getenv("MY_JWT_TOKEN"))

//var mySigningKey = []byte("charizard010")

func GenerateJWT(id int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user_id"] = id
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySecretKey)

	if err != nil {
		fmt.Printf("Something went wrong: %s", err.Error())
		return "", err
	}

	return tokenString, err
}

func SignIn(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	//var str = string(body)
	var creds Credentials
	json.Unmarshal(body, &creds)

	password_bytes := []byte(creds.Password)
	hash := sha256.New()
	hash.Write(password_bytes)
	hash.Sum(nil)
	//hashed_password := sha256.Sum256(password_bytes)
	//fmt.Println(hashed_password)
	str_hashed_pass := hex.EncodeToString(hash.Sum(nil))
	var user User

	database.Connector.Table("users").Where("email = ? AND password = ?", creds.Email, str_hashed_pass).Select("id").Scan(&user)

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

func GetUserData(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	//var str = string(body)
	var auth_token AuthorizationToken
	json.Unmarshal(body, &auth_token)

	//fmt.Println(reflect.TypeOf(auth_token.Token))
	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(auth_token.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(mySecretKey), nil
	})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(claims)

	var user_id = claims["user_id"]
	stmt, err := database.Connector.DB().Prepare("SELECT first_name, last_name FROM users JOIN addresses ON  users.address_id = addresses.address_id WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	//var user User
	fmt.Println(user_id)

	result, err := stmt.Query(user_id)
	if err != nil {
		fmt.Println(err)
	}
	var users []User
	//result.Scan(&user.First_Name, &user.Last_Name)
	for result.Next() {
		var user User
		err := result.Scan(&user.First_Name, &user.Last_Name)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	json_response, err := json.Marshal(users[0])

	if err != nil {
		fmt.Println(err)
	}

	w.Write([]byte(json_response))

}
