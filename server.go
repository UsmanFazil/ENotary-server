package main

import (
	"ENOTARY-Server/DB"
	"ENOTARY-Server/Hashing"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"upper.io/db.v3/mysql"
)

func main() {

	// Database connection here
	var settings = mysql.ConnectionURL{
		User:     "root",
		Host:     "localhost",
		Password: "mypass",
		Database: "ENotary",
	}
	db, err := DB.Dbinit(settings)
	if err != nil {
		log.Fatal("error : ", err.Error())
	}
	log.Print("Maria DB server started ....")
	defer db.CloseSession()

	var dir = "./Files"

	//routes for the server
	r := mux.NewRouter()
	r.HandleFunc("/login", db.Login).Methods(http.MethodPost)
	r.HandleFunc("/signup", db.Signup).Methods(http.MethodPost)
	r.HandleFunc("/hashfile", Hashing.Servehash).Methods(http.MethodPost)
	r.HandleFunc("/verifyemail", db.AccountVerif).Methods(http.MethodPost)
	r.HandleFunc("/resendcode", db.ResendCode).Methods(http.MethodPost)
	r.HandleFunc("/inbox", db.InboxData).Methods(http.MethodGet)
	r.HandleFunc("/sent", db.SentContract).Methods(http.MethodGet)
	r.HandleFunc("/uploadprofilepic", db.ProfilePic).Methods(http.MethodPost)
	r.PathPrefix("/Files/").Handler(http.StripPrefix("/Files/", http.FileServer(http.Dir(dir))))

	log.Println("Go-lang server started at port 8000 ....")
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r)))
}
