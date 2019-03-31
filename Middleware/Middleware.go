package MiddleWare

import (
	"ENOTARY-Server/DB"
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte("secretkey")

// Middle ware to check every request (JWT VALIDATION)
func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {

			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return mySigningKey, nil
			})

			if err != nil {
				DB.RenderError(w, "Invalid user request")
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {
			DB.RenderError(w, "Invalid user request")
		}
	})
}
