package Hashing

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Filepath : struct to store path of the file on server
type Filepath struct {
	Path string `json:"path"`
}

// hexString : function to convert bytes array into hexadecimal string
func hexString(filename []byte) (string, error) {
	value := "0x" + hex.EncodeToString(filename)
	return value, nil
}

// hasher : calculates the hash of the file and outputs in bytes array
func hasher(filename string) ([]byte, error) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	h := sha256.New()
	h.Write(bs)
	a := h.Sum([]byte{})
	return a, nil
}

// Gethash : function to get ouput the final hash of a file
func Gethash(filename string) (string, error) {
	h1, err := hasher(filename)
	if err != nil {
		return "null", err
	}
	finalhash, err := hexString(h1)
	if err != nil {
		return "null", err
	}
	return finalhash, nil

}

// Servehash : function to serve file hash to the user
func Servehash(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var path Filepath
	_ = json.NewDecoder(r.Body).Decode(&path)
	filehash, err := Gethash(path.Path)

	if err != nil {
		json.NewEncoder(w).Encode("error here 2")
	}
	json.NewEncoder(w).Encode(filehash)
}
