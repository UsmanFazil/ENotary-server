package DB

const userCollection = "Users"
const ContractCollection = "Contract"
const SignerCollection = "Signer"

type User struct {
	Userid       string `db:"userid"`
	Email        string `db:"email"`
	Password     string `db:"password"`
	Name         string `db:"name"`
	Company      string `db:"company"`
	Phone        string `db:"phone"`
	Picture      string `db:"picture"`
	Sign         string `db:"sign"`
	Initials     string `db:"initials"`
	Verified     int    `db:"verified"`
	CreationTime string `db:"creationTime"`
}

type Contract struct {
	ContractID           string `db:"ContractID"`
	Filepath             string `db:"filepath"`
	Status               string `db:"status"`
	ContractcreationTime string `db:"creationTime"`
	Creator              string `db:"Creator"`
	DelStatus            int    `db:"delStatus"`
	UpdateTime           string `db:"updateTime"`
	ContractName         string `db:"contractName"`
	ExpirationTime       string `db:"ExpirationTime"`
	Blockchain           int    `db:"Blockchain"`
	Message              string `db:"Message"`
}

type Signer struct {
	ContractID    string `db:"ContractID"`
	UserID        string `db:"userID"`
	SignStatus    string `db:"SignStatus"`
	SignDate      string `db:"SignDate"`
	DeleteApprove int    `db:"DeleteApprove"`
	Message       string `db:"Message"`
	Access        int    `db:"Access"`
}

type VerifUser struct {
	UserID           string `db:"userID"`
	VerificationCode string `db:"VerificationCode"`
	ExpTime          string `db:"expTime"`
}

type JwtToken struct {
	Token string `json:"token"`
}

type Exception struct {
	Message string `json:"message"`
}

type LoginStruct struct {
	Userdata     User
	WaitingME    uint64
	WaitingOther uint64
	Token        string
}

type LogCheck struct {
	Email    string `json:"email"`
	Password string `json:password`
}
