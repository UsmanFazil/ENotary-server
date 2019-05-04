package DB

import (
	"encoding/json"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"

	"github.com/nfnt/resize"
	"upper.io/db.v3"
)

func (d *dbServer) SignIt(w http.ResponseWriter, r *http.Request) {

	var contract Contract
	var user User
	tokenstring := r.Header["Token"][0]
	claims, cBool := GetClaims(tokenstring)
	if !cBool {
		RenderError(w, "Invalid user request")
		Logger("Invalid user request")
		return
	}
	uID := claims["userid"].(string)

	_ = json.NewDecoder(r.Body).Decode(&contract)

	userCollection := d.sess.Collection(UserCollection)
	contractCollection := d.sess.Collection(ContractCollection)

	res := userCollection.Find(db.Cond{"userid": uID})
	err := res.One(&user)

	if err != nil {
		return
	}
	res1 := contractCollection.Find(db.Cond{"ContractID": contract.ContractID})
	err1 := res1.One(&contract)
	if err1 != nil {
		return
	}

	ContractManipulation(contract.Filepath, user.Sign)

}

func ContractManipulation(contractpath string, signpath string) {

	image1, err := os.Open(contractpath)
	if err != nil {
		log.Fatalf("failed to open: %s", err)
	}
	// spliter := strings.Split(contractpath, ".")
	// ext := spliter[1]

	first, err := jpeg.Decode(image1)
	if err != nil {
		log.Fatalf("failed to decode: %s", err)
	}
	defer image1.Close()

	image2, err := os.Open(signpath)
	if err != nil {
		log.Fatalf("failed to open: %s", err)
	}
	second, err := png.Decode(image2)
	if err != nil {
		log.Fatalf("failed to decode: %s", err)
	}
	defer image2.Close()

	m := resize.Resize(150, 150, second, resize.Lanczos3)

	offset := image.Pt(292, 410)

	// b := first.Bounds()
	c := m.Bounds()

	image3 := image.NewRGBA(first.Bounds())
	draw.Draw(image3, first.Bounds(), first, image.ZP, draw.Src)
	draw.Draw(image3, c.Add(offset), m, image.ZP, draw.Over)

	third, err := os.Create(contractpath)
	if err != nil {
		log.Fatalf("failed to create: %s", err)
	}

	jpeg.Encode(third, image3, &jpeg.Options{jpeg.DefaultQuality})
	defer third.Close()

}
