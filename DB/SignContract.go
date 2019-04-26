package DB

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"strings"
	"time"

	"upper.io/db.v3"
)

var Coordinate []Coordinates

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

	Coordinate = append(Coordinate, sc...)
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

	for _, value := range Coordinate {
		if value.ContractID == contract.ContractID && value.UserID == uID {
			cords = append(cords, value)
		}
	}

	fmt.Println(cords)

	json.NewEncoder(w).Encode(cords)

}

func (d *dbServer) DeclineContract(w http.ResponseWriter, r *http.Request) {
	tokenstring := r.Header["Token"][0]
	claims, cBool := GetClaims(tokenstring)
	if !cBool {
		RenderError(w, "Invalid user request")
		Logger("Invalid user request")
		return
	}
	uID := claims["userid"].(string)

	var contract Contract
	var signer Signer

	_ = json.NewDecoder(r.Body).Decode(&contract)

	signerCol := d.sess.Collection(SignerCollection)
	Collection := d.sess.Collection(ContractCollection)
	res := Collection.Find(db.Cond{"ContractID": contract.ContractID})
	err := res.One(&contract)

	if err != nil {
		fmt.Println("contract not found")
		return
	}

	res1 := signerCol.Find(db.Cond{"ContractID": contract.ContractID, "userID": uID})
	errstring := res1.One(&signer)

	if errstring != nil {
		RenderError(w, "User Not found")
		return
	}
	signer.SignStatus = "Declined"
	contract.Status = "Declined"
	contract.UpdateTime = time.Now().Format(time.RFC850)

	err = res1.Update(signer)
	if err != nil {
		RenderError(w, "Signer status not updated")
		return
	}
	err = res.Update(contract)
	if err != nil {
		RenderError(w, "Signer status not updated")
		return
	}

	json.NewEncoder(w).Encode(contract)
	return

}

func (d *dbServer) SignContract(w http.ResponseWriter, r *http.Request) {

	tokenstring := r.Header["Token"][0]
	claims, cBool := GetClaims(tokenstring)
	if !cBool {
		RenderError(w, "Invalid user request")
		Logger("Invalid user request")
		return
	}
	uID := claims["userid"].(string)

	var sc SignContract
	var contract Contract
	var signer Signer
	var signers []Signer

	_ = json.NewDecoder(r.Body).Decode(&sc)

	Collection := d.sess.Collection(ContractCollection)
	Signercollection := d.sess.Collection(SignerCollection)

	res := Collection.Find(db.Cond{"ContractID": sc.ContractID})
	err := res.One(&contract)

	if err != nil {
		fmt.Println("contract not found")
		return
	}

	// path := Contractfilepath + "/" + contract.ContractID + ".png"
	// Save(sc.FileBase64, path)

	//contract.Filepath = path
	contract.UpdateTime = time.Now().Format(time.RFC850)

	signerCol := d.sess.Collection(SignerCollection)
	res1 := signerCol.Find(db.Cond{"ContractID": contract.ContractID, "userID": uID})
	errstring := res1.One(&signer)

	if errstring != nil {
		RenderError(w, "User Not found")
		return
	}
	signer.SignStatus = "Signed"
	signer.SignDate = time.Now().Format(time.RFC850)

	err = res1.Update(signer)
	if err != nil {
		RenderError(w, "Signer status not updated")
		return
	}

	res2 := Signercollection.Find(db.Cond{"ContractID": contract.ContractID, "CC": 0})
	err = res2.All(&signers)

	if err != nil {
		fmt.Println("Can not find signers")
		return
	}
	res.Update(contract)

	count := 0

	for i := 0; i < len(signers); i++ {
		if signers[i].SignStatus == "Signed" {
			count++
		}
	}
	if count == len(signers) {
		contract.Status = "Completed"
		res.Update(contract)
	}

	json.NewEncoder(w).Encode(signers)
	return

}

// err = res1.One(signer)
// if err != nil {
// 	RenderError(w, "Internal error")
// 	return
// }
// fmt.Println(signer)

// signer.SignStatus = "Signed"
// signer.SignDate = time.Now().Format(time.RFC850)

// err = res1.Update(signer)

// if err != nil {
// 	RenderError(w, "Signer status not updated")
// 	return
// }

// res2 := signerCol.Find(db.Cond{"ContractID": contract.ContractID})
// err = res2.All(&signers)

// if err != nil {
// 	RenderError(w, "Other Signer not found")
// 	return
// }
// count := 0

// for i := 0; i < len(signers); i++ {
// 	if signers[i].SignStatus == "Signed" {
// 		count = count + 1
// 	}
// }
// if count == len(signers) {
// 	contract.Status = "Completed"
// 	res.Update(contract)
// }

// RenderResponse(w, "Contract Signed Successfully", http.StatusOK)
// return

func Save(data string, path string) {
	idx := strings.Index(data, ";base64,")
	if idx < 0 {
		panic("InvalidImage")
	}
	ImageType := data[11:idx]

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
