package DB

import (
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

	f, _, err := r.FormFile("userfile")
	userid := r.FormValue("userid")
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

	fileName := userid
	fileEndings, err := mime.ExtensionsByType(filetype)
	if err != nil {
		RenderError(w, "INVALID_FILE")
		return
	}

	// remove users old picture before adding new one
	rop := d.removeOldPic(userid)
	if !rop {
		RenderError(w, "CAN NOT UPDATE PICTURE TRY AGAIN")
		return
	}

	s = string(bs)
	newpath := filepath.Join(Profilepicspath, fileName+fileEndings[0])
	file, err := os.Create(newpath)

	if err != nil {
		RenderError(w, "INVALID_FILE ")
		return
	}
	defer file.Close()
	file.WriteString(s)
	d.updatePicPath(userid, newpath)

	RenderResponse(w, "FILE UPLOADED SUCCESSFULY", http.StatusOK)
}

func (d *dbServer) updatePicPath(userid string, picpath string) bool {
	collection := d.sess.Collection(userCollection)
	res := collection.Find(db.Cond{"userid": userid})
	res.Update(map[string]string{
		"picture": picpath,
	})
	return true
}

func (d *dbServer) removeOldPic(userid string) bool {
	collection := d.sess.Collection(userCollection)
	res := collection.Find(db.Cond{"userid": userid})
	var user User
	err := res.One(&user)
	if err != nil {
		return false
	}

	spliter := strings.Split(user.Picture, "/")
	picName := spliter[2]

	if picName != "default.png" {
		err = os.Remove(user.Picture)
		if err != nil {
			return false
		}
	}
	return true

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
