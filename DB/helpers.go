package DB

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	db "upper.io/db.v3"
)

func (d *dbServer) GetUser(email string) (*User, bool, error) {
	Collection := d.sess.Collection(userCollection)
	res := Collection.Find(db.Cond{"email": email})
	var user User
	err := res.One(&user)

	if err != nil {
		return nil, false, err
	}
	return &user, true, nil
}

// RenderError : creates error response
func RenderError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

// RenderResponse : creates response
func RenderResponse(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func GenerateToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func CredentialValidation(usr User) (bool, string) {
	mailRE, errstring := VerifyEmail(usr.Email)
	if !mailRE {
		return false, errstring
	}
	pswdRE, errstring := VerifyPassword(usr.Password)
	if !pswdRE {
		return false, errstring
	}
	nameRE, errstring := verifyName(usr.Name)
	if !nameRE {
		return false, errstring
	}
	cmpRE, errstring := verifyComp(usr.Company)
	if !cmpRE {
		return false, errstring
	}
	phRE, errstring := verifyPhone(usr.Phone)
	if !phRE {
		return false, errstring
	}

	return true, ""
}

func verifyPhone(ph string) (bool, string) {
	var len int
	for _, _ = range ph {
		len++
	}
	_, err := strconv.Atoi(ph)
	if err != nil || len < 8 || len > 20 {
		return false, "Invalid number"
	}
	return true, ""
}

func verifyName(name string) (bool, string) {
	nonNameRE := regexp.MustCompile("^[0-9!#$%&*+/=?^_`{|}~]+")
	if !(len(name) > 7) || nonNameRE.MatchString(name) {
		return false, "invalid name"
	}
	return true, ""
}

func verifyComp(company string) (bool, string) {
	compRE := regexp.MustCompile("^[a-zA-Z0-9]{3,50}")
	if !compRE.MatchString(company) {
		return false, "invalid company name"
	}
	return true, ""
}

func VerifyEmail(mail string) (bool, string) {
	mailRE := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !mailRE.MatchString(mail) {
		return false, "invalid email address"
	}
	return true, ""
}

func VerifyPassword(password string) (bool, string) {
	var uppercasePresent bool
	var lowercasePresent bool
	var numberPresent bool
	var space bool
	const minPassLength = 8
	const maxPassLength = 16
	var passLen int

	for _, ch := range password {
		passLen++
		switch {
		case unicode.IsNumber(ch):
			numberPresent = true
		case unicode.IsUpper(ch):
			uppercasePresent = true
		case unicode.IsLower(ch):
			lowercasePresent = true
		case ch == ' ':
			space = true
		}
	}

	if space {
		return false, "space fields are not allowed in password"
	}
	if !(minPassLength <= passLen && passLen <= maxPassLength) {
		return false, "password length must be between 8 to 16 characters long"
	}
	if !lowercasePresent {
		return false, "lowercase letter missing in password"
	}
	if !uppercasePresent {
		return false, "uppercase letter missing in password"
	}
	if !numberPresent {
		return false, "atleast one numeric character required in password"
	}

	return true, ""
}

func (d *dbServer) GetimageName(userid string) (string, string, error) {
	collection := d.sess.Collection(userCollection)
	res := collection.Find(db.Cond{"userid": userid})
	var user User
	err := res.One(&user)
	if err != nil {
		return "", "", err
	}
	spliter := strings.Split(user.Picture, "/")
	picName := spliter[3]
	return picName, user.Picture, nil
}