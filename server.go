package main

import (
	"ENOTARY-Server/DB"
	"ENOTARY-Server/Hashing"
	"fmt"
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
	r.HandleFunc("/hashFile", Hashing.Servehash).Methods(http.MethodGet)
	r.HandleFunc("/verifyEmail", db.AccountVerif).Methods(http.MethodPost)
	r.HandleFunc("/resendCode", db.ResendCode).Methods(http.MethodGet)
	r.HandleFunc("/inbox", db.InboxData).Methods(http.MethodGet)
	r.HandleFunc("/sent", db.SentContract).Methods(http.MethodGet)
	r.HandleFunc("/uploadProfilePic", db.ProfilePic).Methods(http.MethodPost)
	r.HandleFunc("/newContract", db.NewContract).Methods(http.MethodPost)
	r.HandleFunc("/addRecipients", db.AddRecipients).Methods(http.MethodPost)

	r.PathPrefix("/Files/").Handler(http.StripPrefix("/Files/", http.FileServer(http.Dir(dir))))

	log.Println("Go-lang server started at port 8000 ....")
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r)))
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] != nil {
			fmt.Println("auth")
		} else {

			fmt.Fprintf(w, "Not Authorized")
		}
	})
}
