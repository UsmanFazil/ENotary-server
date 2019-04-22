package DB

import (
	"encoding/json"
	"net/http"
)

// func (d *dbServer) SignContract(w http.ResponseWriter, r *http.Request) {

// 	var sc SignContract
// 	var signer Signer
// 	var signers []Signer
// 	_ = json.NewDecoder(r.Body).Decode(&sc)

// 	tokenstring := r.Header["Token"][0]
// 	claims, cBool := GetClaims(tokenstring)
// 	if !cBool {
// 		RenderError(w, "Invalid user request")
// 		Logger("Invalid user request")
// 		return
// 	}
// 	uID := claims["userid"].(string)

// 	contractCollection := d.sess.Collection(ContractCollection)
// 	signerCollection := d.sess.Collection(SignerCollection)
// 	userCollection := d.sess.Collection(UserCollection)

// 	res := contractCollection.Find(db.Cond{"ContractID": sc.ContractID})
// 	total, _ := res.Count()
// 	if total != 1 {
// 		RenderError(w, "CONTRACT NOT FOUND")
// 		return
// 	}

// 	res1 := signerCollection.Find(db.Cond{"ContractID": sc.ContractID, "userID": uID})
// 	total1, _ := res1.Count()
// 	if total1 != 1 {
// 		RenderError(w, "CONTRACT NOT FOUND")
// 		return
// 	}
// 	res1.One(&signer)

// 	if signer.SignStatus != "pending" {
// 		RenderError(w, "Contract Already signed by User")
// 		return
// 	}
// 	signer.SignStatus = "Signed"
// 	signer.SignDate = time.Now().Format(time.RFC850)
// 	res.Update(signer)

// 	res2 := signerCollection.Find(db.Cond{"ContractID": sc.ContractID})
// 	res2.All(&signer)

// 	signedusers := 0
// 	for i := 0; i< len(signers); i ++ {
// 		if signers[i].SignStatus
// 	}

// }

func (d *dbServer) SaveCoordinates(w http.ResponseWriter, r *http.Request) {
	var pi []PlaygroundInput
	var signerCords []Coordinates
	_ = json.NewDecoder(r.Body).Decode(&pi)
	json.NewEncoder(w).Encode(pi)

}

// [ { "width": 50, "height": 50, "top": 10, "left": 10, "draggable": true, "resizable": false, "minw": 10, "minh": 10, "axis": "both", "parentLim": true, "snapToGrid": false, "aspectRatio": false, "zIndex": 1, "color": "lightblue url('http://localhost:8000/Files/User_signs/Signs/997c679b-7f54-4908-b254-2925c51d8889.png') no-repeat fixed center", "active": false, "userid": "", "text": "Signature", "recipient": "", "recipientname": "" },
//  { "width": 50, "height": 50, "top": 10, "left": 10, "draggable": true, "resizable": false, "minw": 10, "minh": 10, "axis": "both", "parentLim": true, "snapToGrid": false, "aspectRatio": false, "zIndex": 1, "color": "lightblue url('http://localhost:8000/Files/User_signs/Signs/997c679b-7f54-4908-b254-2925c51d8889.png') no-repeat fixed center", "active": false, "userid": "", "text": "Initial", "recipient": "", "recipientname": "" },
//  { "width": 50, "height": 50, "top": 10, "left": 10, "draggable": true, "resizable": false, "minw": 10, "minh": 10, "axis": "both", "parentLim": true, "snapToGrid": false, "aspectRatio": false, "zIndex": 1, "color": "lightblue url('http://localhost:8000/Files/User_signs/Signs/997c679b-7f54-4908-b254-2925c51d8889.png') no-repeat fixed center", "active": false, "userid": "", "text": "Calender", "recipient": "", "recipientname": "" },
//  { "width": 50, "height": 50, "top": 10, "left": 10, "draggable": true, "resizable": false, "minw": 10, "minh": 10, "axis": "both", "parentLim": true, "snapToGrid": false, "aspectRatio": false, "zIndex": 1, "color": "lightblue url('http://localhost:8000/Files/User_signs/Signs/997c679b-7f54-4908-b254-2925c51d8889.png') no-repeat fixed center", "active": false, "userid": "", "text": "Email", "recipient": "", "recipientname": "" },
//  { "width": 50, "height": 50, "top": 10, "left": 10, "draggable": true, "resizable": false, "minw": 10, "minh": 10, "axis": "both", "parentLim": true, "snapToGrid": false, "aspectRatio": false, "zIndex": 1, "color": "lightblue url('http://localhost:8000/Files/User_signs/Signs/997c679b-7f54-4908-b254-2925c51d8889.png') no-repeat fixed center", "active": false, "userid": "", "text": "Company", "recipient": "", "recipientname": "" },
//  { "width": 50, "height": 50, "top": 10, "left": 10, "draggable": true, "resizable": false, "minw": 10, "minh": 10, "axis": "both", "parentLim": true, "snapToGrid": false, "aspectRatio": false, "zIndex": 1, "color": "lightblue url('http://localhost:8000/Files/User_signs/Signs/997c679b-7f54-4908-b254-2925c51d8889.png') no-repeat fixed center", "active": false, "userid": "", "text": "Text", "recipient": "", "recipientname": "" } ]
