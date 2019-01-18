package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

const uploadPath = "./tmp"

func main() {
	http.HandleFunc("/up", foo)

	fs := http.FileServer(http.Dir(uploadPath))
	http.Handle("/files/", http.StripPrefix("/files", fs))

	fmt.Println("server at 8000")
	http.ListenAndServe(":8000", nil)
}

func foo(w http.ResponseWriter, r *http.Request) {
	var s string

	if r.Method == http.MethodPost {
		f, _, err := r.FormFile("userfile")
		if err != nil {
			log.Fatal("error")
			return
		}
		defer f.Close()

		bs, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal("error")
			return
		}
		filetype := http.DetectContentType(bs)
		if filetype != "image/jpeg" && filetype != "image/jpg" &&
			filetype != "image/gif" && filetype != "image/png" &&
			filetype != "application/pdf" {
			renderError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
			return
		}
		fileName := randToken(12)
		fileEndings, err := mime.ExtensionsByType(filetype)
		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}
		s = string(bs)
		newpath := filepath.Join(uploadPath, fileName+fileEndings[0])
		file, err := os.Create(newpath)
		fmt.Println(fileEndings[0])
		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}
		defer file.Close()
		file.WriteString(s)

		json.NewEncoder(w).Encode("file uploaded successfully")

	}

}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
