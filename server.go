package main

import (
	"ENOTARY-Server/DB"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

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
	r.Handle("/playgroundinput", db.IsAuthorized(db.SaveCoordinates)).Methods(http.MethodPost)
	r.Handle("/serveCoordinates", db.IsAuthorized(db.ServeCoordinates)).Methods(http.MethodPost)
	r.Handle("/signContract", db.IsAuthorized(db.SignContract)).Methods(http.MethodPost)
	r.Handle("/DeclineContract", db.IsAuthorized(db.DeclineContract)).Methods(http.MethodPost)
	r.Handle("/ExportCSV", db.IsAuthorized(db.ExportCSV)).Methods(http.MethodPost)

	r.Handle("/newFolder", db.IsAuthorized(db.NewFolder)).Methods(http.MethodPost)
	r.Handle("/moveContract", db.IsAuthorized(db.AddContract)).Methods(http.MethodPost)
	r.Handle("/folderContractList", db.IsAuthorized(db.FolderContractList)).Methods(http.MethodPost)

	r.HandleFunc("/test", Test)

	r.PathPrefix("/Files/").Handler(http.StripPrefix("/Files/", http.FileServer(http.Dir(dir))))

	log.Println("Go-lang server started at port 8000 ....")
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Token"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r)))
}

func Test(writer http.ResponseWriter, request *http.Request) {
	Filename := request.URL.Query().Get("file")
	if Filename == "" {
		//Get not set, send a 400 bad request
		//http.Error(writer, "Get 'file' not specified in url.", 400)
		return
	}
	fmt.Println("Client requests: " + Filename)

	//Check if file exists and open
	Openfile, err := os.Open(Filename)
	defer Openfile.Close() //Close after function return
	if err != nil {
		//File not found, send 404
		//http.Error(writer, "File not found.", 404)
		return
	}

	//File is found, create and send the correct headers

	//Get the Content-Type of the file
	//Create a buffer to store the header of the file in
	FileHeader := make([]byte, 512)
	//Copy the headers into the FileHeader buffer
	Openfile.Read(FileHeader)
	//Get content type of file
	FileContentType := http.DetectContentType(FileHeader)

	//Get the file size
	FileStat, _ := Openfile.Stat()                     //Get info from file
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

	//Send the headers
	writer.Header().Set("Content-Disposition", "attachment; filename="+Filename)
	writer.Header().Set("Content-Type", FileContentType)
	writer.Header().Set("Content-Length", FileSize)

	//Send the file
	//We read 512 bytes from the file already, so we reset the offset back to 0
	Openfile.Seek(0, 0)
	io.Copy(writer, Openfile) //'Copy' the file to the client
	return

}
