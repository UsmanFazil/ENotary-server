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
	resbool, contracts := d.InboxContractsList(userID, false)

	if !resbool {
		json.NewEncoder(w).Encode(contracts)
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
	resbool, contracts := d.SentContractsList(userID, false, false)
	if !resbool {
		json.NewEncoder(w).Encode(contracts)
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
	resbool, contracts := d.SentContractsList(userID, true, false)
	if !resbool {
		json.NewEncoder(w).Encode(contracts)
		Logger("CAN NOT FIND ANY CONTRACT " + userID)
		return
	}
	json.NewEncoder(w).Encode(contracts)

}

func (d *dbServer) ExpiringsoonContracts(w http.ResponseWriter, r *http.Request) {
	tokenstring := r.Header["Token"][0]
	claims, cBool := GetClaims(tokenstring)
	if !cBool {
		RenderError(w, "Invalid user request")
		Logger("Invalid user request")
		return
	}
	userID := claims["userid"].(string)

	_, errbool, contractsList := d.ExpiringSoon(userID)

	if errbool != nil {
		json.NewEncoder(w).Encode(contractsList)
		return
	}
	json.NewEncoder(w).Encode(contractsList)
	return
}

func (d *dbServer) ActionRequired(w http.ResponseWriter, r *http.Request) {
	tokenstring := r.Header["Token"][0]
	claims, cBool := GetClaims(tokenstring)
	if !cBool {
		RenderError(w, "Invalid user request")
		Logger("Invalid user request")
		return
	}
	userID := claims["userid"].(string)

	resbool, contracts := d.InboxContractsList(userID, true)

	if !resbool {
		json.NewEncoder(w).Encode(contracts)
		Logger("CAN NOT FIND ANY CONTRACT " + userID)
		return
	}
	json.NewEncoder(w).Encode(contracts)
	return
}

func (d *dbServer) WaitingForOthers(w http.ResponseWriter, r *http.Request) {
	tokenstring := r.Header["Token"][0]
	claims, cBool := GetClaims(tokenstring)
	if !cBool {
		RenderError(w, "Invalid user request")
		Logger("Invalid user request")
		return
	}
	userID := claims["userid"].(string)

	resbool, contracts := d.SentContractsList(userID, false, true)
	if !resbool {
		json.NewEncoder(w).Encode(contracts)
		Logger("CAN NOT FIND ANY CONTRACT " + userID)
		return
	}
	json.NewEncoder(w).Encode(contracts)
	return
}
func (d *dbServer) Completed(w http.ResponseWriter, r *http.Request) {
	var contracts []Contract
	tokenstring := r.Header["Token"][0]
	claims, cBool := GetClaims(tokenstring)
	if !cBool {
		RenderError(w, "Invalid user request")
		Logger("Invalid user request")
		return
	}
	uID := claims["userid"].(string)

	resbool1, inboxList := d.InboxContractsList(uID, false)
	resbool2, sentList := d.SentContractsList(uID, false, false)
	resbool3, draftlist := d.SentContractsList(uID, true, false)

	if !resbool1 && !resbool2 && !resbool3 {
		json.NewEncoder(w).Encode(&contracts)
		Logger("No contracts found in search")
		return
	}
	Allcontracts := append(inboxList, sentList...)
	Allcontracts = append(Allcontracts, draftlist...)

	for _, index := range Allcontracts {
		if index.Status == "Completed" {
			contracts = append(contracts, index)
		}
	}
	json.NewEncoder(w).Encode(&contracts)
}

func (d *dbServer) SentContractsList(userid string, drafts bool, completed bool) (bool, []Contract) {
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
		if !completed {
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

		} else {
			res := contractCollection.Find(db.Cond{"Creator": userid, "status": "In Progress"})
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
}

func (d *dbServer) InboxContractsList(userid string, unsigned bool) (bool, []Contract) {
	var signer []Signer
	var tmpContract Contract
	i := 0
	var total uint64

	signercollection := d.sess.Collection(SignerCollection)
	contractCollection := d.sess.Collection(ContractCollection)

	if unsigned {
		res := signercollection.Find(db.Cond{"userID": userid, "Access": 1, "SignStatus": "pending"})
		total, _ = res.Count()
		if total < 1 {
			return false, nil
		}

		err := res.All(&signer)
		if err != nil {
			return false, nil
		}
	} else {
		res := signercollection.Find(db.Cond{"userID": userid, "Access": 1})
		total, _ = res.Count()
		if total < 1 {
			return false, nil
		}

		err := res.All(&signer)
		if err != nil {
			return false, nil
		}
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

func (d *dbServer) Manage(w http.ResponseWriter, r *http.Request) {

	var folderlist []Folder
	tokenstring := r.Header["Token"][0]
	claims, cBool := GetClaims(tokenstring)
	if !cBool {
		RenderError(w, "Invalid user request")
		Logger("Invalid user request")
		return
	}
	uID := claims["userid"].(string)
	_, contracts := d.InboxContractsList(uID, false)

	// if !resbool {
	// 	json.NewEncoder(w).Encode(contracts)
	// 	Logger("CAN NOT FIND ANY CONTRACT " + uID)
	// 	return
	// }

	Collection := d.sess.Collection(FolderCollection)
	res := Collection.Find(db.Cond{"userID": uID})
	_ = res.All(&folderlist)

	var manageobj ManageScreen

	manageobj.InboxContracts = contracts
	manageobj.FolderList = folderlist

	json.NewEncoder(w).Encode(manageobj)
	return

}
