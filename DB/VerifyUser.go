package DB

import (
	"ENOTARY-Server/Email"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	db "upper.io/db.v3"
)

//TODO here : account verification with email not userid

func (d *dbServer) AccountVerif(w http.ResponseWriter, r *http.Request) {
	var VU VerifUser
	var temp VerifUser

	_ = json.NewDecoder(r.Body).Decode(&temp)

	Collection := d.sess.Collection(VerifCollection)
	res := Collection.Find(db.Cond{"userid": temp.UserID, "VerificationCode": temp.VerificationCode})
	err := res.One(&VU)
	if err != nil {
		RenderError(w, "INVALID_CODE")
		return
	}
	expTime, err := strconv.ParseInt(VU.ExpTime, 10, 64)
	if err != nil {
		RenderError(w, "INTERNAL ERROR TRY AGAIN")
		return
	}
	if expTime < time.Now().Unix() {
		RenderResponse(w, "VERIFICATION CODE HAS EXPIRED", http.StatusOK)
		return
	}
	d.sess.Query(`Update Users set verified ="1" where userID= ?`, temp.UserID)
	RenderResponse(w, "USER EMAIL VERIFIED SUCCESSFULLY ", http.StatusOK)
	return
}

func (d *dbServer) ResendCode(w http.ResponseWriter, r *http.Request) {
	var user User
	var temp User
	_ = json.NewDecoder(r.Body).Decode(&temp)

	// find user user id using his email
	UColletion := d.sess.Collection(UserCollection)
	res1 := UColletion.Find(db.Cond{"email": temp.Email})
	err := res1.One(&user)
	if err != nil {
		RenderError(w, "INTERNAL ERROR (USER NOT FOUND")
		return
	}
	// update user verification code and exp using his userid
	Collection := d.sess.Collection(VerifCollection)
	res := Collection.Find(db.Cond{"userID": user.Userid})
	vcode := GenerateToken(3)
	fmt.Println(vcode)
	res.Update(map[string]string{
		"VerificationCode": vcode,
		"expTime":          strconv.FormatInt(time.Now().Add(2*time.Hour).Unix(), 10),
	})
	_, err = Email.SendMail(user.Email, vcode)
	if err != nil {
		RenderError(w, "CAN NOT SEND MAIL TRY AGAIN")
		return
	}
	RenderResponse(w, "NEW CODE SENT TO YOUR EMAIL", http.StatusOK)
	return
}
