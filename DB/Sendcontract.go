package DB

import (
	"ENOTARY-Server/Email"
	"encoding/json"
	"net/http"

	"upper.io/db.v3"
)

func (d *dbServer) SendContract(w http.ResponseWriter, r *http.Request) {
	var contractinfo SendContract
	var signers []Signer
	var user User

	_ = json.NewDecoder(r.Body).Decode(&contractinfo)

	signerCollection := d.sess.Collection(SignerCollection)
	userCollection := d.sess.Collection(UserCollection)
	res := signerCollection.Find(db.Cond{"ContractID": contractinfo.ContractID})

	err := res.All(&signers)
	if err != nil {
		RenderResponse(w, "CONTRACT HAS NO RECEPIENTS", http.StatusOK)
		return
	}
	res.Update(map[string]int{
		"Access": 1,
	})

	for i := 0; i < len(signers); i++ {
		res := userCollection.Find(db.Cond{"userid": signers[i].UserID})
		_ = res.One(&user)

		go Email.ContractEmail(user.Email, contractinfo.EmailSubj, contractinfo.EmailMsg)

	}
	RenderResponse(w, "EMAIL SENT TO ALL RECEPIENTS", http.StatusOK)
	return

}
