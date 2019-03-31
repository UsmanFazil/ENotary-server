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
		return
	}

	tokenstring := r.Header["Token"][0]
	claims, cBool := GetClaims(tokenstring)
	if !cBool {
		RenderError(w, "Invalid user request")
		return
	}
	userID := claims["userid"].(string)

	folderID, err := uuid.NewV4()
	if err != nil {
		RenderError(w, "CAN_NOT_GENERATE_FOLDER_ID")
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
		return
	}
	RenderResponse(w, "NEW FOLDER CREATED", http.StatusOK)
	return
}

func (d *dbServer) AddContract(w http.ResponseWriter, r *http.Request) {

	var CF ContractFolder
	_ = json.NewDecoder(r.Body).Decode(&CF)

	Collection := d.sess.Collection(ContractFolderCollection)

	_, err := Collection.Insert(CF)

	if err != nil {
		RenderError(w, "CAN NOT ADD CONTRACT IN FOLDER, TRY AGAIN")
		return
	}
	RenderResponse(w, "CONTRACT ADDED SUCCESSFULLY", http.StatusOK)
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
	if total < 1 {
		RenderResponse(w, "NO CONTRACTS IN THIS FOLDER", http.StatusOK)
		return
	}
	err := res.All(&CFs)
	if err != nil {
		RenderError(w, "NO CONTRACTS IN THIS FOLDER")
		return
	}

	var contracts = make([]Contract, total)
	i := 0
	for _, v := range CFs {
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
