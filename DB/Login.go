package DB

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	db "upper.io/db.v3"
)

// Login : Method to check weather the user exists on the system or not ()

func (d *dbServer) Login(w http.ResponseWriter, r *http.Request) {
	var user User
	var logcheck LogCheck
	_ = json.NewDecoder(r.Body).Decode(&logcheck)

	emailvalid, errstring := VerifyEmail(logcheck.Email)
	if !emailvalid {
		RenderError(w, errstring)
		return
	}
	pswdvalid, _ := VerifyPassword(logcheck.Password)
	if !pswdvalid {
		RenderError(w, "INVALID PASSWORD")
		Logger("INVALID PASSWORD")
		return
	}
	Collection := d.sess.Collection(UserCollection)
	res := Collection.Find(db.Cond{"email": logcheck.Email})

	err := res.One(&user)
	if err != nil {
		RenderError(w, "INVALID EMAIL")
		Logger("INVALID EMAIL")
		return
	}

	if logcheck.Password == user.Password {
		if user.Verified == 0 {
			RenderResponse(w, "Please verify your email first", http.StatusOK)
			Logger("Verify email first")
			return
		}

		//create JW Token for the user
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userid": user.Userid,
			"exp":    time.Now().Add(time.Minute * 120).Unix(),
			"iat":    time.Now().Unix(),
			"iss":    "ENotary",
		})
		tokenString, error := token.SignedString([]byte("secretkey"))
		if error != nil {
			RenderError(w, "INTERNAL ERROR TRY AGAIN")
			Logger("INTERNAL ERROR")
			return
		}

		//total number of contracts (waiting for me & waiting for others)
		waitingOther, err := d.WaitingforOther(user.Userid)
		if err != nil {
			RenderError(w, "INTERNAL ERROR TRY AGAIN")
			Logger("INTERNAL DB ERROR")
			return
		}

		waitingMe, err := d.WaitingforMe(user.Userid)
		if err != nil {
			RenderError(w, "INTERNAL ERROR TRY AGAIN")
			Logger("INTERNAL DB ERROR")
			return
		}

		Expiring, errbool := d.ExpiringSoon(user.Userid)
		if errbool != nil {
			RenderError(w, "INTERNAL ERROR TRY AGAIN")
			Logger("INTERNAL DB ERROR")
			return
		}

		//remove user password from data struct
		user.Password = ""
		data := LoginStruct{Userdata: user, WaitingME: waitingMe, WaitingOther: waitingOther, ExpiringSoon: Expiring, Token: tokenString}

		Logger("New Login" + user.Userid)
		json.NewEncoder(w).Encode(data)
		return
	}

	RenderError(w, "INVALID PASSWORD")
	Logger("INVALID PASSWORD" + logcheck.Email)
	return
}
