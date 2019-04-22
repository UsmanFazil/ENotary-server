package DB

import (
	"encoding/json"
	"net/http"

	db "upper.io/db.v3"
)

func (d *dbServer) VerifyContract(w http.ResponseWriter, r *http.Request) {
	var wi WalletInfo
	var user User
	_ = json.NewDecoder(r.Body).Decode(&wi)

	walletcollection := d.sess.Collection(WalletsCollection)
	usercollection := d.sess.Collection(UserCollection)

	res := walletcollection.Find(db.Cond{"walletaddress": wi.PublicAddress})
	err := res.One(&wi)

	if err != nil {
		RenderError(w, "No account found")
		return
	}

	res1 := usercollection.Find(db.Cond{"userid": wi.Userid})

	errstring := res1.One(&user)
	if errstring != nil {
		RenderError(w, "No account found")
		return
	}

	json.NewEncoder(w).Encode(user.Email)
	return
}

// func (d *dbServer) VerifyContract(w http.ResponseWriter, r *http.Request) {
// 	var walletinfo WalletInfo

// 	_ = json.NewDecoder(r.Body).Decode(&walletinfo)
// 	fmt.Println(walletinfo)
// }
