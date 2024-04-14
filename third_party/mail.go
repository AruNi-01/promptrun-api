package third_party

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"os"
	"strconv"
)

var (
	Mail *gomail.Dialer
)

func MailInit() {
	d := gomail.NewDialer(
		os.Getenv("MAIL_HOST"),
		func() int {
			port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
			return port
		}(),
		os.Getenv("MAIL_USER"),
		os.Getenv("MAIL_AUTH_CODE"),
	)
	Mail = d
}

func SendMailForBecomeSeller(userName, email string) (bool, error) {
	m := gomail.NewMessage()
	m.SetHeader("From", "PromptRun | promptrun.shop <"+os.Getenv("MAIL_USER")+">")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "祝贺，您已成为 PromptRun 卖家！")
	message := `
		<div style="display: flex; justify-content: center; margin: 35px;">
		  <div style="width: 80%%; text-align: center; border: 2px solid #000; border-radius: 15px; padding: 20px;">
			<h1>Hello %s. </h1>
			<h3>感谢您的等待，非常高兴的通知您，您的申请已通过，现已成为 PromptRun 卖家中的一员，快去发布您的 Prompt 吧！</h3>
			<h3 style="color: red; font-weight: bold;">注意：成为卖家后，您需要在 PromptRun 平台重新登录，以便获取卖家权限。</h3>
			<h3>祝您在 PromptRun 平台上获得更多的收益！</h3>
			<h5>—PromptRun 团队</h5>
		  </div>
		</div>
	`
	m.SetBody("text/html", fmt.Sprintf(message, userName))
	if err := Mail.DialAndSend(m); err != nil {
		return false, fmt.Errorf("send mail failed, err: %s", err.Error())
	}
	return true, nil
}
