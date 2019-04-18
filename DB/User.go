package DB

import (
	"encoding/json"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	db "upper.io/db.v3"
)

func (d *dbServer) UploadSign(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxpicSize)
	err := r.ParseMultipartForm(5000)
	if err != nil {
		RenderError(w, "FILE SHOULD BE LESS THAN 5 MB")
		Logger("FILE SHOULD BE LESS THAN 5 MB")
		return
	}

	//get user id from JWT
	tokenstring := r.Header["Token"][0]
	claims, cBool := GetClaims(tokenstring)
	if !cBool {
		RenderError(w, "Invalid user request")
		Logger("Invalid user request")
		return
	}
	uID := claims["userid"].(string)

	f, _, err := r.FormFile("userSign")
	if err != nil {
		RenderError(w, "INVALID_SIGN_FILE")
		Logger("Invalid file upload")
		return
	}
	defer f.Close()

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		RenderError(w, "INVALID_SIGN_FILE")
		Logger("Invalid file upload")
		return
	}
	filetype := http.DetectContentType(bs)
	if filetype != "image/jpeg" && filetype != "image/jpg" && filetype != "image/bmp" &&
		filetype != "image/gif" && filetype != "image/png" {
		RenderError(w, "INVALID_FILE_TYPE_UPLOAD jpeg,jpg,png OR gif")
		Logger("Invalid file upload")
		return
	}

	fileEndings, err := mime.ExtensionsByType(filetype)
	if err != nil {
		RenderError(w, "INVALID_SIGN_FILE_TYPE")
		Logger("Invalid file upload " + uID)
		return
	}
	_ = d.removeOldSign(uID)

	f1, _, err := r.FormFile("userInitail")
	if err != nil {
		RenderError(w, "INVALID_INITAILS_FILE")
		Logger("Invalid file upload")
		return
	}
	bs1, err := ioutil.ReadAll(f1)
	if err != nil {
		RenderError(w, "INVALID_INITAILS_FILE")
		Logger("Invalid file upload")
		return
	}
	filetype1 := http.DetectContentType(bs1)
	if filetype1 != "image/jpeg" && filetype1 != "image/jpg" && filetype1 != "image/bmp" &&
		filetype1 != "image/gif" && filetype1 != "image/png" {
		RenderError(w, "INVALID_FILE_TYPE_UPLOAD jpeg,jpg,png OR gif")
		Logger("Invalid file upload")
		return
	}
	fileEndings1, err := mime.ExtensionsByType(filetype1)
	if err != nil {
		RenderError(w, "INVALID_INITAILS_FILE_TYPE")
		Logger("Invalid file upload " + uID)
		return
	}
	_ = d.removeOldInitial(uID)

	signdata := string(bs)
	initaildata := string(bs1)
	signpath := filepath.Join(Signpath, uID+fileEndings[0])
	initialspath := filepath.Join(InitialsPath, uID+fileEndings1[0])
	file, err := os.Create(signpath)
	file1, err1 := os.Create(initialspath)

	if err != nil {
		RenderError(w, "CAN NOT CREATE NEW SIGN ")
		Logger("CAN NOT CREATE NEW SIGN " + uID)
		return
	}
	if err1 != nil {
		RenderError(w, "CAN NOT CREATE NEW INITAILS")
		Logger("CAN NOT UPDATE INITAILS " + uID)
		return
	}

	defer file.Close()
	defer file1.Close()
	file.WriteString(signdata)
	file1.WriteString(initaildata)

	d.updateInitialpath(uID, initialspath)
	d.updateSignpath(uID, signpath)

	var signRes SignRes
	signRes.InitialsPath = initialspath
	signRes.Signpath = signpath

	json.NewEncoder(w).Encode(signRes)
	Logger("Sign updated user: " + uID)
	return

}

func (d *dbServer) updateSignpath(userid string, path string) bool {
	collection := d.sess.Collection(UserCollection)
	res := collection.Find(db.Cond{"userid": userid})
	res.Update(map[string]string{
		"sign": path,
	})
	return true
}

func (d *dbServer) updateInitialpath(userid string, path string) bool {
	collection := d.sess.Collection(UserCollection)
	res := collection.Find(db.Cond{"userid": userid})
	res.Update(map[string]string{
		"initials": path,
	})
	return true
}

//function to update user profile pic
func (d *dbServer) ProfilePic(w http.ResponseWriter, r *http.Request) {

	r.Body = http.MaxBytesReader(w, r.Body, MaxpicSize)
	err := r.ParseMultipartForm(5000)
	if err != nil {
		RenderError(w, "FILE SHOULD BE LESS THAN 5 MB")
		Logger("FILE SHOULD BE LESS THAN 5 MB")
		return
	}

	//get user id from JWT
	tokenstring := r.Header["Token"][0]
	claims, cBool := GetClaims(tokenstring)
	if !cBool {
		RenderError(w, "Invalid user request")
		Logger("Invalid user request")
		return
	}
	uID := claims["userid"].(string)

	//get file from form
	f, _, err := r.FormFile("userfile")
	if err != nil {
		RenderError(w, "INVALID_FILE")
		Logger("Invalid file upload")
		return
	}
	defer f.Close()

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		RenderError(w, "INVALID_FILE")
		Logger("Invalid file upload")
		return
	}
	filetype := http.DetectContentType(bs)
	if filetype != "image/jpeg" && filetype != "image/jpg" && filetype != "image/bmp" &&
		filetype != "image/gif" && filetype != "image/png" {
		RenderError(w, "INVALID_FILE_TYPE_UPLOAD jpeg,jpg,png OR gif")
		Logger("Invalid file upload")
		return
	}

	fileEndings, err := mime.ExtensionsByType(filetype)
	if err != nil {
		RenderError(w, "INVALID_FILE")
		Logger("Invalid file upload")
		return
	}

	// remove users old picture before adding new one
	rop := d.removeOldPic(uID)
	if !rop {
		RenderError(w, "CAN NOT UPDATE PICTURE TRY AGAIN")
		Logger("CAN NOT UPDATE PICTURE")
		return
	}

	picdata := string(bs)
	newpath := filepath.Join(Profilepicspath, uID+fileEndings[0])
	file, err := os.Create(newpath)

	if err != nil {
		RenderError(w, "INVALID_FILE ")
		Logger("CAN NOT UPDATE PICTURE")
		return
	}

	defer file.Close()
	file.WriteString(picdata)
	d.updatePicPath(uID, newpath)

	RenderResponse(w, newpath, http.StatusOK)
	Logger("Profile pic updated | userid = " + uID)
	return
}

func (d *dbServer) RemovePic(w http.ResponseWriter, r *http.Request) {
	tokenstring := r.Header["Token"][0]
	claims, cBool := GetClaims(tokenstring)
	if !cBool {
		RenderError(w, "Invalid user request")
		Logger("Invalid user request")
		return
	}
	uID := claims["userid"].(string)

	resbool := d.removeOldPic(uID)

	if !resbool {
		RenderResponse(w, "INTERNAL ERROR TRY AGAIN", http.StatusOK)
		Logger("Remove old pic error " + uID)
		return
	}
	picbool := d.updatePicPath(uID, Def_pic_path)
	if !picbool {
		RenderResponse(w, "INTERNAL ERROR TRY AGAIN", http.StatusOK)
		Logger("Remove old pic error " + uID)
		return
	}
	RenderResponse(w, Def_pic_path, http.StatusOK)
	return
}

//function to update user's profile pic path in DB
func (d *dbServer) updatePicPath(userid string, picpath string) bool {
	collection := d.sess.Collection(UserCollection)
	res := collection.Find(db.Cond{"userid": userid})
	res.Update(map[string]string{
		"picture": picpath,
	})
	return true
}

// func (d *dbServer) removeOldSign(userid string) bool {
// 	collection := d.sess.Collection(UserCollection)
// 	res := collection.Find(db.Cond{"userid": userid})
// 	var user User
// 	err := res.One(&user)
// 	if err != nil {
// 		return false
// 	}
// 	spliter := strings.Split(user.Picture, "/")
// 	picName := spliter[3]
// 	err = os.Remove(user.Picture)
// 	if err != nil {
// 		return false
// 	}
// 	return true
// }

// function to remove user's privious profile pic
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

func (d *dbServer) removeOldSign(userid string) bool {
	collection := d.sess.Collection(UserCollection)
	res := collection.Find(db.Cond{"userid": userid})
	var user User
	err := res.One(&user)
	if err != nil {
		return false
	}

	err = os.Remove(user.Sign)
	if err != nil {
		return false
	}

	return true
}
func (d *dbServer) removeOldInitial(userid string) bool {
	collection := d.sess.Collection(UserCollection)
	res := collection.Find(db.Cond{"userid": userid})
	var user User
	err := res.One(&user)
	if err != nil {
		return false
	}

	err = os.Remove(user.Initials)
	if err != nil {
		return false
	}

	return true
}

//func to update password when user forgets it
func (d *dbServer) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var passrec Passrecovery
	_ = json.NewDecoder(r.Body).Decode(&passrec)

	resBool, err := d.VerfUser(passrec.Email, passrec.Vcode, false)
	if !resBool {
		RenderError(w, err)
		Logger(err)
		return
	}

	Collection := d.sess.Collection(UserCollection)
	res := Collection.Find(db.Cond{"email": passrec.Email})
	resCheck, _ := res.Count()

	if resCheck != 1 {
		RenderError(w, "INVALID USER TRY AGAIN")
		Logger("INVALID USER " + passrec.Email)
		return
	}

	pasres, err := VerifyPassword(passrec.Pass)
	if !pasres {
		RenderError(w, err)
		Logger(err)
		return
	}

	res.Update(map[string]string{
		"password": passrec.Pass,
	})

	RenderResponse(w, "PASSWORD UPDATED SUCCESSFULLY", http.StatusOK)
	Logger("PASSWORD UPDATED " + passrec.Email)
	return
}

func (d *dbServer) Userpreferences(w http.ResponseWriter, r *http.Request) {
	tokenstring := r.Header["Token"][0]
	claims, cBool := GetClaims(tokenstring)
	if !cBool {
		RenderError(w, "Invalid user request")
		Logger("Invalid user request")
		return
	}
	uID := claims["userid"].(string)

	var prefs Preferences
	var user User
	_ = json.NewDecoder(r.Body).Decode(&prefs)

	Collection := d.sess.Collection(UserCollection)
	res := Collection.Find(db.Cond{"userid": uID})
	err := res.One(&user)

	if err != nil {
		RenderResponse(w, "Data not saved try again", http.StatusOK)
		return
	}

	user.Name = prefs.UserName
	user.Company = prefs.Company
	user.Phone = prefs.Phone

	resbool, errstring := CredentialValidation(user)
	if !resbool {
		RenderError(w, errstring)
		return
	}

	res.Update(user)

	RenderResponse(w, "Updated successfully", http.StatusOK)
	Logger("user data updated :" + uID)
	return
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
