package DB

const UserCollection = "Users"
const ContractCollection = "Contract"
const SignerCollection = "Signer"
const VerifCollection = "Verification"
const FolderCollection = "Folder"
const ContractFolderCollection = "ContractFolder"
const BlackListCollection = "BlackList"

const Profilepicspath = "Files/Profile_pics"
const Signpath = "Files/User_signs/Signs"
const InitialsPath = "Files/User_signs/Initials"
const Contractfilepath = "Files/Contracts"
const Def_pic_path = "Files/Profile_pics/default.jpeg"

const MaxpicSize = 5 * 1024 * 1024
const MaxContractSize = 10 * 1024 * 1024
const RFC850 = "Monday, 02-Jan-06 15:04:05 MST"

var MySigningKey = []byte("secretkey")

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
	Name          string `db:"name"`
	SignStatus    string `db:"SignStatus"`
	SignDate      string `db:"SignDate"`
	DeleteApprove int    `db:"DeleteApprove"`
	Message       string `db:"Message"`
	Access        int    `db:"Access"`
	CC            int    `db:"CC"`
}

type VerifUser struct {
	UserID           string `db:"userid"`
	VerificationCode string `db:"VerificationCode"`
	ExpTime          string `db:"exptime"`
}

type EmailVerf struct {
	Email            string `json:"email"`
	VerificationCode string `json:"VerificationCode"`
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
	ExpiringSoon int
	Token        string
}

type LogCheck struct {
	Email    string `json:"email"`
	Password string `json:password`
}

type Signerdata struct {
	ContractID string `json:"cid"`
	Email      string `json:"email"`
	Name       string `json:"name"`
}

type Passrecovery struct {
	Email string `json:"email"`
	Vcode string `json:"vcode"`
	Pass  string `json:"password"`
}

type Folder struct {
	FolderID     string `db:"folderID"`
	FolderName   string `db:"folderName"`
	FolderType   string `db:"folderType"`
	UserID       string `db:"userID"`
	ParentFolder string `db:"parentFolder"`
}
type ContractFolder struct {
	FolderID   string `db:"folderID"`
	ContractID string `db:"ContractID"`
}

type SearchInput struct {
	ContractName string `json:"ContractName"`
	Status       string `json:"Status"`
	Date         string `json:"Date"`
}

type BlackList struct {
	TokenString string `db:"token"`
	ExpTime     string `db:"exptime"`
}

//CONTRACT STATUS TYPES
// DRAFT
// In Progress
// Completed
// Declined
// Voided

//Signer Status
// default pending
// Signed
