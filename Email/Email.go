package Email

import (
	"log"

	gomail "gopkg.in/gomail.v2"
)

func SendMail(usermail string, vcode string) {

	m := gomail.NewMessage()
	m.SetHeader("From", "eNotaryOfficial@gmail.com")
	m.SetHeader("To", usermail)
	m.SetHeader("Subject", "E-Notary says Hi !!")
	m.SetBody("text/html", signupMsg(vcode))

	d := gomail.NewDialer("smtp.gmail.com", 587, "eNotaryOfficial@gmail.com", "Enotary360")

	if err := d.DialAndSend(m); err != nil {
		log.Println("CAN NOT GENERATE EMAIL:", err)
		return
	}
	log.Println("EMAIL SENT SUCCESSFULLY")

}

func signupMsg(vcode string) string {
	return "Hello! <br/> Please use following verification code to verify your email, Thanks. <br/> 	This Code will expire in two Hour. <br/> Verification Code = " + "<b>" + vcode + "<b>"
}
