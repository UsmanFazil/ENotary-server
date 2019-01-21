package Email

import (
	gomail "gopkg.in/gomail.v2"
)

func SendMail(usermail string, vcode string) (bool, error) {

	m := gomail.NewMessage()
	m.SetHeader("From", "eNotaryOfficial@gmail.com")
	m.SetHeader("To", usermail)
	m.SetHeader("Subject", "E-Notary says Hi !!")
	m.SetBody("text/html", signupMsg(vcode))

	d := gomail.NewDialer("smtp.gmail.com", 587, "eNotaryOfficial@gmail.com", "Enotary360")

	if err := d.DialAndSend(m); err != nil {
		return false, err
	}
	return true, nil
}

func signupMsg(vcode string) string {
	return "Hello! Welcome to E-Notary Platform. <br/> Please use follwoing verification code to verify your email, Thanks. <br/> 	This Code will expire in one Hour. <br/> Verification Code = " + "<b>" + vcode + "<b>"
}
