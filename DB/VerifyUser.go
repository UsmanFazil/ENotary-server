package DB

import (
	"ENOTARY-Server/Email"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	db "upper.io/db.v3"
)

//TODO here : account verification with email not userid

func (d *dbServer) EmailVerification(w http.ResponseWriter, r *http.Request) {
	var VU VerifUser
	var temp EmailVerf
	var user User

	_ = json.NewDecoder(r.Body).Decode(&temp)

	userCol := d.sess.Collection(UserCollection)
	res := userCol.Find(db.Cond{"email": temp.Email})
	err := res.One(&user)

	if err != nil {
		RenderError(w, "INVALID USER")
		return
	}

	Collection := d.sess.Collection(VerifCollection)
	res1 := Collection.Find(db.Cond{"userid": user.Userid, "VerificationCode": temp.VerificationCode})
	errstring := res1.One(&VU)
	if errstring != nil {
		RenderError(w, "INVALID_CODE")
		return
	}
	expTime, err := strconv.ParseInt(VU.ExpTime, 10, 64)
	if err != nil {
		RenderError(w, "INTERNAL ERROR TRY AGAIN")
		return
	}
	if expTime < time.Now().Unix() {
		RenderError(w, "VERIFICATION CODE HAS EXPIRED")
		return
	}
	d.sess.Query(`Update Users set verified ="1" where userID= ?`, user.Userid)
	RenderResponse(w, "USER EMAIL VERIFIED SUCCESSFULLY ", http.StatusOK)
	return
}

func (d *dbServer) SendCode(w http.ResponseWriter, r *http.Request) {
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
	res.Update(map[string]string{
		"VerificationCode": vcode,
		"expTime":          strconv.FormatInt(time.Now().Add(2*time.Hour).Unix(), 10),
	})
	go Email.SendMail(user.Email, vcode)
	RenderResponse(w, "NEW CODE SENT TO YOUR EMAIL", http.StatusOK)
	return
}
