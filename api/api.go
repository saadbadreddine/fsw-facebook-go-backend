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
	ID         int    `json:"id"`
	First_Name string `json:"first_name"`
	Last_Name  string `json:"last_name"`
	Dob        string `json:"dob"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Timestamp  string `json:"timestamp"`
	Address_ID int    `json:"address_id"`
}

type Post struct {
	First_Name string `json:"first_name"`
	Last_Name  string `json:"last_name"`
	Post       string `json:"post"`
	Timestamp  string `json:"timestamp"`
	User_ID    string `json:"user_id"`
}

type Action struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
}

type ActionResponse struct {
	Status string `json:"status"`
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

	var creds Credentials
	json.Unmarshal(body, &creds)
	//fmt.Println(reflect.TypeOf(creds.Password))
	password_bytes := []byte(creds.Password)
	hash := sha256.New()
	hash.Write(password_bytes)
	hash.Sum(nil)

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

	var auth_token AuthorizationToken
	json.Unmarshal(body, &auth_token)

	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(auth_token.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(mySecretKey), nil
	})
	if err != nil {
		fmt.Println(err)
	}

	var user_id = claims["user_id"]
	stmt, err := database.Connector.DB().Prepare(`SELECT first_name, last_name FROM users 
	JOIN addresses ON  users.address_id = addresses.address_id WHERE id = ?`)
	if err != nil {
		log.Fatal(err)
	}

	result, err := stmt.Query(user_id)
	if err != nil {
		fmt.Println(err)
	}
	var users []User

	for result.Next() {
		var user User
		err := result.Scan(&user.First_Name, &user.Last_Name)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	stmt.Close()

	json_response, err := json.Marshal(users[0])

	if err != nil {
		fmt.Println(err)
	}

	w.Write([]byte(json_response))

}

func GetPosts(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var auth_token AuthorizationToken
	json.Unmarshal(body, &auth_token)

	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(auth_token.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(mySecretKey), nil
	})
	if err != nil {
		fmt.Println(err)
	}

	var user_id = claims["user_id"]
	stmt, err := database.Connector.DB().Prepare(`SELECT DISTINCT posts.post, posts.timestamp, posts.user_id, users.first_name, users.last_name 
	FROM posts JOIN users ON posts.user_id = users.id JOIN friendships ON (posts.user_id = friendships.sender OR posts.user_id = friendships.receiver) 
	WHERE(friendships.sender = ? OR friendships.receiver = ?) AND friendships.accepted = 1 AND users.id NOT IN (SELECT blocks.receiver FROM blocks WHERE blocks.receiver = ? OR blocks.sender = ?) 
	ORDER BY timestamp DESC;`)

	if err != nil {
		log.Fatal(err)
	}

	result, err := stmt.Query(user_id, user_id, user_id, user_id)
	if err != nil {
		fmt.Println(err)
	}
	var posts []Post
	//result.Scan(&user.First_Name, &user.Last_Name)
	for result.Next() {
		var post Post
		err := result.Scan(&post.Post, &post.Timestamp, &post.User_ID, &post.First_Name, &post.Last_Name)
		if err != nil {
			log.Fatal(err)
		}
		posts = append(posts, post)
	}

	stmt.Close()

	json_response, err := json.Marshal(posts)

	if err != nil {
		fmt.Println(err)
	}

	w.Write([]byte(json_response))

}

func GetFriends(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var auth_token AuthorizationToken
	json.Unmarshal(body, &auth_token)

	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(auth_token.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(mySecretKey), nil
	})
	if err != nil {
		fmt.Println(err)
	}

	var user_id = claims["user_id"]

	stmt, err := database.Connector.DB().Prepare(`SELECT Distinct id, first_name, last_name 
	FROM users INNER JOIN friendships ON users.id = friendships.sender OR users.id = friendships.receiver 
	LEFT JOIN blocks ON  users.id = blocks.receiver OR users.id = blocks.sender 
	WHERE (friendships.sender = ? OR friendships.receiver = ?) AND friendships.accepted = 1 AND id != ? 
	AND id NOT IN (SELECT blocks.sender FROM blocks WHERE blocks.receiver = ?) 
	AND id NOT IN (SELECT blocks.receiver FROM blocks WHERE blocks.sender = ?)`)

	if err != nil {
		log.Fatal(err)
	}

	result, err := stmt.Query(user_id, user_id, user_id, user_id, user_id)
	if err != nil {
		fmt.Println(err)
	}
	var users []User
	for result.Next() {
		var user User
		err := result.Scan(&user.ID, &user.First_Name, &user.Last_Name)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	stmt.Close()

	json_response, err := json.Marshal(users)

	if err != nil {
		fmt.Println(err)
	}

	w.Write([]byte(json_response))

}

func GetFriendRequests(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var auth_token AuthorizationToken
	json.Unmarshal(body, &auth_token)

	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(auth_token.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(mySecretKey), nil
	})
	if err != nil {
		fmt.Println(err)
	}

	var user_id = claims["user_id"]

	stmt, err := database.Connector.DB().Prepare(`SELECT id, first_name, last_name 
	FROM users INNER JOIN friendships ON users.id = friendships.sender OR users.id = friendships.receiver 
	WHERE friendships.receiver = ? AND friendships.accepted = 0 AND id != ?`)

	if err != nil {
		log.Fatal(err)
	}

	result, err := stmt.Query(user_id, user_id)
	if err != nil {
		fmt.Println(err)
	}
	var users []User
	for result.Next() {
		var user User
		err := result.Scan(&user.ID, &user.First_Name, &user.Last_Name)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	stmt.Close()

	json_response, err := json.Marshal(users)

	if err != nil {
		fmt.Println(err)
	}

	w.Write([]byte(json_response))

}

func GetBlockedUsers(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var auth_token AuthorizationToken
	json.Unmarshal(body, &auth_token)

	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(auth_token.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(mySecretKey), nil
	})
	if err != nil {
		fmt.Println(err)
	}

	var user_id = claims["user_id"]

	stmt, err := database.Connector.DB().Prepare(`SELECT id, first_name, last_name 
	FROM users INNER JOIN blocks ON users.id = blocks.receiver WHERE blocks.sender = ?`)

	if err != nil {
		log.Fatal(err)
	}

	result, err := stmt.Query(user_id)
	if err != nil {
		fmt.Println(err)
	}
	var users []User
	for result.Next() {
		var user User
		err := result.Scan(&user.ID, &user.First_Name, &user.Last_Name)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	stmt.Close()

	json_response, err := json.Marshal(users)

	if err != nil {
		fmt.Println(err)
	}

	w.Write([]byte(json_response))

}

func AcceptFriendRequest(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var action Action
	json.Unmarshal(body, &action)

	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(action.Sender, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(mySecretKey), nil
	})
	if err != nil {
		fmt.Println(err)
	}

	var user_id = claims["user_id"]

	accepted := true
	var action_response ActionResponse

	stmt, err := database.Connector.DB().Prepare(`UPDATE friendships SET accepted = ? WHERE sender = ? AND receiver = ?`)

	if err != nil {
		log.Fatal(err)
	}

	stmt.Query(accepted, action.Receiver, user_id)
	if err != nil {
		action_response.Status = "Failed"
		fmt.Println(err)
	}

	stmt.Close()
	action_response.Status = "Success"
	json_response, err := json.Marshal(action_response)

	if err != nil {
		fmt.Println(err)
	}

	w.Write([]byte(json_response))

}
