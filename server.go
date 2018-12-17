package main

import (
	"Server/DB"
	"Server/Hashing"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"upper.io/db.v3/mysql"
)

func main() {

	var settings = mysql.ConnectionURL{
		User:     "root",
		Host:     "localhost",
		Password: "",
		Database: "ENotary",
	}
	db, err := DB.Dbinit(settings)
	if err != nil {
		log.Fatal("error : ", err.Error())
	}
	log.Print("Maria DB server started ....")
	defer db.CloseSession()

	r := mux.NewRouter()

	//routes for the server
	r.HandleFunc("/login", db.Login).Methods("POST")
	r.HandleFunc("/Signup", db.Newuser).Methods("POST")
	r.HandleFunc("/hashfile", Hashing.Servehash).Methods("POST")
	r.HandleFunc("/validateuser/{email}", db.Validateuser).Methods("GET")

	log.Println("Go-lang server started at port 8000 ...")
	log.Println(http.ListenAndServe(":8000", r))

}
