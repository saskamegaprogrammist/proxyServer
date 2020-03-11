package main

import (
	"github.com/gorilla/mux"
	"github.com/saskamegaprogrammist/proxyServer/db"
	"github.com/saskamegaprogrammist/proxyServer/handlers"
	"log"
	"net/http"
)

func main() {
	db.CreateDataBaseConnection("docker", "docker", "localhost", "docker", 20)
	//db.CreateDataBaseConnection("alexis", "sinope27", "localhost", "proxy", 20)
	r := mux.NewRouter()
	handlers.RepeaterClient = &http.Client{}
	r.HandleFunc("/requests", handlers.GetLastRequests).Methods("GET")
	r.HandleFunc("/requests/{id}", handlers.MakeRequest).Methods("GET")
	


	err := http.ListenAndServe(":5001", r)
	if err != nil {
		log.Fatal(err)
		return
	}

}