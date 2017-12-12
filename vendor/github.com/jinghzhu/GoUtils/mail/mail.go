package mail

import (
	"errors"
	"fmt"
	"net/smtp"
	"strings"

	catcher "github.com/jinghzhu/GoUtils/panic"
)

const (
	server          = "mail.server.com:25"
	contentTypeText = "Content-Type: text/plain; charset=UTF-8"
)

func Send(from string, fromAddr string, to []string, subject string, body string) (err error) {
	defer catcher.CatchPanic(&err)
	if len(to) == 0 {
		return errors.New("Need 1 recipient at least.")
	}
	if from == "" || !strings.Contains(fromAddr, "@") {
		errMsg := fmt.Sprintf("The mail sender is not correct: from = %s, from_address = %q\n", from, fromAddr)
		return errors.New(errMsg)
	}

	msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " + from +
		"<" + fromAddr + ">\r\nSubject: " + subject + "\r\n" + contentTypeText + "\r\n\r\n" + body)
	return smtp.SendMail(server, nil, from, to, msg)

}
