package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
)

type Email struct {
	To      string
	Subject string
	Msg     string
}

var pass = os.Getenv("GMAIL_APP_PASSWORD")

func Send(e Email) error {
	hostAddress := "smtp.gmail.com"
	from := "grafana1000@gmail.com"

	fullServerAddress := hostAddress + ":465"
	headerMap := make(map[string]string)
	headerMap["From"] = from
	headerMap["To"] = e.To
	headerMap["Cc"] = from
	headerMap["Subject"] = e.Subject
	headerMap["MIME-version"] = "1.0"
	headerMap["Content-Type"] = "text/html; charset=\"UTF-8\""
	msg := ""
	for k, v := range headerMap {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	msg += e.Msg

	authenticate := smtp.PlainAuth("", from, pass, hostAddress)

	tlsConfigurations := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: false,
		ServerName:         hostAddress,
	}
	conn, err := tls.Dial("tcp", fullServerAddress, tlsConfigurations)
	if err != nil {
		return err
	}
	client, err := smtp.NewClient(conn, hostAddress)
	if err != nil {
		return err
	}
	// Auth
	if err = client.Auth(authenticate); err != nil {
		return err
	}
	// To && From
	if err = client.Mail(from); err != nil {
		return err
	}
	if err = client.Rcpt(e.To); err != nil {
		return err
	}
	if err = client.Rcpt(from); err != nil { // CC email message
		return err
	}
	// Data
	writer, err := client.Data()
	if err != nil {
		return err
	}
	if _, err = writer.Write([]byte(msg)); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	if err := client.Quit(); err != nil {
		return err
	}
	return nil
}
