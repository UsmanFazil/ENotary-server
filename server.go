package main

import (
	"ENOTARY-Server/DB"
	"ENotary-server/Hashing"
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

	//API ENDPOINTS
	r := mux.NewRouter()
	r.HandleFunc("/login", db.Login).Methods(http.MethodPost)
	r.HandleFunc("/signup", db.Signup).Methods(http.MethodPost)
	r.HandleFunc("/verifyEmail", db.EmailVerification).Methods(http.MethodPost)
	r.HandleFunc("/sendCode", db.SendCode).Methods(http.MethodPost)
	r.HandleFunc("/updatePass", db.UpdatePassword).Methods(http.MethodPost)
	r.Handle("/inbox", db.IsAuthorized(db.InboxData)).Methods(http.MethodGet)
	r.Handle("/sent", db.IsAuthorized(db.SentContract)).Methods(http.MethodGet)
	r.Handle("/uploadProfilePic", db.IsAuthorized(db.ProfilePic)).Methods(http.MethodPost)
	r.Handle("/newContract", db.IsAuthorized(db.NewContract)).Methods(http.MethodPost)
	r.Handle("/addRecipients", db.IsAuthorized(db.AddRecipients)).Methods(http.MethodPost)
	r.Handle("/hashFile", db.IsAuthorized(Hashing.Servehash)).Methods(http.MethodPost)
	r.Handle("/newFolder", db.IsAuthorized(db.NewFolder)).Methods(http.MethodPost)
	r.Handle("/moveContract", db.IsAuthorized(db.AddContract)).Methods(http.MethodPost)
	r.Handle("/folderContractList", db.IsAuthorized(db.FolderContractList)).Methods(http.MethodPost)
	r.Handle("/searchContract", db.IsAuthorized(db.GenericSearch)).Methods(http.MethodPost)
	r.Handle("/Logout", db.IsAuthorized(db.Logout)).Methods(http.MethodGet)

	r.PathPrefix("/Files/").Handler(http.StripPrefix("/Files/", http.FileServer(http.Dir(dir))))

	log.Println("Go-lang server started at port 8000 ....")
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r)))
}
