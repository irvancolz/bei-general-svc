package email

import (
	"be-idx-tsg/internal/app/helper"
	"be-idx-tsg/internal/app/httprest/model"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
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
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("ADMIN_USERNAME")
	from := os.Getenv("ADMIN_EMAIL")
	smtpPass := os.Getenv("ADMIN_EMAIL_PASS")
	smtpHost := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(smtpPort)

	var body bytes.Buffer

	template, err := ParseTemplateDir("templates")
	if err != nil {
		log.Fatal("Could not parse template", err)
	}

	errExecuteTemplate := template.ExecuteTemplate(&body, emailTemp, &data)
	if errExecuteTemplate != nil {
		log.Fatal("Could not Excute Template: ", errExecuteTemplate)
	}
	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(smtpHost, port, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(m); err != nil {
		log.Fatal("Could not send email: ", err)
	}
}

func SendEmail2(email string, data *EmailData, emailTemp string) {

	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("ADMIN_USERNAME")
	from := os.Getenv("ADMIN_EMAIL")
	smtpPass := os.Getenv("ADMIN_EMAIL_PASS")
	smtpHost := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(smtpPort)

	var body bytes.Buffer

	template, err := ParseTemplateDir("templates")
	if err != nil {
		log.Fatal("Could not parse template", err)
	}

	errExecuteTemplate := template.ExecuteTemplate(&body, emailTemp, &data)
	if errExecuteTemplate != nil {
		log.Fatal("Could not Excute Template: ", errExecuteTemplate)
	}
	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", email)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(smtpHost, port, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(m); err != nil {
		log.Fatal("Could not send email: ", err)
	}
}

type EmailConfig struct {
	Receiver string
	Subject  string
	URL      string
}

func (c EmailConfig) SendEmailWithTemplate(wg *sync.WaitGroup, data interface{}, pathToTemplate string) {
	defer wg.Done()
	wg.Add(1)

	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("ADMIN_USERNAME")
	adminName := os.Getenv("ADMIN_EMAIL")
	smtpPass := os.Getenv("ADMIN_EMAIL_PASS")
	smtpHost := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(smtpPort)

	htmltemplate, error_generate := generateTemplatesfromHTML(data, pathToTemplate)
	if error_generate != nil {
		log.Println("failed to generate email templates :", error_generate)
		return
	}

	var htmlTemplateString = htmltemplate.String()

	if strings.Contains(os.Getenv("SMTP_HOST"), "develmail") {
		htmlTemplateString = strings.ReplaceAll(htmlTemplateString, "cid:idx-email-header.png", ImageToBase64("static/idx-email-header.png"))
	}

	message := gomail.NewMessage()
	message.SetHeader("From", adminName)
	message.SetHeader("To", c.Receiver)
	message.SetHeader("Subject", c.Subject)
	message.AddAlternative("text/plain", html2text.HTML2Text(htmlTemplateString))
	message.SetBody("text/html", htmlTemplateString)
	message.Embed("static/idx-email-header.png")

	email := gomail.NewDialer(smtpHost, port, smtpUser, smtpPass)
	email.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	result := email.DialAndSend(message)

	if result != nil {
		log.Println("failed to send email : ", result)
		return
	}
}

func ImageToBase64(directory string) string {
	bytes, err := os.ReadFile(directory)
	if err != nil {
		log.Fatal(err)
	}

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(bytes)
}

func generateTemplatesfromHTML(data interface{}, html string) (*bytes.Buffer, error) {
	var result bytes.Buffer

	template, err := ParseTemplateDir(filepath.Join("static", "templates"))
	if err != nil {
		log.Fatal("Could not parse template", err)
		return nil, err
	}

	err_create_html := template.ExecuteTemplate(&result, html, data)
	if err_create_html != nil {
		log.Println("failed to generate html : ", err_create_html)
		return nil, err_create_html
	}

	return &result, nil
}

func SendEmailNotification(receiver model.UsersIdWithEmail, subject, message string) {
	emailConfig := EmailConfig{
		Receiver: receiver.Email,
		Subject:  subject,
	}
	wg := sync.WaitGroup{}
	go emailConfig.SendEmailWithTemplate(&wg, struct {
		Message string
		Name    string
	}{Message: message, Name: receiver.Username}, "dummy-general-notif.html")
}

func GetAllUserInternalBursa(c *gin.Context) []model.UsersIdWithEmail {
	result := []model.UsersIdWithEmail{}
	dbConn, errInitDb := helper.InitDBConn("auth")
	if errInitDb != nil {
		log.Println(errInitDb)
		return result
	}
	defer dbConn.Close()

	query := ` 
		SELECT
			id,
			username,
			email
		FROM users WHERE 
		type = 'Internal'
		AND deleted_by IS NULL
	`
	queryRes, errQuery := dbConn.Queryx(query)
	if errQuery != nil {
		log.Println("failed to get user internal bursa :", errQuery)
		return result
	}
	defer queryRes.Close()

	for queryRes.Next() {
		var users model.UsersIdWithEmail
		if errScan := queryRes.StructScan(&users); errScan != nil {
			log.Println("failed to read user data :", errScan)
			return []model.UsersIdWithEmail{}
		}
		result = append(result, users)
	}

	return result
}

func GetUserAdminApp(c *gin.Context) ([]model.UsersIdWithEmail, error) {
	var result []model.UsersIdWithEmail
	dbConn, errInitDb := helper.InitDBConn("auth")

	if errInitDb != nil {
		log.Println(errInitDb)
		return nil, errInitDb
	}
	defer dbConn.Close()

	query := ` 
	SELECT
		u id,
		u username,
		u email
	FROM users u JOIN roles r ON r.id::text = u.role_id 
	WHERE r.role = 'Admin App.'
	AND u.deleted_by IS NULL
	`

	queryRes, errQuery := dbConn.Queryx(query)
	if errQuery != nil {
		log.Println("failed to get user Admin app :", errQuery)
		return nil, errQuery
	}

	defer queryRes.Close()

	for queryRes.Next() {
		var user model.UsersIdWithEmail

		if errScan := queryRes.StructScan(&user); errScan != nil {
			log.Println("failed to read user admin app :", errScan)
			return nil, errScan
		}

		result = append(result, user)
	}

	return result, nil
}

func GetUser(c *gin.Context, id string) (*model.UsersIdWithEmail, error) {
	result := model.UsersIdWithEmail{}
	dbConn, errInitDb := helper.InitDBConn("auth")

	if errInitDb != nil {
		log.Println(errInitDb)
		return nil, errInitDb
	}
	defer dbConn.Close()

	query := ` 
		SELECT
			id,
			username,
			email
		FROM users WHERE id = $1
	`
	queryRes := dbConn.QueryRowx(query, id)

	if errScan := queryRes.StructScan(&result); errScan != nil {
		log.Println("failed to read user data :", errScan)
		return nil, errScan
	}

	return &result, nil
}
