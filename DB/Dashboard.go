package DB

import (
	"encoding/json"
	"net/http"

	db "upper.io/db.v3"
)

func (d *dbServer) InboxData(w http.ResponseWriter, r *http.Request) {

	tokenstring := r.Header["Token"][0]
	claims, cBool := GetClaims(tokenstring)
	if !cBool {
		RenderError(w, "Invalid user request")
		Logger("Invalid user request")
		return
	}
	userID := claims["userid"].(string)
	resbool, contracts := d.InboxContractsList(userID)

	if !resbool {
		RenderResponse(w, "CAN NOT FIND CONTRACT FOR THE USER", http.StatusOK)
		Logger("CAN NOT FIND ANY CONTRACT " + userID)
		return
	}
	json.NewEncoder(w).Encode(contracts)
	return

}

func (d *dbServer) SentContract(w http.ResponseWriter, r *http.Request) {

	tokenstring := r.Header["Token"][0]
	claims, cBool := GetClaims(tokenstring)
	if !cBool {
		RenderError(w, "Invalid user request")
		Logger("Invalid user request")
		return
	}
	userID := claims["userid"].(string)
	resbool, contracts := d.SentContractsList(userID, false)
	if !resbool {
		RenderResponse(w, "CAN NOT FIND CONTRACT FOR THE USER", http.StatusOK)
		Logger("CAN NOT FIND ANY CONTRACT " + userID)
		return
	}

	json.NewEncoder(w).Encode(contracts)

}

func (d *dbServer) DraftContracts(w http.ResponseWriter, r *http.Request) {
	tokenstring := r.Header["Token"][0]
	claims, cBool := GetClaims(tokenstring)
	if !cBool {
		RenderError(w, "Invalid user request")
		Logger("Invalid user request")
		return
	}
	userID := claims["userid"].(string)
	resbool, contracts := d.SentContractsList(userID, true)
	if !resbool {
		RenderResponse(w, "CAN NOT FIND CONTRACT FOR THE USER", http.StatusOK)
		Logger("CAN NOT FIND ANY CONTRACT " + userID)
		return
	}
	json.NewEncoder(w).Encode(contracts)

}

func (d *dbServer) SentContractsList(userid string, drafts bool) (bool, []Contract) {
	var contracts []Contract
	contractCollection := d.sess.Collection(ContractCollection)

	if drafts {
		res := contractCollection.Find(db.Cond{"Creator": userid, "status": "DRAFT"})
		total, _ := res.Count()
		if total < 1 {
			return false, nil
		}
		err := res.All(&contracts)
		if err != nil {
			return false, nil
		}
		return true, contracts

	} else {
		res := contractCollection.Find(db.Cond{"Creator": userid, "status !=": "DRAFT"})
		total, _ := res.Count()
		if total < 1 {
			return false, nil
		}
		err := res.All(&contracts)
		if err != nil {
			return false, nil
		}
		return true, contracts
	}
}

func (d *dbServer) InboxContractsList(userid string) (bool, []Contract) {
	var signer []Signer
	var tmpContract Contract
	i := 0

	signercollection := d.sess.Collection(SignerCollection)
	contractCollection := d.sess.Collection(ContractCollection)

	res := signercollection.Find(db.Cond{"userID": userid, "Access": 1})
	total, _ := res.Count()
	if total < 1 {
		return false, nil
	}

	err := res.All(&signer)
	if err != nil {
		return false, nil
	}

	var contracts = make([]Contract, total)
	for _, v := range signer {
		res1 := contractCollection.Find(db.Cond{"ContractID": v.ContractID})
		err := res1.One(&tmpContract)
		if err != nil {
			return false, nil
		}
		contracts[i] = tmpContract
		i++
	}

	return true, contracts

}
