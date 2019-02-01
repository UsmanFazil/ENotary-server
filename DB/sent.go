package DB

import (
	"encoding/json"
	"net/http"

	"upper.io/db.v3"
)

func (d dbServer) SentContract(w http.ResponseWriter, r *http.Request) {
	var tmp User
	var contracts []Contract
	_ = json.NewDecoder(r.Body).Decode(&tmp)

	contractCollection := d.sess.Collection(ContractCollection)
	res := contractCollection.Find(db.Cond{"Creator": tmp.Userid})
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
