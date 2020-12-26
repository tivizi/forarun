package base

import (
	"fmt"
	"log"

	"github.com/tivizi/forarun/approot/config"
	"gopkg.in/gomail.v2"
)

var dialer *gomail.Dialer
var conf *config.SMTPConfig

func init() {
	c := config.GetContext().SMTPConfig
	conf = &c
	if conf.Enabled {
		log.Println("Email: Enabled")
		dialer = &gomail.Dialer{
			Host:     conf.Host,
			Port:     conf.Port,
			Username: conf.Account,
			Password: conf.Password,
			SSL:      conf.SSL,
		}
	}
}

// SendActiveEmail 发送激活帐号邮件
func SendActiveEmail(siteName, siteHost, activeToken, email string) error {
	activeLink := fmt.Sprintf("https://%s/active/%s", siteHost, activeToken)
	m := gomail.NewMessage()
	m.SetHeader("From", conf.Account)
	m.SetHeader("To", email)
	m.SetHeader("Subject", fmt.Sprintf("欢迎注册【%s】，激活你的帐号确认！", siteName))
	m.SetBody("text/html",
		`
		<div style="color: green;font-size: 18px;display: block;margin: 10px 0">点击下面的链接，立即激活你的帐号</div>
		<a href="`+activeLink+`">`+activeLink+`</a>
		<div style="color: #ccc">如果不是你本人操作，请忽略该邮件，不要感到困扰，谢谢</div>
	`)
	return dialer.DialAndSend(m)
}

// SendResetPasswordEmail 发送重置密码邮件
func SendResetPasswordEmail(siteName, siteHost, userName, email, activeToken, passwd string) error {
	activeLink := fmt.Sprintf("https://%s/common/reset-password/%s/%s", siteHost, passwd, activeToken)
	m := gomail.NewMessage()
	m.SetHeader("From", conf.Account)
	m.SetHeader("To", email)
	m.SetHeader("Subject", fmt.Sprintf("【%s】重置密码请求！", siteName))
	m.SetBody("text/html",
		`
		<div style="color: green;font-size: 18px;display: block;margin: 10px 0">点击下面的链接，立即重置你的密码</div>
		<a href="`+activeLink+`">`+activeLink+`</a>
		<div style="color: #ccc">如果不是你本人操作，请忽略该邮件，不要感到困扰，谢谢</div>
	`)
	return dialer.DialAndSend(m)
}
