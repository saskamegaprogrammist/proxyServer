package main

import (
	"github.com/gorilla/mux"
	"github.com/saskamegaprogrammist/proxyServer/db"
	"github.com/saskamegaprogrammist/proxyServer/handlers"
	"log"
	"net/http"
)

func main() {
	db.CreateDataBaseConnection("alexis", "sinope27", "localhost", "proxy", 20)
	//db.InitDataBase()
	r := mux.NewRouter()
	r.HandleFunc("/lastreqs", handlers.GetLastRequests).Methods("GET")
	r.HandleFunc("/request/{id}", handlers.MakeRequest).Methods("GET")

	err := http.ListenAndServe(":5001", r)
	if err != nil {
		log.Fatal(err)
		return
	}

}