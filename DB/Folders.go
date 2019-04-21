package DB

import (
	"encoding/json"
	"net/http"
	"strings"

	uuid "github.com/satori/go.uuid"
	db "upper.io/db.v3"
)

func (d *dbServer) NewFolder(w http.ResponseWriter, r *http.Request) {

	var newFolder Folder
	_ = json.NewDecoder(r.Body).Decode(&newFolder)

	folderName := strings.TrimSpace(newFolder.FolderName)

	if len(folderName) < 1 {
		RenderError(w, "INVALID NAME")
		Logger("INVALID FOLDER NAME")
		return
	}

	tokenstring := r.Header["Token"][0]
	claims, cBool := GetClaims(tokenstring)
	if !cBool {
		RenderError(w, "Invalid user request")
		Logger("Invalid user request")
		return
	}
	userID := claims["userid"].(string)

	folderID, err := uuid.NewV4()
	if err != nil {
		RenderError(w, "CAN_NOT_GENERATE_FOLDER_ID")
		Logger("UUID ERROR")
		return
	}

	newFolder.FolderID = folderID.String()
	newFolder.UserID = userID
	newFolder.ParentFolder = "non"
	newFolder.FolderType = "Contract"
	newFolder.FolderName = folderName

	Collection := d.sess.Collection(FolderCollection)
	_, err = Collection.Insert(newFolder)

	if err != nil {
		RenderError(w, "Can not create new folder, try again")
		Logger("DB ERROR (FOLDER CREATION)")
		return
	}
	json.NewEncoder(w).Encode(newFolder)
	Logger("NEW FOLDER CREATED " + newFolder.FolderID)
	return
}

func (d *dbServer) AddContract(w http.ResponseWriter, r *http.Request) {

	var CF ContractFolder
	_ = json.NewDecoder(r.Body).Decode(&CF)

	Collection := d.sess.Collection(ContractFolderCollection)

	_, err := Collection.Insert(CF)

	if err != nil {
		RenderError(w, "CAN NOT ADD CONTRACT IN FOLDER, TRY AGAIN")
		Logger("CAN NOT ADD CONTRACT IN FOLDER, TRY AGAIN | ContractID :" + CF.ContractID)
		return
	}
	RenderResponse(w, "CONTRACT ADDED SUCCESSFULLY", http.StatusOK)
	Logger("CONTRACT ADDED TO FOLDER | Folderid: " + CF.FolderID + " Contractdid: " + CF.ContractID)
	return
}

func (d *dbServer) FolderContractList(w http.ResponseWriter, r *http.Request) {
	var folder Folder
	var CFs []ContractFolder
	var tmpContract Contract

	_ = json.NewDecoder(r.Body).Decode(&folder)

	cfCollection := d.sess.Collection(ContractFolderCollection)
	contractCollection := d.sess.Collection(ContractCollection)

	res := cfCollection.Find(db.Cond{"folderID": folder.FolderID})
	total, _ := res.Count()
	var contracts = make([]Contract, total)
	if total < 1 {
		json.NewEncoder(w).Encode(contracts)
		Logger("NO CONTRACTS FOUND | FolderID :" + folder.FolderID)
		return
	}
	err := res.All(&CFs)
	if err != nil {
		json.NewEncoder(w).Encode(contracts)
		Logger("NO CONTRACTS FOUND | FolderID :" + folder.FolderID)
		return
	}

	i := 0
	for _, v := range CFs {
		res1 := contractCollection.Find(db.Cond{"ContractID": v.ContractID})
		err := res1.One(&tmpContract)
		if err != nil {
			json.NewEncoder(w).Encode(contracts)
			Logger("NO CONTRACTS FOUND | FolderID :" + folder.FolderID)
			return
		}
		contracts[i] = tmpContract
		i++
	}
	json.NewEncoder(w).Encode(contracts)
	return
}
