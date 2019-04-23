package DB

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"os"

	"upper.io/db.v3"
)

func (d *dbServer) ExportCSV(w http.ResponseWriter, r *http.Request) {
	var contract Contract

	_ = json.NewDecoder(r.Body).Decode(&contract)

	Collection := d.sess.Collection(ContractCollection)
	res := Collection.Find(db.Cond{"ContractID": contract.ContractID})

	err := res.One(&contract)
	if err != nil {
		RenderError(w, "Error")
	}

	filepath := CSVpath + "/" + contract.ContractID + ".csv"
	file, err := os.Create(filepath)
	defer file.Close()

	var data = [][]string{{"ContractID", contract.ContractID}, {"Name", contract.ContractName}, {"Status", contract.Status}, {"Created on", contract.ContractcreationTime}, {"LAST UPDATE", contract.UpdateTime}, {"CONTRACT OWNER", contract.Creator}, {"File exported from", "E-NOTARY"}}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			RenderError(w, "Error")
		}
	}
	json.NewEncoder(w).Encode(filepath)
	Logger("NEW CSV GENERATED " + contract.ContractID)
	return
}
