package DB

import (
	"net/http"
	"time"
)

func (d *dbServer) Logout(w http.ResponseWriter, r *http.Request) {

	var blackList BlackList
	blackList.TokenString = r.Header["Token"][0]
	blackList.ExpTime = time.Now().Add(time.Hour * 2).Format(time.RFC850)

	Collection := d.sess.Collection(BlackListCollection)

	_, err := Collection.Insert(blackList)

	if err != nil {
		RenderError(w, "CAN NOT LOGOUT USER! TRY AGAIN")
		return
	}
	RenderResponse(w, "LOGGED OUT ", http.StatusOK)
	return
}
