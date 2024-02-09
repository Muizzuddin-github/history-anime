package utility

import (
	"crypto/tls"
	"errors"
	"net/smtp"
)

type Email struct {
	From    string
	To      string
	Subject string
	Html    string
}

func SendEmail(config Email) error {
	const host = "smtp.gmail.com"
	const port = "465"
	const user = "muizzuddin332@gmail.com"
	const auth = "ipna ydyx ugxi ztvs"

	// buat config tls
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         host,
	}

	// buat koneksi tls
	conn, err := tls.Dial("tcp", host+":"+port, tlsConfig)
	if err != nil {
		return errors.New(err.Error())
	}

	// buat client
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return errors.New(err.Error())
	}

	// buat authentication client
	authtentication := smtp.PlainAuth("", user, auth, host)
	err = client.Auth(authtentication)
	if err != nil {
		return errors.New(err.Error())
	}

	// sender
	err = client.Mail(user)
	if err != nil {
		return errors.New(err.Error())
	}

	// kirimkan email ke
	err = client.Rcpt(config.To)
	if err != nil {
		return errors.New(err.Error())

	}

	w, err := client.Data()
	if err != nil {
		return errors.New(err.Error())

	}

	msg := "From: " + config.From + "\n" +
		"To:" + config.To + "\n" +
		"Subject: " + config.Subject + "\n" +
		"MIME-version: 1.0\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\n\n" +
		config.Html

	_, err = w.Write([]byte(msg))
	if err != nil {
		return errors.New(err.Error())

	}

	err = w.Close()
	if err != nil {
		return errors.New(err.Error())

	}

	client.Quit()
	return nil
}
