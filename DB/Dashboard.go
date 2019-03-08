package DB

import (
	"encoding/json"
	"net/http"

	db "upper.io/db.v3"
)

var MySigningKey = []byte("secretkey")

func (d *dbServer) InboxData(w http.ResponseWriter, r *http.Request) {
	var signer []Signer
	var tmpContract Contract
	i := 0

	signercollection := d.sess.Collection(SignerCollection)
	contractCollection := d.sess.Collection(ContractCollection)

	tokenstring := r.Header["Token"][0]
	claims, cBool := GetClaims(tokenstring)
	if !cBool {
		RenderError(w, "Invalid user request")
		return
	}
	uID := claims["userid"]

	res := signercollection.Find(db.Cond{"userID": uID, "Access": 1})
	total, _ := res.Count()
	if total < 1 {
		RenderResponse(w, "CAN NOT FIND ANY CONTRACT FOR THE USER", http.StatusOK)
		return
	}

	err := res.All(&signer)
	if err != nil {
		RenderError(w, "CAN NOT FIND ANY CONTRACT FOR THE USER")
		return
	}

	var contracts = make([]Contract, total)
	for _, v := range signer {
		res1 := contractCollection.Find(db.Cond{"ContractID": v.ContractID})
		err := res1.One(&tmpContract)
		if err != nil {
			RenderError(w, "CAN NOT FIND ANY CONTRACT FOR THE USER")
			return
		}
		contracts[i] = tmpContract
		i++
	}
	json.NewEncoder(w).Encode(contracts)
	return

}

func (d dbServer) SentContract(w http.ResponseWriter, r *http.Request) {
	var contracts []Contract

	tokenstring := r.Header["Token"][0]
	claims, cBool := GetClaims(tokenstring)
	if !cBool {
		RenderError(w, "Invalid user request")
		return
	}
	uID := claims["userid"]

	contractCollection := d.sess.Collection(ContractCollection)
	res := contractCollection.Find(db.Cond{"Creator": uID})
	total, _ := res.Count()
	if total < 1 {
		RenderResponse(w, "CAN NOT FIND ANY CONTRACT FOR THE USER", http.StatusOK)
		return
	}
	err := res.All(&contracts)
	if err != nil {
		RenderError(w, "CAN NOT FIND ANY CONTRACT FOR THE USER")
		return
	}
	json.NewEncoder(w).Encode(contracts)

}
