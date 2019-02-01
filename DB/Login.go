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
		return
	}
	Collection := d.sess.Collection(userCollection)
	res := Collection.Find(db.Cond{"email": logcheck.Email})

	err := res.One(&user)
	if err != nil {
		RenderError(w, "INVALID EMAIL")
		return
	}

	if logcheck.Password == user.Password {
		if user.Verified == 0 {
			RenderResponse(w, "Please verify your email first", http.StatusOK)
			return
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userid": user.Userid,
			"exp":    time.Now().Add(time.Minute * 10).Unix(),
			"iat":    time.Now().Unix(),
			"iss":    "ENotary",
		})
		tokenString, error := token.SignedString([]byte("secretkey"))
		if error != nil {
			RenderError(w, "INTERNAL ERROR TRY AGAIN")
			return
		}
		waitingOther, err := d.WaitingforOther(user.Userid)
		if err != nil {
			RenderError(w, "INTERNAL ERROR TRY AGAIN")
			return
		}
		waitingMe, err := d.WaitingforMe(user.Userid)
		if err != nil {
			RenderError(w, "INTERNAL ERROR TRY AGAIN")
			return
		}
		data := LoginStruct{Userdata: user, WaitingME: waitingMe, WaitingOther: waitingOther, Token: tokenString}
		json.NewEncoder(w).Encode(data)
		return
	}
	RenderError(w, "INVALID PASSWORD")
}
