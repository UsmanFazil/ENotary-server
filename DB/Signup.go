package DB

import (
	"ENOTARY-Server/Email"
	"encoding/json"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

func (d *dbServer) Signup(w http.ResponseWriter, r *http.Request) {
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)

	// check weather email already exists or not
	_, exists, _ := d.GetUser(user.Email)
	if exists == true {
		RenderError(w, "EMAIL_ALREADY_EXISTS")
		Logger("EMAIL_ALREADY_EXISTS" + user.Email)
		return
	}

	// verify user input with regex
	verify, errmsg := CredentialValidation(user)
	if !verify {
		RenderError(w, errmsg)
		Logger(errmsg + " " + user.Email)
		return
	}

	// create new userid
	id, err := uuid.NewV4()
	if err != nil {
		RenderError(w, "CAN_NOT_GENERATE_USER_ID")
		Logger("UUID issue at signup")
		return
	}
	// TODO : create user struct and use "collection.Insert(user)"

	user.Userid = id.String()
	user.Picture = Def_pic_path
	user.CreationTime = time.Now().Format(time.RFC850)
	user.Verified = 0
	user.Sign = "default"
	user.Initials = "default"

	Collection := d.sess.Collection(UserCollection)
	_, errstrng := Collection.Insert(user)
	if errstrng != nil {
		RenderError(w, "CAN_NOT_GENERATE_USER_ID_TRY_AGAIN")
		Logger("INTERNAL DB ERROR")
		return
	}

	// create a verification code and store in db
	verfCode := GenerateToken(3)
	svc := d.InsertVerfCode(id.String(), verfCode)
	if !svc {
		RenderError(w, "USER CREATED BUT CAN NOT GENERATE VERIFICATION EMAIL TRY LOGIN")
		Logger("INTERNAL DB ERROR")
		return
	}

	// send verification code email to user
	go Email.SendMail(user.Email, verfCode)
	RenderResponse(w, "YOUR ACCOUNT HAS BEEN CREATED AND A VERIFICATION EMAIL HAS BEEN SENT TO YOUR EMAIL ADDRESS", http.StatusOK)
	Logger("NEW USER SIGNUP " + user.Email)
}

func (d *dbServer) InsertVerfCode(userid string, verfcode string) bool {
	_, err := d.sess.InsertInto(VerifCollection).Values(userid, verfcode, time.Now().Add(2*time.Hour).Unix()).
		Exec()
	if err != nil {
		return false
	}
	return true
}
