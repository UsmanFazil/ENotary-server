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
	var temp EmailVerf
	_ = json.NewDecoder(r.Body).Decode(&temp)

	resBool, err := d.VerfUser(temp.Email, temp.VerificationCode, true)

	if !resBool {
		RenderError(w, err)
		return
	}

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

func (d *dbServer) VerfUser(Email string, vcode string, access bool) (bool, string) {
	var user User
	var VU VerifUser
	userCol := d.sess.Collection(UserCollection)
	res := userCol.Find(db.Cond{"email": Email})
	err := res.One(&user)

	if err != nil {
		return false, "INVALID USER"
	}

	Collection := d.sess.Collection(VerifCollection)
	res1 := Collection.Find(db.Cond{"userid": user.Userid, "VerificationCode": vcode})
	errstring := res1.One(&VU)
	if errstring != nil {
		return false, "INVALID_CODE"
	}

	expTime, err := strconv.ParseInt(VU.ExpTime, 10, 64)
	if err != nil {
		return false, "INTERNAL ERROR TRY AGAIN"
	}
	if expTime < time.Now().Unix() {
		return false, "VERIFICATION CODE HAS EXPIRED"
	}

	if access {
		res.Update(map[string]int{
			"verified": 1,
		})
		return true, ""
	}

	return true, ""
}
