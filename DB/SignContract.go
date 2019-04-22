package DB

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// func (d *dbServer) SignContract(w http.ResponseWriter, r *http.Request) {

// 	var sc SignContract
// 	var signer Signer
// 	var signers []Signer
// 	_ = json.NewDecoder(r.Body).Decode(&sc)

// 	tokenstring := r.Header["Token"][0]
// 	claims, cBool := GetClaims(tokenstring)
// 	if !cBool {
// 		RenderError(w, "Invalid user request")
// 		Logger("Invalid user request")
// 		return
// 	}
// 	uID := claims["userid"].(string)

// 	contractCollection := d.sess.Collection(ContractCollection)
// 	signerCollection := d.sess.Collection(SignerCollection)
// 	userCollection := d.sess.Collection(UserCollection)

// 	res := contractCollection.Find(db.Cond{"ContractID": sc.ContractID})
// 	total, _ := res.Count()
// 	if total != 1 {
// 		RenderError(w, "CONTRACT NOT FOUND")
// 		return
// 	}

// 	res1 := signerCollection.Find(db.Cond{"ContractID": sc.ContractID, "userID": uID})
// 	total1, _ := res1.Count()
// 	if total1 != 1 {
// 		RenderError(w, "CONTRACT NOT FOUND")
// 		return
// 	}
// 	res1.One(&signer)

// 	if signer.SignStatus != "pending" {
// 		RenderError(w, "Contract Already signed by User")
// 		return
// 	}
// 	signer.SignStatus = "Signed"
// 	signer.SignDate = time.Now().Format(time.RFC850)
// 	res.Update(signer)

// 	res2 := signerCollection.Find(db.Cond{"ContractID": sc.ContractID})
// 	res2.All(&signer)

// 	signedusers := 0
// 	for i := 0; i< len(signers); i ++ {
// 		if signers[i].SignStatus
// 	}

// }

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
