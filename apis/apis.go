package apis

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

var mySigningKey = []byte(os.Getenv("MY_JWT_TOKEN"))

//var mySigningKey = []byte("charizard010")

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

	password_bytes := []byte(creds.Password)
	hash := sha256.New()
	hash.Write(password_bytes)
	hash.Sum(nil)
	//hashed_password := sha256.Sum256(password_bytes)
	//fmt.Println(hashed_password)
	str_hashed_pass := hex.EncodeToString(hash.Sum(nil))
	json.Unmarshal([]byte(str), &creds)
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

	var str = string(body)
	var auth_token AuthorizationToken
	json.Unmarshal([]byte(str), &auth_token)

	json_response, err := json.Marshal(auth_token)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write([]byte(json_response))

}
