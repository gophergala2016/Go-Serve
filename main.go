package main

import (
	"github.com/gophergala2016/Go-Serve/api/v1/controllers/account"
	"github.com/gophergala2016/Go-Serve/api/v1/controllers/service"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	// Account Routes
	r.HandleFunc("/sign_up", account.Registration.Create).Methods("POST")
	r.HandleFunc("/service/create", service.Service.Create).Methods("POST")

	r.HandleFunc("/sign_in", account.Session.Create).Methods("POST")

	http.Handle("/", r)
	// HTTP Listening Port
	log.Println("main : Started : Listening on: http://localhost:3000 ...")
	log.Fatal(http.ListenAndServe("0.0.0.0:3000", nil))
}
