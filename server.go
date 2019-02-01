package main

import (
	"ENOTARY-Server/DB"
	"ENOTARY-Server/Hashing"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"upper.io/db.v3/mysql"
)

func main() {

	// Database connection here
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

	//routes for the server
	r := mux.NewRouter()
	r.HandleFunc("/login", db.Login).Methods(http.MethodPost)
	r.HandleFunc("/signup", db.Signup).Methods(http.MethodPost)
	r.HandleFunc("/hashfile", Hashing.Servehash).Methods(http.MethodPost)
	r.HandleFunc("/verifyemail", db.AccountVerif).Methods(http.MethodPost)
	r.HandleFunc("/resendcode", db.ResendCode).Methods(http.MethodPost)
	r.HandleFunc("/inbox", db.InboxData).Methods(http.MethodGet)
	r.HandleFunc("/sent", db.SentContract).Methods(http.MethodGet)

	log.Println("Go-lang server started at port 8000 ....")
	log.Println(http.ListenAndServe(":8000", r))

}
