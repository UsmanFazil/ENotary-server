package DB

import (
	"encoding/json"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

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
		RenderResponse(w, "EMAIL_ALREADY_EXISTS", http.StatusOK)
		return
	}
	verify, errmsg := CredentialValidation(user)
	if !verify {
		RenderError(w, errmsg)
		return
	}
	_, err = d.sess.InsertInto(userCollection).
		Values(id, user.Email, user.Password, user.Name, user.Company, user.Phone, "non", "non", "non", time.Now().Format(time.RFC850), 0).
		Exec()
	if err != nil {
		RenderError(w, "CAN_NOT_GENERATE_USER_ID")
		return
	}
	RenderResponse(w, "ACCOUNT_CREATED_PLEASE_VALIDATE_YOUR_EMAIL_TO_LOGIN", http.StatusOK)
}
