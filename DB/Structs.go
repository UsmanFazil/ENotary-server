package DB

const UserCollection = "Users"
const ContractCollection = "Contract"
const SignerCollection = "Signer"
const VerifCollection = "Verification"
const FolderCollection = "Folder"
const ContractFolderCollection = "ContractFolder"
const BlackListCollection = "BlackList"
const WalletsCollection = "Wallets"
const CoordinatesCol = "Coordinates"
const ContractHtmlCol = "ContractHtml"

const Profilepicspath = "Files/Profile_pics"
const Signpath = "Files/User_signs/Signs"
const InitialsPath = "Files/User_signs/Initials"
const Contractfilepath = "Files/Contracts"
const CSVpath = "Files/CSV"
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
	ContractID  string `json:"cid"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	NeedtoSign  bool   `json:"NeedtoSign"`
	ReceiveCopy bool   `json:"Receiveacopy"`
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

type Preferences struct {
	UserName string `json:"username"`
	Company  string `json:"company"`
	Phone    string `json:"phone"`
}

type ContractDetail struct {
	ContractData Contract
	Signers      []Signer
}

type ContractDetailHash struct {
	ContractData Contract
	Signers      []Signer
	ContractHash string
}

type ContractBasic struct {
	ContractID string
	Path       string
}

type ManageScreen struct {
	InboxContracts []Contract
	FolderList     []Folder
}

type SignRes struct {
	Signpath     string
	InitialsPath string
}

type SendContract struct {
	ContractID string `json:"ContractID"`
	EmailSubj  string `json:"EmailSubj"`
	EmailMsg   string `json:"EmailMsg"`
}

type WalletInfo struct {
	Userid        string `db:"userid"`
	PublicAddress string `db: "walletaddress"`
}

type SaveWalletinput struct {
	ContractID    string `json: "ContractID"`
	UserID        string `json:"userid"`
	PublicAddress string `json:"publicAddress"`
}

type Base64 struct {
	SignBase64     string
	InitialsBase64 string
}

type SignContract struct {
	FileBase64 string `json: "FileBase64"`
	ContractID string `json: "ContractID"`
}

type Coordinates struct {
	ContractID string `db:"ContractID"`
	UserID     string `db: "userID"`
	Name       string `db:"name"`
	Topcord    int    `db: "topcord"`
	Leftcord   int    `db:"leftcord"`
}

type PlaygroundInput struct {
	Contractid string `json:"contractid"`
	Top        int    `json :"top"`
	Left       int    `json :"left"`
	Recipient  string `json :"recipient"`
	Text       string `json: "text"`
}

type Testing struct {
	ContractID string `db:"ContractID"`
	UserID     string `db:"userID"`
	Name       string `db:"name"`
	Topcord    string `db: "topcord"`
	Leftcord   string `db: "leftcord"`
}

//CONTRACT STATUS TYPES
// DRAFT
// In Progress
// Completed
// Declined
// Voided

//Signer Status
// pending
// Signed
// Declined
