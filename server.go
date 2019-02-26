package main

import (
	"ENOTARY-Server/DB"
	"ENOTARY-Server/Hashing"
	MW "ENOTARY-Server/Middleware"
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

	//API ENDPOINTS HERE
	r := mux.NewRouter()
	r.HandleFunc("/login", db.Login).Methods(http.MethodPost)
	r.HandleFunc("/signup", db.Signup).Methods(http.MethodPost)
	r.HandleFunc("/verifyEmail", db.EmailVerification).Methods(http.MethodPost)
	r.HandleFunc("/sendCode", db.SendCode).Methods(http.MethodGet)
	r.Handle("/inbox", MW.IsAuthorized(db.InboxData)).Methods(http.MethodGet)
	r.Handle("/sent", MW.IsAuthorized(db.SentContract)).Methods(http.MethodGet)
	r.Handle("/uploadProfilePic", MW.IsAuthorized(db.ProfilePic)).Methods(http.MethodPost)
	r.Handle("/newContract", MW.IsAuthorized(db.NewContract)).Methods(http.MethodPost)
	r.Handle("/addRecipients", MW.IsAuthorized(db.AddRecipients)).Methods(http.MethodPost)
	r.Handle("/hashFile", MW.IsAuthorized(Hashing.Servehash)).Methods(http.MethodPost)
	r.Handle("/logout", MW.IsAuthorized(db.Logout)).Methods(http.MethodGet)

	r.PathPrefix("/Files/").Handler(http.StripPrefix("/Files/", http.FileServer(http.Dir(dir))))

	log.Println("Go-lang server started at port 8000 ....")
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r)))
}
