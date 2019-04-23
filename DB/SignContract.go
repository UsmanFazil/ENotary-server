package DB

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"strings"

	"upper.io/db.v3"
)

func (d *dbServer) SaveCoordinates(w http.ResponseWriter, r *http.Request) {
	var pi []PlaygroundInput
	_ = json.NewDecoder(r.Body).Decode(&pi)

	//Collection := d.sess.Collection(CoordinatesCol)

	var sc = make([]Coordinates, len(pi))

	for i := 0; i < len(pi); i++ {
		sc[i].ContractID = pi[i].Contractid
		sc[i].UserID = pi[i].Recipient
		sc[i].Name = pi[i].Text
		sc[i].Topcord = pi[i].Top
		sc[i].Leftcord = pi[i].Left

		q := d.sess.InsertInto("Coordinates").Columns("ContractID", "userID", "name", "topcord", "leftcord").Values(sc[i].ContractID, sc[i].UserID, sc[i].Name, sc[i].Topcord, sc[i].Leftcord)
		_, err := q.Exec()

		if err != nil {
			RenderError(w, "CAN NOT UPDATE SIGNERS COORDINATES")
			Logger("Can't add coordinates, ContractID :" + sc[0].ContractID)
			fmt.Println(err)
			return
		}
	}
	json.NewEncoder(w).Encode(sc)
	Logger("Signer cordinates added ContractID :" + sc[0].ContractID)
	return
}

func (d *dbServer) ServeCoordinates(w http.ResponseWriter, r *http.Request) {
	tokenstring := r.Header["Token"][0]
	claims, cBool := GetClaims(tokenstring)
	if !cBool {
		RenderError(w, "Invalid user request")
		Logger("Invalid user request")
		return
	}
	uID := claims["userid"].(string)

	var contract Contract
	var cords []Coordinates
	_ = json.NewDecoder(r.Body).Decode(&contract)

	Collection := d.sess.Collection(CoordinatesCol)

	res := Collection.Find(db.Cond{"UserID": uID, "ContractID": contract.ContractID})
	err := res.All(&cords)

	if err != nil {
		RenderError(w, "CAN NOT FIND COORDINATES")
		return
	}

	json.NewEncoder(w).Encode(cords)

}

func (d *dbServer) SignContract(w http.ResponseWriter, r *http.Request) {

	var sc SignContract
	var contract Contract

	_ = json.NewDecoder(r.Body).Decode(&sc)

	Collection := d.sess.Collection(ContractCollection)
	res := Collection.Find(db.Cond{"ContractID": sc.ContractID})
	err := res.One(&contract)

	if err != nil {
		fmt.Println("contract not found")
		return
	}

	path := Contractfilepath + "/" + contract.ContractID + ".png"
	Save(sc.FileBase64, path)

	contract.Filepath = path

	res.Update(contract)

}

func Save(data string, path string) {
	idx := strings.Index(data, ";base64,")
	if idx < 0 {
		panic("InvalidImage")
	}
	ImageType := data[11:idx]
	log.Println(ImageType)

	unbased, err := base64.StdEncoding.DecodeString(data[idx+8:])
	if err != nil {
		panic("Cannot decode b64")
	}
	x := bytes.NewReader(unbased)
	switch ImageType {
	case "png":
		im, err := png.Decode(x)
		if err != nil {
			panic("Bad png")
		}

		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			panic("Cannot open file")
		}

		png.Encode(f, im)
	case "jpeg":
		im, err := jpeg.Decode(x)
		if err != nil {
			panic("Bad jpeg")
		}

		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			panic("Cannot open file")
		}

		jpeg.Encode(f, im, nil)
	case "gif":
		im, err := gif.Decode(x)
		if err != nil {
			panic("Bad gif")
		}

		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			panic("Cannot open file")
		}

		gif.Encode(f, im, nil)
	}
}

// usersignpath := filepath.Join(Signpath, uID)
// userinitialspath := filepath.Join(InitialsPath, uID)

// resbool, _ := saveFile(input.SignBase64, "./Files")

// if !resbool {
// 	fmt.Println("error")
// }

// if resbool {
// 	d.updateSignpath(uID, usersignpath+signExt)
// }

// resbool1, initialExt := saveFile(input.InitialsBase64, userinitialspath)
// if resbool1 {
// 	d.updateInitialpath(uID, userinitialspath+initialExt)
// }
// var signRes SignRes
// signRes.InitialsPath = userinitialspath + initialExt
// signRes.Signpath = usersignpath + signExt

// json.NewEncoder(w).Encode(signRes)
// Logger("Sign updated user: " + uID)

// var sc SignContract
//var contract Contract
//var signer Signer
//var signers []Signer
// _ = json.NewDecoder(r.Body).Decode(&sc)

// saveFile(sc.ContractID, "./Files")

// tokenstring := r.Header["Token"][0]
// claims, cBool := GetClaims(tokenstring)
// if !cBool {
// 	RenderError(w, "Invalid user request")
// 	Logger("Invalid user request")
// 	return
// }
// uID := claims["userid"].(string)

// contractCollection := d.sess.Collection(ContractCollection)
// signerCollection := d.sess.Collection(SignerCollection)

// res := contractCollection.Find(db.Cond{"ContractID": sc.ContractID})
// res.One(&contract)
// total, _ := res.Count()
// if total != 1 {
// 	RenderError(w, "CONTRACT NOT FOUND")
// 	return
// }

// res1 := signerCollection.Find(db.Cond{"ContractID": sc.ContractID, "userID": uID})
// total1, _ := res1.Count()
// if total1 != 1 {
// 	RenderError(w, "CONTRACT NOT FOUND")
// 	return
// }
// res1.One(&signer)

// if signer.SignStatus != "pending" {
// 	RenderError(w, "Contract Already signed by User")
// 	return
// }
// signer.SignStatus = "Signed"
// signer.SignDate = time.Now().Format(time.RFC850)
// res.Update(signer)

// newpath := filepath.Join(Contractfilepath, contract.ContractID)
// err, ext := saveFile(sc.FileBase64, newpath)
// if !err {
// 	fmt.Println(err)
// 	return
// }
// contract.Filepath = newpath + ext
// res.Update(contract)

// fmt.Println(contract.Filepath)
