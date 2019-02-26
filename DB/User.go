package DB

import (
	"encoding/json"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	db "upper.io/db.v3"
)

//TODO HERE : get userid from the token
func (d *dbServer) ProfilePic(w http.ResponseWriter, r *http.Request) {

	var s string
	r.Body = http.MaxBytesReader(w, r.Body, MaxpicSize)
	err := r.ParseMultipartForm(5000)
	if err != nil {
		RenderError(w, "FILE SHOULD BE LESS THAN 5 MB")
		return
	}

	tokenstring := r.Header["Token"][0]
	claims, cBool := GetClaims(tokenstring)
	if !cBool {
		RenderError(w, "Invalid user request")
		return
	}
	uID := claims["userid"].(string)

	f, _, err := r.FormFile("userfile")
	if err != nil {
		RenderError(w, "INVALID_FILE")
		return
	}
	defer f.Close()

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		RenderError(w, "INVALID_FILE")
		return
	}
	filetype := http.DetectContentType(bs)
	if filetype != "image/jpeg" && filetype != "image/jpg" &&
		filetype != "image/gif" && filetype != "image/png" {
		RenderError(w, "INVALID_FILE_TYPE_UPLOAD jpeg,jpg,png OR gif")
		return
	}

	fileEndings, err := mime.ExtensionsByType(filetype)
	if err != nil {
		RenderError(w, "INVALID_FILE")
		return
	}

	// remove users old picture before adding new one
	rop := d.removeOldPic(uID)
	if !rop {
		RenderError(w, "CAN NOT UPDATE PICTURE TRY AGAIN")
		return
	}

	s = string(bs)
	newpath := filepath.Join(Profilepicspath, uID+fileEndings[0])

	file, err := os.Create(newpath)

	if err != nil {
		RenderError(w, "INVALID_FILE ")
		return
	}
	defer file.Close()
	file.WriteString(s)
	d.updatePicPath(uID, newpath)

	RenderResponse(w, "FILE UPLOADED SUCCESSFULY", http.StatusOK)
}

func (d *dbServer) updatePicPath(userid string, picpath string) bool {
	collection := d.sess.Collection(UserCollection)
	res := collection.Find(db.Cond{"userid": userid})
	res.Update(map[string]string{
		"picture": picpath,
	})
	return true
}

func (d *dbServer) removeOldPic(userid string) bool {
	collection := d.sess.Collection(UserCollection)
	res := collection.Find(db.Cond{"userid": userid})
	var user User
	err := res.One(&user)
	if err != nil {
		return false
	}

	spliter := strings.Split(user.Picture, "/")
	picName := spliter[2]

	if picName != "default.jpeg" {
		err = os.Remove(user.Picture)
		if err != nil {
			return false
		}
	}
	return true

}

func (d *dbServer) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var temp EmailVerf
	var user User
	var VU VerifUser

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
		RenderResponse(w, "VERIFICATION CODE HAS EXPIRED", http.StatusOK)
		return
	}

}

// picname, picpath, errstring := d.GetimageName(userid)
// if errstring != nil {
// 	RenderError(w, "CAN NOT REPLACE PICTURE TRY LATER")
// 	return
// }

// if picname != "default.png" {
// 	err := os.Remove(picpath)
// 	if err != nil {
// 		RenderError(w, "CAN NOT REPLACE PICTURE TRY LATER")
// 		return
// 	}
// }
