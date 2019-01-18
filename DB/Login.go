package DB

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	db "upper.io/db.v3"
)

type LoginStruct struct {
	Userdata     User
	WaitingME    uint64
	WaitingOther uint64
	Token        string
}

type LogCheck struct {
	Email    string `json:"email"`
	Password string `json:password`
}

const userCollection = "Users"

// Login : Method to check weather the user exists on the system or not ()

func (d *dbServer) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	var user User
	var logcheck LogCheck
	_ = json.NewDecoder(r.Body).Decode(&logcheck)

	Collection := d.sess.Collection(userCollection)
	res := Collection.Find(db.Cond{"email": logcheck.Email})

	err := res.One(&user)
	if err != nil {
		json.NewEncoder(w).Encode("Invalid email")
		return
	}

	if logcheck.Password == user.Password {

		if user.Verified == 0 {
			json.NewEncoder(w).Encode("Please verify your email first")
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
			// create a proper error response here
			return

		}
		waitingOther, err := d.WaitingforOther(user.Userid)
		if err != nil {
			// create a proper error response here
			return
		}
		waitingMe, err := d.WaitingforMe(user.Userid)
		if err != nil {
			// create a proper error response here
			return
		}
		data := LoginStruct{Userdata: user, WaitingME: waitingMe, WaitingOther: waitingOther, Token: tokenString}
		json.NewEncoder(w).Encode(data)
		return
	}
	json.NewEncoder(w).Encode("Invalid password")
}
