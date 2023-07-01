package util

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/config"
)

type Mail struct {
	SenderAddress string
	SenderName    string
	ToAddress     string
	ToName        string
	Subject       string
	Body          string
}

func SMTPSendMail(mailStruct Mail) error {
	confSMTP := config.Config.SmtpConfig

	servername := fmt.Sprintf("%s:%s", confSMTP.Host, confSMTP.Port)
	host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", confSMTP.Username, confSMTP.Password, host)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	if err = c.Auth(auth); err != nil {
		return err
	}

	if err = c.Mail(mailStruct.SenderAddress); err != nil {
		return err
	}

	if err = c.Rcpt(mailStruct.ToAddress); err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(BuildMessage(mailStruct)))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	c.Quit()

	return nil
}

func BuildMessage(mail Mail) string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s <%s>\r\n", mail.SenderName, mail.SenderAddress)
	msg += fmt.Sprintf("To: %s <%s>;\r\n", mail.ToName, mail.ToAddress)
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

	return msg
}
