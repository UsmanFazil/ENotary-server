package DB

import (
	"ENOTARY-Server/Email"
	"ENOTARY-Server/Hashing"
	"encoding/json"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"

	uuid "github.com/satori/go.uuid"
	db "upper.io/db.v3"
)

// TODO new contract creation process here. with file upload
func (d *dbServer) NewContract(w http.ResponseWriter, r *http.Request) {
	var returndata ContractBasic
	r.Body = http.MaxBytesReader(w, r.Body, MaxContractSize)
	err := r.ParseMultipartForm(5000)

	tokenstring := r.Header["Token"][0]
	claims, cBool := GetClaims(tokenstring)
	if !cBool {
		RenderError(w, "Invalid user request")
		Logger("Invalid user request")
		return
	}
	uID := claims["userid"].(string)

	if err != nil {
		RenderError(w, "FILE SHOULD BE LESS THAN 10 MB")
		Logger("CONTRACT FILE SHOULD BE LESS THAN 10 MB")
		return
	}

	// f, header, err := r.FormFile("contractFile")
	f, _, err := r.FormFile("contractFile")
	contractName := r.FormValue("contractName")
	if err != nil {
		RenderError(w, "INVALID_FILE")
		Logger("INVALID CONTRACT FILE")
		return
	}
	defer f.Close()

	// upFileName := header.Filename

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		RenderError(w, "INVALID_FILE")
		Logger("INVALID CONTRACT FILE")
		return
	}
	filetype := http.DetectContentType(bs)
	if filetype != "image/jpeg" && filetype != "image/jpg" &&
		filetype != "image/png" {
		RenderError(w, "INVALID_FILE_TYPE_UPLOAD jpeg,jpg,png OR pdf")
		Logger("INVALID CONTRACT FILE")
		return
	}
	contractID, err := uuid.NewV4()
	if err != nil {
		RenderError(w, "CAN_NOT_GENERATE_CONTRACT_ID")
		Logger("UUID ERROR")
		return
	}

	filepathName := contractID.String()
	fileEndings, err := mime.ExtensionsByType(filetype)
	if err != nil {
		RenderError(w, "INVALID_FILE")
		Logger("INVALID CONTRACT FILE")
		return
	}

	s := string(bs)
	newpath := filepath.Join(Contractfilepath, filepathName+fileEndings[0])
	file, err := os.Create(newpath)

	if err != nil {
		RenderError(w, "INVALID_FILE ")
		Logger("INVALID CONTRACT FILE")
		return
	}
	defer file.Close()
	file.WriteString(s)

	cid := d.ContractInDB(contractName, filepathName, uID, newpath)
	if !cid {
		RenderError(w, "CAN NOT ADD CONTRACT FILE TRY AGAIN")
		Logger("CAN NOT SAVE CONTRACT FILE ON SERVER")
		return
	}
	returndata.ContractID = filepathName
	returndata.Path = newpath
	json.NewEncoder(w).Encode(returndata)
	Logger("NEW CONTRACT ADDED" + filepathName)
	return

}

func (d *dbServer) ContractInDB(cName string, cID string, userid string, filepath string) bool {
	var contract Contract
	contract.ContractID = cID
	contract.Creator = userid
	contract.Filepath = filepath
	contract.Status = "DRAFT"
	contract.ContractcreationTime = time.Now().Format(time.RFC850)
	contract.DelStatus = 0
	contract.Blockchain = 0
	contract.ContractName = cName
	contract.ExpirationTime = time.Now().AddDate(0, 0, 60).Format(time.RFC850)

	Collection := d.sess.Collection(ContractCollection)
	_, err := Collection.Insert(contract)
	if err != nil {
		return false
	}
	return true
}

func (d *dbServer) AddRecipients(w http.ResponseWriter, r *http.Request) {

	var input []Signerdata
	var signer Signer

	_ = json.NewDecoder(r.Body).Decode(&input)
	signerCollection := d.sess.Collection(SignerCollection)

	var result = make([]Signer, len(input))

	for i := 0; i < len(input); i++ {
		user, _, err := d.GetUser(input[i].Email)
		if err != nil {
			RenderError(w, "INVALID RECIPIENT! DOES NOT EXIST ON PLATFORM")
			Logger("INVALID RECIPIENT")
			return
		}
		signer.UserID = user.Userid
		signer.ContractID = input[i].ContractID
		signer.Name = input[i].Name
		signer.SignStatus = "pending"
		signer.Access = 0
		signer.DeleteApprove = 0
		signer.Message = ""
		signer.SignDate = ""

		if input[i].ReceiveCopy {
			signer.CC = 1
		} else {
			signer.CC = 0
		}
		_, errstring := signerCollection.Insert(signer)

		if errstring != nil {
			RenderError(w, "User not exists on platform")
			Logger("cannot add signer")
			return
		}
		result[i] = signer
	}

	json.NewEncoder(w).Encode(result)
	return

}

// var input []Signerdata
// var signer Signer

// _ = json.NewDecoder(r.Body).Decode(&input)
// Collection := d.sess.Collection(SignerCollection)

// for _, s := range input {
// 	user, _, err := d.GetUser(s.Email)
// 	if err != nil {
// 		RenderError(w, "INVALID RECIPIENT")
// 		Logger("INVALID RECIPIENT")
// 		return
// 	}

// 	signer.ContractID = s.ContractID
// 	signer.UserID = user.Userid
// 	signer.Name = s.Name
// 	signer.Access = 0
// 	signer.SignStatus = "Not Signed"
// 	signer.DeleteApprove = 0
// 	if s.ReceiveCopy == true {
// 		signer.CC = 1
// 	} else {
// 		signer.CC = 0
// 	}

// 	_, errstring := Collection.Insert(signer)

// 	if errstring != nil {
// 		RenderError(w, "User does not exist on the platform")
// 		Logger("DB INSERTION ERROR AT RECIPIENTS",)
// 		return
// 	}
// }

// RenderResponse(w, "SIGNERS ADDED", http.StatusOK)
// Logger("RECIPIENTS ADDED" + signer.ContractID)
// return

func (d *dbServer) WaitingforOther(userid string) (uint64, error) {
	Collection := d.sess.Collection(ContractCollection)
	res := Collection.Find(db.Cond{"Creator": userid, "delStatus": 0, "status": "In Progress"})
	total, err := res.Count()

	if err != nil {
		return 0, err
	}
	return total, nil
}

func (d *dbServer) WaitingforMe(userid string) (uint64, error) {
	Collection := d.sess.Collection(SignerCollection)
	res := Collection.Find(db.Cond{"userID": userid, "Access": 1, "SignStatus": "needs to sign", "DeleteApprove": 0})
	total, err := res.Count()

	if err != nil {
		return 0, err
	}
	return total, nil
}

func (d *dbServer) ExpiringSoon(userid string) (int, error, []Contract) {
	resbool, inboxList := d.InboxContractsList(userid, false)
	res2bool, sentList := d.SentContractsList(userid, false, false)
	var expContractsList []Contract
	count := 0
	if !resbool {
		return 0, nil, nil
	}
	if !res2bool {
		return 0, nil, nil
	}

	contractList := inboxList
	contractList = append(contractList, sentList...)

	for _, index := range contractList {
		if index.Status != "Completed" {
			t, _ := time.Parse(RFC850, index.ExpirationTime)
			timeNow := time.Now()
			diff := t.Sub(timeNow)
			exptime := timeNow.AddDate(0, 0, 7)
			diffexp := exptime.Sub(timeNow)

			if diff > 0 && diff < diffexp {
				expContractsList = append(expContractsList, index)
				count++
			}
		}
	}
	return count, nil, expContractsList

}

func TimeSearch(Allcontracts []Contract, timeframe string) []Contract {

	oneyear := time.Now().AddDate(-1, 0, 0)
	sixmonths := time.Now().AddDate(0, -6, 0)
	onemonth := time.Now().AddDate(0, -1, 0)
	oneweek := time.Now().AddDate(0, 0, -7)
	oneday := time.Now().AddDate(0, 0, -1)
	var searchList []Contract

	if timeframe == "Last one year" {
		for _, index := range Allcontracts {
			t, _ := time.Parse(RFC850, index.ContractcreationTime)
			diff := t.Sub(oneyear)
			if diff > 0 {
				searchList = append(searchList, index)
			}
		}
	}
	if timeframe == "Last six months" {
		for _, index := range Allcontracts {
			t, _ := time.Parse(RFC850, index.ContractcreationTime)
			diff := t.Sub(sixmonths)
			if diff > 0 {
				searchList = append(searchList, index)
			}
		}
	}
	if timeframe == "Last one month" {
		for _, index := range Allcontracts {
			t, _ := time.Parse(RFC850, index.ContractcreationTime)
			diff := t.Sub(onemonth)
			if diff > 0 {
				searchList = append(searchList, index)
			}
		}
	}
	if timeframe == "Last one week" {
		for _, index := range Allcontracts {
			t, _ := time.Parse(RFC850, index.ContractcreationTime)
			diff := t.Sub(oneweek)
			if diff > 0 {
				searchList = append(searchList, index)
			}
		}
	}
	if timeframe == "Last one day" {
		for _, index := range Allcontracts {
			t, _ := time.Parse(RFC850, index.ContractcreationTime)
			diff := t.Sub(oneday)
			if diff > 0 {
				searchList = append(searchList, index)
			}
		}
	}
	return searchList
}

func (d *dbServer) SearchAlgo(w http.ResponseWriter, r *http.Request) {
	var searchInput SearchInput
	var searchList []Contract

	_ = json.NewDecoder(r.Body).Decode(&searchInput)

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
		RenderResponse(w, "NO CONTRACT FOUND FOR THE USER", http.StatusOK)
		Logger("No contracts found in search")
		return
	}
	Allcontracts := append(inboxList, sentList...)
	Allcontracts = append(Allcontracts, draftlist...)

	if searchInput.Status == "All" {
		if searchInput.Date == "All" {
			for _, index := range Allcontracts {
				if index.ContractName == searchInput.ContractName {
					searchList = append(searchList, index)
				}
			}
		}
		if searchInput.ContractName == "" {
			searchList = Allcontracts

		} else {
			tmp := TimeSearch(Allcontracts, searchInput.Date)
			for _, index := range tmp {
				if index.ContractName == searchInput.ContractName {
					searchList = append(searchList, index)
				}
			}
		}
	}

	if searchInput.ContractName == "" {
		if searchInput.Date == "All" {
			for _, index := range Allcontracts {
				if index.Status == searchInput.Status {
					searchList = append(searchList, index)
				}
			}

		} else {
			tmp := TimeSearch(Allcontracts, searchInput.Date)
			for _, index := range tmp {
				if index.Status == searchInput.Status {
					searchList = append(searchList, index)
				}
			}
		}

	}

	if searchInput.Date == "All" && searchInput.ContractName != "" && searchInput.Status != "All" {
		for _, index := range Allcontracts {
			if searchInput.ContractName == index.ContractName {
				if searchInput.Status == index.Status {
					searchList = append(searchList, index)
				}
			}
		}
	}

	if searchInput.ContractName != "" && searchInput.Date != "All" && searchInput.Status != "All" {
		tmp := TimeSearch(Allcontracts, searchInput.Date)

		for _, index := range tmp {
			if searchInput.ContractName == index.ContractName {
				if searchInput.Status == index.Status {
					searchList = append(searchList, index)
				}
			}
		}
	}
	json.NewEncoder(w).Encode(&searchList)
	Logger("Contract Search successful | userID: " + uID)
	return
}

func (d *dbServer) DeleteDraft(w http.ResponseWriter, r *http.Request) {
	tokenstring := r.Header["Token"][0]
	claims, cBool := GetClaims(tokenstring)
	if !cBool {
		RenderError(w, "Invalid user request")
		Logger("Invalid user request")
		return
	}
	uID := claims["userid"].(string)

	var contract Contract
	_ = json.NewDecoder(r.Body).Decode(&contract)

	Collection := d.sess.Collection(ContractCollection)
	res := Collection.Find(db.Cond{"ContractID": contract.ContractID, "Creator": uID, "status": "DRAFT"})

	total, _ := res.Count()

	if total != 1 {
		RenderError(w, "Can not delete contract, TRY AGAIN")
		return
	}
	err := res.Delete()
	if err != nil {
		RenderError(w, "Can not delete contract, TRY AGAIN")
		return
	}

	RenderResponse(w, "CONTRACT DELETED", http.StatusOK)
	return

}

func (d *dbServer) ContractDetails(w http.ResponseWriter, r *http.Request) {
	var contract Contract
	var tmp Contract
	var signers []Signer
	var CD ContractDetail
	_ = json.NewDecoder(r.Body).Decode(&contract)

	Collection := d.sess.Collection(ContractCollection)
	Signercollection := d.sess.Collection(SignerCollection)
	res := Collection.Find(db.Cond{"ContractID": contract.ContractID})
	err := res.One(&tmp)

	if err != nil {
		json.NewEncoder(w).Encode(CD)
		Logger("cannot find contract " + contract.ContractID)
		return
	}

	res1 := Signercollection.Find(db.Cond{"ContractID": tmp.ContractID})
	errstring := res1.All(&signers)

	if errstring != nil {
		json.NewEncoder(w).Encode(CD)
		Logger("cannot find Signers " + contract.ContractID)
		return
	}

	CD.ContractData = tmp
	CD.Signers = signers

	json.NewEncoder(w).Encode(CD)
	return

}

func (d *dbServer) ContractHashDetails(w http.ResponseWriter, r *http.Request) {
	var contract Contract
	var tmp Contract
	var signers []Signer
	var CH ContractDetailHash
	_ = json.NewDecoder(r.Body).Decode(&contract)

	Collection := d.sess.Collection(ContractCollection)
	Signercollection := d.sess.Collection(SignerCollection)
	res := Collection.Find(db.Cond{"ContractID": contract.ContractID})
	err := res.One(&tmp)

	if err != nil {
		json.NewEncoder(w).Encode(CH)
		Logger("cannot find contract " + contract.ContractID)
		return
	}

	res1 := Signercollection.Find(db.Cond{"ContractID": tmp.ContractID})
	errstring := res1.All(&signers)

	if errstring != nil {
		json.NewEncoder(w).Encode(CH)
		Logger("cannot find Signers " + contract.ContractID)
		return
	}

	CH.ContractData = tmp
	CH.Signers = signers
	CH.ContractHash = Hashing.FindHash(tmp.Filepath)

	json.NewEncoder(w).Encode(CH)
	return

}
func (d *dbServer) UpdateBlockchainstatus(w http.ResponseWriter, r *http.Request) {

	var swi SaveWalletinput
	var contract Contract
	var wallet WalletInfo
	var signers []Signer
	var user User

	_ = json.NewDecoder(r.Body).Decode(&swi)
	Collection := d.sess.Collection(ContractCollection)
	signerCollection := d.sess.Collection(SignerCollection)
	userCollection := d.sess.Collection(UserCollection)
	//WalletCollection := d.sess.Collection(WalletsCollection)

	res := Collection.Find(db.Cond{"ContractID": swi.ContractID})
	err := res.One(&contract)
	if err != nil {
		RenderError(w, "CAN NOT UPDATE CONTRACT STATUS, Please contact at enotary99@gmail.com")
		Logger("CAN NOT UPDATE CONTRACT STATUS " + swi.ContractID)
		return
	}
	res.Update(map[string]int{
		"Blockchain": 1,
	})

	wallet.Userid = swi.UserID
	wallet.PublicAddress = swi.PublicAddress

	q := d.sess.InsertInto("Wallets").Columns("userid", "walletaddress").Values(wallet.Userid, wallet.PublicAddress)
	_, _ = q.Exec()

	res1 := signerCollection.Find(db.Cond{"ContractID": swi.ContractID})
	_ = res1.All(&signers)

	for i := 0; i < len(signers); i++ {
		res := userCollection.Find(db.Cond{"userid": signers[i].UserID})
		_ = res.One(&user)

		go Email.BlockchainEmail(user.Email, "CONTRACT SAVED IN BLOCKCHAIN", "YOUR CONTRACT "+swi.ContractID+" is saved in blockchain.")
	}

	json.NewEncoder(w).Encode(contract)
	Logger("CONTRACT SAVED IN BLOCKCHAIN " + swi.ContractID)
	return
}

// func (d *dbServer) SearchContract(w http.ResponseWriter, r *http.Request) {
// 	var searchInput SearchInput
// 	var contracts []Contract

// 	_ = json.NewDecoder(r.Body).Decode(&searchInput)
// 	Collection := d.sess.Collection(ContractCollection)

// 	if searchInput.Status == "All" && searchInput.Date == "All" {
// 		res := Collection.Find(db.Cond{"contractName": searchInput.ContractName})
// 		total, _ := res.Count()
// 		if total < 1 {
// 			RenderResponse(w, "CAN NOT FIND ANY CONTRACT", http.StatusOK)
// 			return
// 		}
// 		err := res.All(&contracts)
// 		if err != nil {
// 			RenderResponse(w, "CAN NOT FIND ANY CONTRACT", http.StatusOK)
// 			return
// 		}
// 		json.NewEncoder(w).Encode(contracts)
// 		return
// 	}

// 	if searchInput.ContractName == "" && searchInput.Date == "All" {
// 		res1 := Collection.Find(db.Cond{"status": searchInput.Status})
// 		total, _ := res1.Count()
// 		if total < 1 {
// 			RenderResponse(w, "CAN NOT FIND ANY CONTRACT", http.StatusOK)
// 			return
// 		}
// 		err := res1.All(&contracts)
// 		if err != nil {
// 			RenderResponse(w, "CAN NOT FIND ANY CONTRACT", http.StatusOK)
// 			return
// 		}
// 		json.NewEncoder(w).Encode(contracts)
// 		return
// 	}

// 	res4 := Collection.Find(db.Cond{"contractName": searchInput.ContractName})
// 	err := res4.All(&contracts)

// 	if err != nil {
// 		RenderError(w, "invalid")
// 		return
// 	}

// to do: date search

// t, _ := time.Parse(RFC850, contracts[0].ContractcreationTime)
// a := time.Now().AddDate(-1, 0, 0)

// diff := t.Sub(a)
// fmt.Printf("Lifespan is %+v", diff)
// if diff > 0 {
// 	fmt.Println("yes print it")
// } else {
// 	fmt.Println("nop")
// }

// }
