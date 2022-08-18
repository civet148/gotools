package smtp

import (
	"crypto/tls"
	"net"
	"net/smtp"
	"strings"
)

type SmtpHelper struct {
	UserName string `json:"username" xml:"username" form:"username"`
	PassWord string `json:"password" xml:"password" form:"password"`
	HostName string `json:"hostname" xml:"hostname" form:"hostname"`
	To       string `json:"sendto" xml:"sendto" form:"sendto"`
	Cc       string `json:"cc" xml:"cc" form:"cc"`
	Bcc      string `json:"bcc" xml:"bcc" form:"bcc"`
	Subject  string `json:"subject" xml:"subject" form:"subject"`
	Body     string `json:"body" xml:"body" form:"body"`
	IsHtml   bool   `json:"ishtml" xml:"ishtml" form:"ishtml"`
}

func (helper *SmtpHelper) SendMail() error {
	hp := strings.Split(helper.HostName, ":")
	auth := smtp.PlainAuth("", helper.UserName, helper.PassWord, hp[0])
	var content_type string
	if helper.IsHtml == true {
		content_type = "Content-Type: text/html; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain; charset=UTF-8"
	}

	strmsg := "From: " + helper.UserName + "\r\n" + "To: " + helper.To + "\r\n"
	send_to := strings.Split(helper.To, ";")
	if len(helper.Cc) > 0 {
		strmsg += "Cc: " + helper.Cc + "\r\n"
		send_to = append(send_to, strings.Split(helper.Cc, ";")...)
	}
	if len(helper.Bcc) > 0 {
		strmsg += "Bcc: " + helper.Bcc + "\r\n"
		send_to = append(send_to, strings.Split(helper.Bcc, ";")...)
	}
	strmsg += "Subject: " + helper.Subject + "\r\n" + content_type + "\r\n\r\n" + helper.Body

	msg := []byte(strmsg)
	err := smtp.SendMail(helper.HostName, auth, helper.UserName, send_to, msg)
	return err
}

// tls 连接smtp服务器

//return a smtp client
func TlsDial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

//参考net/smtp的func SendMail()
//使用net.Dial连接tls(ssl)端口时,smtp.NewClient()会卡住且不提示err
//len(to)>1时,to[1]开始提示是密送
func SendMailUsingTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {

	//create smtp client
	c, err := TlsDial(addr)
	if err != nil {
		return err
	}
	defer c.Close()

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				return err
			}
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}

func (helper *SmtpHelper) SendMailTls() error {
	hp := strings.Split(helper.HostName, ":")
	auth := smtp.PlainAuth("", helper.UserName, helper.PassWord, hp[0])
	var content_type string
	if helper.IsHtml == true {
		content_type = "Content-Type: text/html; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain; charset=UTF-8"
	}

	msg := []byte("From: " + helper.UserName + "\r\n" +
		"To: " + helper.To + "\r\n" +
		"Cc: " + helper.Cc + "\r\n" +
		"Bcc: " + helper.Bcc + "\r\n" +
		"Subject: " + helper.Subject + "\r\n" + content_type + "\r\n\r\n" +
		helper.Body)

	send_to := strings.Split(helper.To, ";")
	send_to = append(send_to, strings.Split(helper.Cc, ";")...)
	send_to = append(send_to, strings.Split(helper.Bcc, ";")...)
	err := SendMailUsingTLS(helper.HostName, auth, helper.UserName, send_to, msg)
	return err
}
