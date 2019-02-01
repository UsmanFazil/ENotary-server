package DB

import (
	"encoding/json"
	"net/http"

	db "upper.io/db.v3"
)

func (d *dbServer) InboxData(w http.ResponseWriter, r *http.Request) {
	var tmp User
	var signer []Signer
	var tmpContract Contract
	i := 0
	_ = json.NewDecoder(r.Body).Decode(&tmp)

	signercollection := d.sess.Collection(SignerCollection)
	contractCollection := d.sess.Collection(ContractCollection)

	res := signercollection.Find(db.Cond{"userID": tmp.Userid, "Access": 1})
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
