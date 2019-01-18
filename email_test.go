package learning

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/smtp"
	"os"
	"testing"

	"github.com/bradfitz/go-smtpd/smtpd"
	gomail "gopkg.in/gomail.v2"
)

type emailReceiver struct{}

func (e *emailReceiver) AddRecipient(rcpt smtpd.MailAddress) error {
	return nil
}

func (e *emailReceiver) BeginData() error {
	return nil
}

func (e *emailReceiver) Write(line []byte) error {
	return nil
}

func (e *emailReceiver) Close() error {
	return nil
}

func runSMTPServer(t *testing.T) {
	srv := smtpd.Server{
		OnNewMail: func(c smtpd.Connection, from smtpd.MailAddress) (smtpd.Envelope, error) {
			return &emailReceiver{}, nil
		},
	}
	srv.ListenAndServe()
}

func TestSMTPClient(t *testing.T) {
	go runSMTPServer(t)

	c, err := smtp.Dial("localhost:25")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	if err := c.Mail("sender@example.org"); err != nil {
		t.Fatal(err)
	}
	if err := c.Rcpt("recipient@example.net"); err != nil {
		t.Fatal(err)
	}
	wc, err := c.Data()
	if err != nil {
		t.Fatal(err)
	}
	_, err = fmt.Fprintf(wc, "This is the email body")
	if err != nil {
		t.Fatal(err)
	}
	err = wc.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = c.Quit()
	if err != nil {
		t.Fatal(err)
	}
}

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unkown fromServer")
		}
	}
	return nil, nil
}

func TestConnectingToOutlook(t *testing.T) {
	t.Skip("In favor of using gomail")
	pwd := os.Getenv("OUTLOOK_PASSWORD")
	if pwd == "" {
		t.Skip("OUTLOOK_PASSWORD not set. (OUTLOOK_PASSWORD=PWD go test -v email_test.go)")
	}
	hostname := "outlook.office365.com"
	auth := LoginAuth("arun.barua@e2open.com", pwd)

	from := "arun.barua@e2open.com"
	recipients := []string{"arun.barua@e2open.com", "arunsworld@gmail.com"}
	msg := []byte(`To: arun.barua@e2open.com
Subject: this is a great test.

This is the email body!!`)

	err := smtp.SendMail(hostname+":587", auth, from, recipients, msg)
	if err != nil {
		t.Fatal(err)
	}
}

func stopOnError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestConnectingToOutlookExplicit(t *testing.T) {
	t.Skip("Explicit not required if SendMail works...")
	pwd := os.Getenv("OUTLOOK_PASSWORD")
	if pwd == "" {
		t.Skip("OUTLOOK_PASSWORD not set. (OUTLOOK_PASSWORD=PWD go test -v email_test.go)")
	}
	hostname := "outlook.office365.com"
	c, err := smtp.Dial(hostname + ":587")
	stopOnError(t, err)
	defer c.Close()

	config := tls.Config{
		ServerName: hostname,
	}
	err = c.StartTLS(&config)
	stopOnError(t, err)

	auth := LoginAuth("arun.barua@e2open.com", pwd)
	err = c.Auth(auth)
	stopOnError(t, err)

	err = c.Mail("arun.barua@e2open.com")
	stopOnError(t, err)

	err = c.Rcpt("arun.barua@e2open.com")
	stopOnError(t, err)
	err = c.Rcpt("arunsworld@gmail.com")
	stopOnError(t, err)
}

func TestSendingViaOutlookUsingGOMAIL(t *testing.T) {
	pwd := os.Getenv("OUTLOOK_PASSWORD")
	if pwd == "" {
		t.Skip("OUTLOOK_PASSWORD not set. (OUTLOOK_PASSWORD=PWD go test -v email_test.go)")
	}
	hostname := "outlook.office365.com"

	m := gomail.NewMessage()
	m.SetHeader("From", "arun.barua@e2open.com")
	m.SetHeader("To", "arun.barua@e2open.com", "arunsworld@gmail.com")
	m.SetHeader("Subject", "Email from gomail")
	m.SetBody("text/plain", "Body with plain text")

	d := gomail.NewDialer(hostname, 587, "arun.barua@e2open.com", pwd)

	if err := d.DialAndSend(m); err != nil {
		t.Fatal(err)
	}
}
