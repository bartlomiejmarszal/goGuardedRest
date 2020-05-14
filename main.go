package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type JWT struct {
	Token string `json:"token"`
}

type Error struct {
	Message string `json:"message"`
}

var db *sql.DB

func main() {
	pgURL, err := pq.ParseURL("postgres://bart:bart@localhost/home?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", pgURL)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/signup", signup).Methods("POST")
	r.HandleFunc("/login", login).Methods("POST")

	// r.HandleFunc("/protected", TokenVerifyMiddleware(protectedEndpoint)).Methods("GET")

	log.Println("Web server has been started...")
	log.Fatal(http.ListenAndServe(":8800", r))

}

func respondWithError(w http.ResponseWriter, status int, error Error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)
}

func responseJSON(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}

func signup(w http.ResponseWriter, req *http.Request) {
	var user User
	var error Error
	json.NewDecoder(req.Body).Decode(&user)

	if user.Password == "" {
		error.Message = "Password is missing"
		respondWithError(w, http.StatusBadRequest, error)
		return
	}

	if user.Email == "" {
		error.Message = "Email is missing"
		respondWithError(w, http.StatusBadRequest, error)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil {
		log.Fatal(err)
		return
	}

	spew.Dump(&user.ID, user.Email, user.Password)

	user.Password = string(hash)

	var userid int

	// err = db.QueryRow(`INSERT INTO users(email, password)
	// VALUES('$1', '$2') RETURNING id`, user.Email, user.Password).Scan(&userid)  // doesn't work either

	err = db.QueryRow(`insert into users (email, password) values($1, $2) RETURNING id;`, user.Email, user.Password).Scan(&result)

	if err != nil {
		error.Message = "SQL Server error"
		respondWithError(w, http.StatusInternalServerError, error)
	}

	user.Password = ""
	json.NewEncoder(w).Encode(userid)
}

func login(w http.ResponseWriter, req *http.Request) { fmt.Println("login invoked") }

func protectedEndpoint(w http.ResponseWriter, req *http.Request) {
	// fmt.Println("ProtectedEndpoint invoked")
}

func TokenVerifyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	// fmt.Println("TokenVerifyMiddleware invoked")
	return nil
}
