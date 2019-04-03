package DB

import (
	"fmt"
	"net/http"

	"upper.io/db.v3"

	jwt "github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte("secretkey")

func (d *dbServer) IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {

			Collection := d.sess.Collection(BlackListCollection)

			res := Collection.Find(db.Cond{"token": r.Header["Token"][0]})
			total, _ := res.Count()

			if total > 0 {
				RenderError(w, "Invalid user request")
				return
			}

			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return mySigningKey, nil
			})

			if err != nil {
				RenderError(w, "Invalid user request")
				return
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {
			RenderError(w, "Invalid user request")
		}
	})
}
