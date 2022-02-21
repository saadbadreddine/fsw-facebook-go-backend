package main

import (
	"net/http"

	"github.com/saadbadreddine/fsw-facebook-go-backend/api"
	"github.com/saadbadreddine/fsw-facebook-go-backend/database"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {

	//Connect creates MySQL connection
	config :=
		database.Config{
			ServerName: "localhost:3306",
			User:       "debian-sys-maint",
			Password:   "7LRTlMIJFQQH3tSc",
			DB:         "facebookdb",
		}

	connectionString := database.GetConnectionString(config)
	err := database.Connect(connectionString)
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
	r.HandleFunc("/signin", api.SignIn).Methods("POST")
	r.HandleFunc("/getdata", api.GetUserData).Methods("POST")
	r.HandleFunc("/getposts", api.GetPosts).Methods("POST")
	r.HandleFunc("/getfriends", api.GetFriends).Methods("POST")
	r.HandleFunc("/getfriendrequests", api.GetFriendRequests).Methods("POST")
	r.HandleFunc("/getblockedusers", api.GetBlockedUsers).Methods("POST")
	r.HandleFunc("/acceptfriendrequest", api.AcceptFriendRequest).Methods("POST")
	r.HandleFunc("/rejectfriendrequest", api.RejectFriendRequest).Methods("POST")
	r.HandleFunc("/blockfriend", api.BlockFriend).Methods("POST")
	r.HandleFunc("/unblockfriend", api.UnblockFriend).Methods("POST")
	r.HandleFunc("/removefriend", api.RemoveFriend).Methods("POST")
	r.HandleFunc("/addfriend", api.AddFriend).Methods("POST")

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
