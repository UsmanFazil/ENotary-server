package DB

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	db "upper.io/db.v3"
)

const userCollection = "Users"

type User struct {
	Userid       string `db:"userid"`
	Email        string `db:"email"`
	Password     string `db:"password"`
	Name         string `db:"name"`
	Company      string `db:"company"`
	Phone        string `db:"phone"`
	Picture      string `db:"picture"`
	Sign         string `db:"sign"`
	Initials     string `db:"initials"`
	Verified     int    `db:"verified"`
	CreationTime string `db:"creationTime"`
}

type LogCheck struct {
	Email    string `json:"email"`
	Password string `json:password`
}

type JwtToken struct {
	Token string `json:"token"`
}

type Exception struct {
	Message string `json:"message"`
}

// Login : Method to check weather the user exists on the system or not
func (d *dbServer) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
			"email":  user.Email,
		})
		tokenString, error := token.SignedString([]byte("secretkey"))
		if error != nil {
			fmt.Println(error)
		}
		json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
		return
	}

	json.NewEncoder(w).Encode("Invalid password")
}


// Newuser : Method to add new user into the system
func (d *dbServer) Newuser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)

	id, err := uuid.NewV4()
	if err != nil {
		json.NewEncoder(w).Encode("error")
		return
	}

	res, err := d.sess.InsertInto(userCollection).
		Values(id, user.Email, user.Password, user.Name, user.Company, user.Phone, "non", "non", "non", time.Now().Format(time.RFC850), 0).
		Exec()
	if err != nil {
		json.NewEncoder(w).Encode("error")
	}
	log.Println("new user entered in DB", res)

	json.NewEncoder(w).Encode("Successfull entered")

}

func (d *dbServer) Validateuser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	Collection := d.sess.Collection(userCollection)
	var user User
	res := Collection.Find(db.Cond{"email": params["email"]})

	err := res.One(&user)
	if err != nil {
		json.NewEncoder(w).Encode("invalid email")
		return
	}

	if user.Verified == 1 {
		json.NewEncoder(w).Encode("user already verified")
		return
	}
	user.Verified = 1
	if err := res.Update(user); err != nil {
		json.NewEncoder(w).Encode("internal error please try again")
	}
	json.NewEncoder(w).Encode("User verified successfully")

}
