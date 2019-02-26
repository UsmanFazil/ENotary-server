package DB

import (
	"ENOTARY-Server/Email"
	"encoding/json"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

const def_pic_path = "Files/Profile_pics/default.jpeg"

func (d *dbServer) Signup(w http.ResponseWriter, r *http.Request) {
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	mailID := user.Email

	id, err := uuid.NewV4()
	if err != nil {
		RenderError(w, "CAN_NOT_GENERATE_USER_ID")
		return
	}

	_, exists, _ := d.GetUser(mailID)

	if exists == true {
		RenderError(w, "EMAIL_ALREADY_EXISTS")

		return
	}
	verify, errmsg := CredentialValidation(user)
	if !verify {
		RenderError(w, errmsg)
		return
	}
	// TODO : create user struct and use "collection.Insert(user)"
	_, err = d.sess.InsertInto(UserCollection).
		Values(id, user.Email, user.Password, user.Name, user.Company, user.Phone, def_pic_path, "non", "non", time.Now().Format(time.RFC850), 0).
		Exec()
	if err != nil {
		RenderError(w, "CAN_NOT_GENERATE_USER_ID_TRY_AGAIN")
		return
	}
	verfCode := GenerateToken(3)

	//save verification code in DB
	svc := d.InsertVerfCode(id.String(), verfCode)
	if !svc {
		RenderError(w, "USER CREATED BUT CAN NOT GENERATE VERIFICATION EMAIL TRY LOGIN")
		return
	}

	_, err = Email.SendMail(user.Email, verfCode)
	if err != nil {
		RenderError(w, "ACCOUNT GENERATED BUT CAN NOT GENERATE VERIFICATION EMAIL TRY LOGIN")
		return
	}
	RenderResponse(w, "YOUR ACCOUNT HAS BEEN CREATED AND A VERIFICATION EMAIL HAS BEEN SENT TO YOUR EMAIL ADDRESS", http.StatusOK)
}

func (d *dbServer) InsertVerfCode(userid string, verfcode string) bool {
	_, err := d.sess.InsertInto(VerifCollection).Values(userid, verfcode, time.Now().Add(2*time.Hour).Unix()).
		Exec()
	if err != nil {
		return false
	}
	return true
}
