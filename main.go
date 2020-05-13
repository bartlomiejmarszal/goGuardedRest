package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/signup", signup).Methods("POST")
	r.HandleFunc("/login", login).Methods("POST")

	r.HandleFunc("/protected", TokenVerifyMiddleware(protectedEndpoint)).Methods("GET")

	log.Println("Web server has been started...")
	log.Fatal(http.ListenAndServe(":8000", r))

}

func signup(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Success"))
}

func login(w http.ResponseWriter, req *http.Request) { fmt.Println("login invoked") }

func protectedEndpoint(w http.ResponseWriter, req *http.Request) {
	fmt.Println("ProtectedEndpoint invoked")
}

func TokenVerifyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	fmt.Println("TokenVerifyMiddleware invoked")
	return nil
}
