package DB

import (
	"ENOTARY-Server/Email"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	jwt "github.com/dgrijalva/jwt-go"
	db "upper.io/db.v3"
)

func (d *dbServer) GetUser(email string) (*User, bool, error) {
	Collection := d.sess.Collection(UserCollection)
	res := Collection.Find(db.Cond{"email": email})
	var user User
	err := res.One(&user)

	if err != nil {
		return nil, false, err
	}
	return &user, true, nil
}

// RenderError : creates error response
func RenderError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

// RenderResponse : creates response
func RenderResponse(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func GenerateToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func CredentialValidation(usr User) (bool, string) {
	mailRE, errstring := VerifyEmail(usr.Email)
	if !mailRE {
		return false, errstring
	}
	pswdRE, errstring := VerifyPassword(usr.Password)
	if !pswdRE {
		return false, errstring
	}
	nameRE, errstring := verifyName(usr.Name)
	if !nameRE {
		return false, errstring
	}
	cmpRE, errstring := verifyComp(usr.Company)
	if !cmpRE {
		return false, errstring
	}
	phRE, errstring := verifyPhone(usr.Phone)
	if !phRE {
		return false, errstring
	}

	return true, ""
}

func verifyPhone(ph string) (bool, string) {
	var len int
	for _, _ = range ph {
		len++
	}
	_, err := strconv.Atoi(ph)
	if err != nil || len < 8 || len > 20 {
		return false, "Invalid number"
	}
	return true, ""
}

func verifyName(name string) (bool, string) {
	if len(name) < 5 {
		return false, "invalid name"
	}

	splitter := strings.Split(name, " ")

	if len(splitter) != 2 {
		return false, "Invalid name"
	}
	if splitter[1] == "" || splitter[0] == "" {
		return false, "invalid name"
	}

	for _, ch := range splitter[0] {
		if !unicode.IsLetter(ch) {
			return false, "invalid name"
		}
	}
	for _, ch := range splitter[1] {
		if !unicode.IsLetter(ch) {
			return false, "invalid name"
		}
	}

	return true, ""
}

func verifyComp(company string) (bool, string) {
	compRE := regexp.MustCompile("^[a-zA-Z0-9]{2,50}")

	if !compRE.MatchString(company) {
		return false, "invalid company name"
	}
	return true, ""
}

func VerifyEmail(mail string) (bool, string) {
	mailRE := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !mailRE.MatchString(mail) {
		return false, "invalid email address"
	}
	return true, ""
}

func VerifyPassword(password string) (bool, string) {
	//var uppercasePresent bool
	//var lowercasePresent bool
	//	var numberPresent bool
	var space bool
	const minPassLength = 8
	const maxPassLength = 24
	var passLen int

	for _, ch := range password {
		passLen++
		switch {
		// case unicode.IsNumber(ch):
		// 	numberPresent = true
		// case unicode.IsUpper(ch):
		// 	uppercasePresent = true
		// case unicode.IsLower(ch):
		// 	lowercasePresent = true
		case ch == ' ':
			space = true
		}
	}

	if space {
		return false, "space fields are not allowed in password"
	}
	if !(minPassLength <= passLen && passLen <= maxPassLength) {
		return false, "password length must be between 8 to 24 characters long"
	}

	// if !uppercasePresent {
	// 	return false, "uppercase letter missing in password"
	// }
	// if !numberPresent {
	// 	return false, "atleast one numeric character required in password"
	// }

	return true, ""
}

func GetClaims(tokenstring string) (jwt.MapClaims, bool) {

	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return MySigningKey, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		return nil, false
	}
}

func Logger(logstring string) {
	f, err := os.OpenFile("ENotary-logs.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logger := log.New(f, "ENotary-log ", log.LstdFlags)
	logger.Println(logstring)
}

func (d *dbServer) Getemails(cd ContractDetail) (bool, []string) {
	var emails = make([]string, len(cd.Signers)+1)
	var user User
	var idlists = make([]string, len(cd.Signers)+1)

	idlists[0] = cd.ContractData.Creator

	for i := 0; i < len(cd.Signers); i++ {
		idlists[i+1] = cd.Signers[i].UserID
	}

	collection := d.sess.Collection(UserCollection)

	for j := 0; j < len(idlists); j++ {
		res := collection.Find(db.Cond{"userid": idlists[j]})
		err := res.One(&user)

		if err != nil {
			return false, nil
		}
		emails[j] = user.Email
	}
	return true, emails
}

// func (d *dbServer) GetimageName(userid string) (string, string, error) {
// 	collection := d.sess.Collection(userCollection)
// 	res := collection.Find(db.Cond{"userid": userid})
// 	var user User
// 	err := res.One(&user)
// 	if err != nil {
// 		return "", "", err
// 	}
// 	spliter := strings.Split(user.Picture, "/")
// 	picName := spliter[3]
// 	return picName, user.Picture, nil
// }

func DownloadFile(writer http.ResponseWriter, request *http.Request) {
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
func (d *dbServer) TemperingEmail(w http.ResponseWriter, r *http.Request) {
	var contract Contract
	var signers []Signer

	_ = json.NewDecoder(r.Body).Decode(&contract)

	contractCollection := d.sess.Collection(ContractCollection)
	Signercollection := d.sess.Collection(SignerCollection)

	res := contractCollection.Find(db.Cond{"ContractID": contract.ContractID})
	res.One(&contract)

	res2 := Signercollection.Find(db.Cond{"ContractID": contract.ContractID, "CC": 0})
	res2.All(&signers)

	var cd ContractDetail
	cd.ContractData = contract
	cd.Signers = signers

	resbool, emails := d.Getemails(cd)
	if !resbool {
		RenderResponse(w, "CAN NOT SEND EMAIL TO THE RECIPIENTS", http.StatusOK)
		return
	}

	for _, index := range emails {
		go Email.TemperMail(index, "CONTRACT TEMPERED", contract.ContractID)
	}
	RenderResponse(w, "EMAIL SENT", http.StatusOK)
	return
}
