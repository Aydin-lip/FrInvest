package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"path/filepath"
	"recruitment-api/config"
	"strconv"
)

type Mailer struct {
	host     string
	port     int
	username string
	password string
	from     string
}

func NewMailer() *Mailer {
	port, _ := strconv.Atoi(config.AppConfig.SMTPPort)
	return &Mailer{
		host:     config.AppConfig.SMTPHost,
		port:     port,
		username: config.AppConfig.SMTPUsername,
		password: config.AppConfig.SMTPPassword,
		from:     config.AppConfig.SMTPUsername,
	}
}

func (m *Mailer) sendEmail(to, subject, body string) error {
	fmt.Println(m)
	auth := smtp.PlainAuth("", m.username, m.password, m.host)

	headers := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"utf-8\"\r\n\r\n",
		m.from, to, subject)

	msg := []byte(headers + body)

	addr := fmt.Sprintf("%s:%d", m.host, m.port)
	return smtp.SendMail(addr, auth, m.from, []string{to}, msg)
}

func (m *Mailer) SendVerificationEmail(toEmail, token string) error {
	tmplPath := filepath.Join("templates", "verification_email.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to parse verification template: %w", err)
	}

	verifyLink := fmt.Sprintf("%s:%s/api/auth/verify-email?token=%s", config.AppConfig.AppUrl, config.AppConfig.AppPort, token)

	data := map[string]string{
		"VerifyLink": verifyLink,
		"Email":      toEmail,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute verification template: %w", err)
	}

	return m.sendEmail(toEmail, "Verify Your Email", buf.String())
}

func (m *Mailer) SendWebinarEmail(toEmail, firstName string) error {
	tmplPath := filepath.Join("templates", "webinar_email.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to parse webinar template: %w", err)
	}

	data := map[string]string{
		"FirstName":       firstName,
		"Email":           toEmail,
		"WebinarDateTime": config.AppConfig.WebinarDateTime,
		"WebinarLink":     config.AppConfig.WebinarLink,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute webinar template: %w", err)
	}

	return m.sendEmail(toEmail, "Webinar Registration Confirmed", buf.String())
}
