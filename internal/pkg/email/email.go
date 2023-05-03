package email

import (
	"be-idx-tsg/internal/app/httprest/model"
	"bytes"
	"crypto/tls"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	URL      string
	Name     string
	UserName string
	Subject  string
	Password string
}

// Email template parser
func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

func SendEmail(user *model.AuthConfirmPasswordResponse, data *EmailData, emailTemp string) {

	// Sender data.
	from := "admin1@admin.com"
	smtpPass := "A7UBMHZDZEPV32E5QL4RRGW3NI"
	smtpUser := "CGQ4D3QKOZDEIXVR5IEBYTGFNI"
	to := user.Email
	smtpHost := "smtp.develmail.com"
	smtpPort := 587

	var body bytes.Buffer

	template, err := ParseTemplateDir("templates")
	if err != nil {
		log.Fatal("Could not parse template", err)
	}

	errExecuteTemplate := template.ExecuteTemplate(&body, emailTemp, &data)
	if errExecuteTemplate != nil{
		log.Fatal("Could not Excute Template: ", errExecuteTemplate)
	}
	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(m); err != nil {
		log.Fatal("Could not send email: ", err)
	}
}

func SendEmail2(email string, data *EmailData, emailTemp string) {

	// Sender data.
	from := "admin1@admin.com"
	smtpPass := "A7UBMHZDZEPV32E5QL4RRGW3NI"
	smtpUser := "CGQ4D3QKOZDEIXVR5IEBYTGFNI"
	to := email
	smtpHost := "smtp.develmail.com"
	smtpPort := 587

	var body bytes.Buffer

	template, err := ParseTemplateDir("templates")
	if err != nil {
		log.Fatal("Could not parse template", err)
	}

	errExecuteTemplate := template.ExecuteTemplate(&body, emailTemp, &data)
	if errExecuteTemplate != nil{
		log.Fatal("Could not Excute Template: ", errExecuteTemplate)
	}
	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(m); err != nil {
		log.Fatal("Could not send email: ", err)
	}
}
