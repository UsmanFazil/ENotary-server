package DB

import (
	"ENOTARY-Server/Email"
	"encoding/json"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/nfnt/resize"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"upper.io/db.v3"
)

func (d *dbServer) SignIt(w http.ResponseWriter, r *http.Request) {

	var contract Contract
	var user User
	var signers []Signer
	var signer Signer

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
	Signercollection := d.sess.Collection(SignerCollection)

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

	//------------ Get signer
	res3 := Signercollection.Find(db.Cond{"ContractID": contract.ContractID, "userID": uID})
	errstring := res3.One(&signer)

	if errstring != nil {
		RenderError(w, "User Not found")
		return
	}

	//------------ Place sign & coordinates

	for _, value := range AllCoordinates {
		if value.Contractid == contract.ContractID && value.Recipient == uID {
			if value.Text == "Signature" {
				errbool := ContractManipulation(contract.Filepath, user.Sign, value.Top, value.Left)
				if !errbool {
					RenderError(w, "CAN NOT SIGN CONTRACT!")
					return
				}
			}
			if value.Text == "Initial" {
				errbool := ContractManipulation(contract.Filepath, user.Initials, value.Top, value.Left)
				if !errbool {
					RenderError(w, "CAN NOT SIGN CONTRACT!")
					return
				}
			}
			if value.Text == "DateSigned" {
				addLabel(contract.Filepath, value.Left, value.Top, time.Now().Format("2006-01-02"))
			}
			if value.Text == "Name" {
				addLabel(contract.Filepath, value.Left, value.Top, user.Name)
			}
			if value.Text == "Email" {
				addLabel(contract.Filepath, value.Left, value.Top, user.Email)
			}
			if value.Text == "Company" {
				addLabel(contract.Filepath, value.Left, value.Top, user.Company)
			}
			if value.Text == "Phone" {
				addLabel(contract.Filepath, value.Left, value.Top, user.Phone)
			}
		}
	}
	// update signer status
	signer.SignStatus = "Signed"
	signer.SignDate = time.Now().Format(time.RFC850)

	err = res3.Update(signer)
	if err != nil {
		RenderError(w, "Signer status not updated")
		return
	}

	// update contract time
	contract.UpdateTime = time.Now().Format(time.RFC850)
	res1.Update(contract)

	//get other signers for email
	res4 := Signercollection.Find(db.Cond{"ContractID": contract.ContractID, "CC": 0})
	err = res4.All(&signers)

	if err != nil {
		return
	}
	var cd ContractDetail
	cd.ContractData = contract
	cd.Signers = signers

	resbool, emails := d.Getemails(cd)
	if !resbool {
		RenderResponse(w, "CONTRACT SIGNED BUT CAN NOT GENERATE EMAIL RESPONSE", http.StatusOK)
		return
	}

	count := 0
	for i := 0; i < len(signers); i++ {
		if signers[i].SignStatus == "Signed" {
			count++
		}
	}
	if count == len(signers) {
		contract.Status = "Completed"
		res1.Update(contract)
		for _, index := range emails {
			go Email.CompletedEmail(index, "CONTRACT COMPLETED", contract.ContractID)
		}

	} else {

		for _, index := range emails {
			go Email.StatusEmail(index, "CONTRACT STATUS UPDATE", contract.ContractID, false)
		}
	}

	RenderResponse(w, "Contract Signed Successfully", http.StatusOK)
	return

}

func ContractManipulation(contractpath string, signpath string, top int, left int) bool {

	image1, err := os.Open(contractpath)
	if err != nil {
		return false
	}

	first, err := jpeg.Decode(image1)
	if err != nil {
		pngManip(contractpath, signpath, top, left)
		return true
	}
	defer image1.Close()

	image2, err := os.Open(signpath)
	if err != nil {
		log.Fatalf("failed to decode: %s", err)
		return false
	}
	second, err := png.Decode(image2)
	if err != nil {
		log.Fatalf("failed to decode: %s", err)
		return false
	}
	defer image2.Close()

	m := resize.Resize(108, 96, second, resize.Lanczos3)
	offset := image.Pt(left, top)

	// b := first.Bounds()
	c := m.Bounds()

	image3 := image.NewRGBA(first.Bounds())
	draw.Draw(image3, first.Bounds(), first, image.ZP, draw.Src)
	draw.Draw(image3, c.Add(offset), m, image.ZP, draw.Over)

	// addLabel(image3, 100, 100, "ye karo na jani")

	third, err := os.Create(contractpath)
	if err != nil {
		log.Fatalf("failed to create: %s", err)
		return false
	}

	jpeg.Encode(third, image3, &jpeg.Options{jpeg.DefaultQuality})

	defer third.Close()
	return true

}

func pngManip(contractpath string, signpath string, top int, left int) bool {
	image1, err := os.Open(contractpath)
	if err != nil {
		log.Fatalf("failed to open: %s", err)
	}
	// spliter := strings.Split(contractpath, ".")
	// ext := spliter[1]

	first, err := png.Decode(image1)
	if err != nil {
		return false
	}

	defer image1.Close()

	image2, err := os.Open(signpath)
	if err != nil {
		log.Fatalf("failed to open: %s", err)
		return false
	}
	second, err := png.Decode(image2)
	if err != nil {
		log.Fatalf("failed to decode: %s", err)
		return false
	}
	defer image2.Close()

	m := resize.Resize(108, 96, second, resize.Lanczos3)

	offset := image.Pt(left, top)

	c := m.Bounds()

	image3 := image.NewRGBA(first.Bounds())
	draw.Draw(image3, first.Bounds(), first, image.ZP, draw.Src)
	draw.Draw(image3, c.Add(offset), m, image.ZP, draw.Over)

	// addLabel(image3, 100, 100, "ye karo na jani")

	third, err := os.Create(contractpath)
	if err != nil {
		log.Fatalf("failed to create: %s", err)
		return false
	}

	png.Encode(third, image3)

	defer third.Close()
	return true
}

func addLabel(contractpath string, x int, y int, label string) {

	image1, err := os.Open(contractpath)
	if err != nil {
		log.Fatalf("failed to open: %s", err)
	}

	first, err := jpeg.Decode(image1)
	if err != nil {
		addpngLabel(contractpath, x, y, label)
		return
	}
	defer image1.Close()

	image3 := image.NewRGBA(first.Bounds())
	draw.Draw(image3, first.Bounds(), first, image.ZP, draw.Src)

	col := color.RGBA{0, 0, 0, 255}
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  image3,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)

	third, err := os.Create(contractpath)
	jpeg.Encode(third, image3, &jpeg.Options{jpeg.DefaultQuality})
}

func addpngLabel(contractpath string, x int, y int, label string) {

	image1, err := os.Open(contractpath)
	if err != nil {
		log.Fatalf("failed to open: %s", err)
	}

	first, err := png.Decode(image1)
	if err != nil {
		return
	}
	defer image1.Close()

	image3 := image.NewRGBA(first.Bounds())
	draw.Draw(image3, first.Bounds(), first, image.ZP, draw.Src)

	col := color.RGBA{0, 0, 0, 255}
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  image3,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)

	third, err := os.Create(contractpath)
	png.Encode(third, image3)
}

// col := color.RGBA{0, 0, 0, 255}
// point := fixed.Point26_6{fixed.Int26_6(20 * 64), fixed.Int26_6(40 * 64)}

// d := &font.Drawer{
// 	Dst:  image3,
// 	Src:  image.NewUniform(col),
// 	Face: basicfont.Face7x13,
// 	Dot:  point,
// }
// d.DrawString("main time stamp hn")

// spliter := strings.Split(contract.Filepath, ".")
// ext := spliter[1]

// func (d *dbServer) Tester(w http.ResponseWriter, r *http.Request) {
// 	var t Testing
// 	//	var x Testing
// 	_ = json.NewDecoder(r.Body).Decode(&t)

// 	//col := d.sess.Collection(CoordinatesCol)

// 	q := d.sess.InsertInto("Coordinates").Columns("ContractID", "userID", "name", "topcord", "leftcord").Values("ba8daa9f-7e1b-470a-874e-e841602d3c5a", "c8ecfed8-1838-4f3d-b587-d04836867598", "Signature", "aaa", "bnns")
// 	_, err := q.Exec()

// 	if err != nil {
// 		fmt.Println(err)
// 	}

// }
// func (d *dbServer) ReadTester(w http.ResponseWriter, r *http.Request) {
// 	var x Testing

// 	_ = json.NewDecoder(r.Body).Decode(&x)

// 	//col := d.sess.Collection(CoordinatesCol)
// 	c := "ba8daa9f-7e1b-470a-874e-e841602d3c5a"
// 	b := "c8ecfed8-1838-4f3d-b587-d04836867598"
// 	rows, _ := d.sess.Query("SELECT * FROM Coordinates Where ContractID = ? AND userID = ?", c, b)

// 	iter := sqlbuilder.NewIterator(rows)
// 	iter.One(&x)

// 	fmt.Println(x)

// }
