package smtp

import (
	"fmt"
	"testing"
)

func Test_SendMail(t *testing.T) {
	var m_smtp SmtpHelper
	m_smtp = SmtpHelper{
		UserName: "songsy@chargerlink.com",
		PassWord: "Chargerlinkssyzer0",
		HostName: "smtp.exmail.qq.com:25",
		To:       "425826621@qq.com",
		Cc:       "songsyzero@163.com",
		Bcc:      "915766980@qq.com",
		Subject:  "快点加工区",
		Body:     `时空房间的饭卡的设计感！`,
		IsHtml:   false,
	}

	fmt.Println("send mail start ... ")
	if err := m_smtp.SendMail(); err != nil {
		fmt.Println("send mail err: ", err)
	} else {
		fmt.Println("send mail success ")
	}
}
