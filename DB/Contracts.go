package DB

import (
	"net/http"
)

const ContractCollection = "Contract"

type Contract struct {
	ContractID           string `db:"ContractID"`
	Filepath             string `db:"filepath"`
	Status               string `db:"status"`
	ContractcreationTime string `db:"creationTime"`
	Creator              User   `db:"Creator"`
	DelStatus            string `db:"delStatus"`
	UpdateTime           string `db:"updateTime"`
	ContractName         string `db:"contractName"`
	ExpirationTime       string `db:"ExpirationTime"`
	Blockchain           int    `db:"Blockchain"`
	Message              string `db:"Message"`
}

func (d *dbServer) CreateContract(w http.ResponseWriter, r *http.Request) {

}
