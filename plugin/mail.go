package plugin

import (
	"fmt"
	"os"
	"regexp"
	"net/smtp"
	"strings"
)

const (
	Q1 = iota    // 输入收件人
	Q2 = iota    // 输入邮件标题
	Q3 = iota    // 输入邮件正文
)

type Mail struct{}

type Params struct {
	stage   int32
	address string
	title   string
	body    string
}

func (w *Mail) Matches(text string) bool {
	return strings.HasPrefix(text, "mail") ||
	strings.Contains(text, "发邮件") ||
	strings.Contains(text, "发送邮件")
}

func (w *Mail) Respond(msg *Message) error {
	var err error
	switch msg.Status {
	case S_INIT:
		msg.Status = S_PROCESSING
		msg.Handler = w
		params := new(Params)
		tokens := strings.Fields(msg.Text)
		if len(tokens) < 3 {
			msg.Send("你想给谁发邮件?")
			params.stage = Q1
		} else {
			params.address, err = parse_mail_address(tokens[2])
			if err != nil {
				msg.Send("邮箱地址好像不对哦!")
				params.stage = Q1
			} else {
				msg.Send("邮件标题是什么?")
				params.stage = Q2
			}
		}
		msg.params = params
	case S_PROCESSING:
		params, _ := msg.params.(*Params)
		switch params.stage {
		case Q1:
			params.address, err = parse_mail_address(msg.Text)
			if err != nil {
				msg.Send("邮箱地址好像不对哦!")
				params.stage = Q1
			} else {
				msg.Send("邮件标题是什么?")
				params.stage = Q2
			}
		case Q2:
			params.title = msg.Text
			msg.Send("邮件正文是什么?")
			params.stage = Q3
		case Q3:
			params.body = msg.Text
			err := send_mail(params)
			if err != nil {
				msg.ResetSession()
				return err
			} else {
				msg.Send("发送邮件成功!")
			}
			msg.ResetSession()
		}
	}

	msg.Done()
	return nil
}

func (w *Mail) Help() string {
	return "mail - 发送邮件."
}

func send_mail(p *Params) error {
	auth := smtp.PlainAuth(
		"",
		"t34@qq.com",
		os.Getenv("SLACK_BOT_MAIL_PASSWORD"),
		"smtp.qq.com",
	)

	var message string
	message += fmt.Sprintf("From: %s\n", "t34@qq.com")
	message += fmt.Sprintf("To: %s\n", p.address)
	message += fmt.Sprintf("Subject: %s\n", p.title)
	message += fmt.Sprintf("\n%s\n", p.body)
	err := smtp.SendMail(
		"smtp.qq.com:25",
		auth,
		"t34@qq.com",
		[]string{p.address},
		[]byte(message),
	)

	return err
}

func parse_mail_address(text string) (string, error) {
	mail_addr_pattern := `(\w[-._\w]*\w@\w[-._\w]*\w\.\w{2,3})`
	r := regexp.MustCompile("^<mailto:" + mail_addr_pattern + ".*")
	s := r.FindStringSubmatch(text)
	if len(s) == 0 {
		return "", fmt.Errorf("parse email address fail: %s", text)
	}
	return s[1], nil
}
