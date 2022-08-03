package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"os"

	gomail "gopkg.in/gomail.v2"

	"github.com/joho/godotenv"

	"github.com/labstack/echo"
)

type info struct {
	Name string
}

func main() {

	e := echo.New()
	e.GET("/", hello)
	e.POST("/", sendMail)
	e.POST("/gomail", sendMailByGoMail)

	e.Logger.Fatal(e.Start(":3000"))

}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Main Page")
}

func sendMail(c echo.Context) error {

	envErr := godotenv.Load(".env")
	if envErr != nil {
		fmt.Printf("Cound not load env")
		os.Exit(1)
	}
	envMap, MapErr := godotenv.Read(".env")
	if MapErr != nil {
		fmt.Printf("Could not Read")
		os.Exit(1)
	}

	// Sender data. เปิดใช้ generate รหัสที่ใช้กับ App ในการส่ง email SMTP กับ Gmail
	from := envMap["senderEmail"]
	//รหัส SMTP
	password := envMap["smtpPwd"]

	receiverEmail := c.FormValue("email")

	// Receiver email address.
	to := []string{
		receiverEmail,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, _ := template.ParseFiles("email.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: หัวข้ออีเมล \n%s\n\n", mimeHeaders)))

	var i = info{"Newzyz"}
	t.Execute(&body, i)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return err
	}
	type Message struct {
		Message string `json:"message"`
	}
	fmt.Println("Email Sent!")
	var message = Message{"Email Sented"}
	return c.JSON(http.StatusOK, message)
}

// GoMail แนบไฟล์ได้
func sendMailByGoMail(c echo.Context) error {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		fmt.Printf("Cound not load env")
		os.Exit(1)
	}
	envMap, MapErr := godotenv.Read(".env")
	if MapErr != nil {
		fmt.Printf("Could not Read")
		os.Exit(1)
	}

	// Sender data. เปิดใช้ generate รหัสที่ใช้กับ App ในการส่ง email SMTP กับ Gmail
	from := envMap["senderEmail"]
	//รหัส SMTP
	password := envMap["smtpPwd"]

	receiverEmail := c.FormValue("email")

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", from)

	// Set E-Mail receivers
	m.SetHeader("To", receiverEmail)

	// Set E-Mail subject
	m.SetHeader("Subject", "Gomail test subject")

	t := template.New("email.html")

	var err error
	t, err = t.ParseFiles("email.html")
	if err != nil {
		log.Println(err)
	}

	var i = info{"Newzyz"}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, i); err != nil {
		log.Println(err)
	}

	result := tpl.String()

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/html", result)

	// Attach Files
	// m.Attach("./img/rabbit.jpg")
	// m.Attach("./img/rabbit.jpg")

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, from, password)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}
	type Message struct {
		Message string `json:"message"`
	}
	fmt.Println("Email Sent!")
	var message = Message{"Email Sented"}
	return c.JSON(http.StatusOK, message)
}
