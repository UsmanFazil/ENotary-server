package main

import (
	"ENOTARY-Server/DB"
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
	r.Handle("/uploadProfilePic", db.IsAuthorized(db.ProfilePic)).Methods(http.MethodPost)
	r.Handle("/removeProfilePic", db.IsAuthorized(db.RemovePic)).Methods(http.MethodGet)
	r.Handle("/Logout", db.IsAuthorized(db.Logout)).Methods(http.MethodGet)
	r.Handle("/uploadSign", db.IsAuthorized(db.UploadSign)).Methods(http.MethodPost)
	r.Handle("/signbase64", db.IsAuthorized(db.SignBase64)).Methods(http.MethodPost)
	r.Handle("/userprefs", db.IsAuthorized(db.Userpreferences)).Methods(http.MethodPost)

	r.Handle("/manage", db.IsAuthorized(db.Manage)).Methods(http.MethodGet)
	r.Handle("/inbox", db.IsAuthorized(db.InboxData)).Methods(http.MethodGet)
	r.Handle("/sent", db.IsAuthorized(db.SentContract)).Methods(http.MethodGet)
	r.Handle("/drafts", db.IsAuthorized(db.DraftContracts)).Methods(http.MethodGet)
	r.Handle("/actionReq", db.IsAuthorized(db.ActionRequired)).Methods(http.MethodGet)
	r.Handle("/expSoon", db.IsAuthorized(db.ExpiringsoonContracts)).Methods(http.MethodGet)
	r.Handle("/waitingForOther", db.IsAuthorized(db.WaitingForOthers)).Methods(http.MethodGet)
	r.Handle("/completed", db.IsAuthorized(db.Completed)).Methods(http.MethodGet)
	r.Handle("/searchContract", db.IsAuthorized(db.SearchAlgo)).Methods(http.MethodPost)

	r.Handle("/newContract", db.IsAuthorized(db.NewContract)).Methods(http.MethodPost)
	r.Handle("/addRecipients", db.IsAuthorized(db.AddRecipients)).Methods(http.MethodPost)
	r.Handle("/delDraft", db.IsAuthorized(db.DeleteDraft)).Methods(http.MethodDelete)
	r.Handle("/ContractDetails", db.IsAuthorized(db.ContractDetails)).Methods(http.MethodPost)
	r.Handle("/SendContract", db.IsAuthorized(db.SendContract)).Methods(http.MethodPost)
	r.Handle("/SaveinBlockchain", db.IsAuthorized(db.ContractHashDetails)).Methods(http.MethodPost)
	r.Handle("/updateBlockchainstatus", db.IsAuthorized(db.UpdateBlockchainstatus)).Methods(http.MethodPost)
	r.Handle("/verifyContract", db.IsAuthorized(db.VerifyContract)).Methods(http.MethodPost)

	r.Handle("/newFolder", db.IsAuthorized(db.NewFolder)).Methods(http.MethodPost)
	r.Handle("/moveContract", db.IsAuthorized(db.AddContract)).Methods(http.MethodPost)
	r.Handle("/folderContractList", db.IsAuthorized(db.FolderContractList)).Methods(http.MethodPost)

	r.PathPrefix("/Files/").Handler(http.StripPrefix("/Files/", http.FileServer(http.Dir(dir))))

	log.Println("Go-lang server started at port 8000 ....")
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Token"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r)))
}
