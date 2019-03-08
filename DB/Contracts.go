package DB

import (
	"encoding/json"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
	db "upper.io/db.v3"
)

// TODO new contract creation process here. with file upload
func (d *dbServer) NewContract(w http.ResponseWriter, r *http.Request) {

	r.Body = http.MaxBytesReader(w, r.Body, MaxContractSize)
	err := r.ParseMultipartForm(5000)

	tokenstring := r.Header["Token"][0]
	claims, cBool := GetClaims(tokenstring)
	if !cBool {
		RenderError(w, "Invalid user request")
		return
	}
	uID := claims["userid"].(string)

	if err != nil {
		RenderError(w, "FILE SHOULD BE LESS THAN 10 MB")
		return
	}

	f, header, err := r.FormFile("contractFile")
	if err != nil {
		RenderError(w, "INVALID_FILE")
		return
	}
	defer f.Close()

	upFileName := header.Filename

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		RenderError(w, "INVALID_FILE")
		return
	}
	filetype := http.DetectContentType(bs)
	if filetype != "image/jpeg" && filetype != "image/jpg" &&
		filetype != "image/png" && filetype != "application/pdf" {
		RenderError(w, "INVALID_FILE_TYPE_UPLOAD jpeg,jpg,png OR pdf")
		return
	}
	contractID, err := uuid.NewV4()
	if err != nil {
		RenderError(w, "CAN_NOT_GENERATE_CONTRACT_ID")
		return
	}

	filepathName := contractID.String()
	fileEndings, err := mime.ExtensionsByType(filetype)
	if err != nil {
		RenderError(w, "INVALID_FILE")
		return
	}

	s := string(bs)
	newpath := filepath.Join(Contractfilepath, filepathName+fileEndings[0])
	file, err := os.Create(newpath)

	if err != nil {
		RenderError(w, "INVALID_FILE ")
		return
	}
	defer file.Close()
	file.WriteString(s)

	cid := d.ContractInDB(upFileName, filepathName, uID, newpath)
	if !cid {
		RenderError(w, "CAN NOT ADD CONTRACT FILE TRY AGAIN")
		return
	}
	RenderResponse(w, filepathName, http.StatusOK)
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

	contract.ExpirationTime = strconv.FormatInt(time.Now().Add(1440*time.Hour).Unix(), 10) // 1440 = 60 days

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
	Collection := d.sess.Collection(SignerCollection)

	for _, s := range input {
		user, _, err := d.GetUser(s.Email)
		if err != nil {
			RenderError(w, "INVALID RECIPIENT")
			return
		}

		signer.ContractID = s.ContractID
		signer.UserID = user.Userid
		signer.Name = s.Name
		signer.Access = 0
		signer.SignStatus = "Not Signed"
		signer.DeleteApprove = 0

		_, err = Collection.Insert(signer)
		if err != nil {
			RenderError(w, "CAN NOT ADD SIGNER TRY AGAIN")
			return
		}
	}
	RenderResponse(w, "SIGNERS ADDED", http.StatusOK)
	return
}

func (d *dbServer) WaitingforOther(userid string) (uint64, error) {
	Collection := d.sess.Collection(ContractCollection)
	res := Collection.Find(db.Cond{"Creator": userid, "delStatus": 0, "status": "in progress"})
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
