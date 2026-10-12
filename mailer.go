package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"html"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type passwordResetMailer interface {
	SendPasswordReset(ctx context.Context, recipient, resetURL string) error
}

type smtpPasswordResetMailer struct {
	host     string
	port     int
	username string
	password string
	from     mail.Address
}

var passwordMailer passwordResetMailer
var passwordResetFrontendURL string

func initializeMailer() error {
	values := map[string]string{
		"SMTP_HOST":      strings.TrimSpace(os.Getenv("SMTP_HOST")),
		"SMTP_PORT":      strings.TrimSpace(os.Getenv("SMTP_PORT")),
		"SMTP_USERNAME":  strings.TrimSpace(os.Getenv("SMTP_USERNAME")),
		"SMTP_PASSWORD":  strings.TrimSpace(os.Getenv("SMTP_PASSWORD")),
		"SMTP_FROM":      strings.TrimSpace(os.Getenv("SMTP_FROM")),
		"SMTP_FROM_NAME": strings.TrimSpace(os.Getenv("SMTP_FROM_NAME")),
	}
	for name, value := range values {
		if value == "" {
			return fmt.Errorf("%s must be configured", name)
		}
	}
	port, err := strconv.Atoi(values["SMTP_PORT"])
	if err != nil || port < 1 || port > 65535 {
		return errors.New("SMTP_PORT must be a valid port number")
	}
	from, err := mail.ParseAddress(values["SMTP_FROM"])
	if err != nil || from.Address != values["SMTP_FROM"] {
		return errors.New("SMTP_FROM must be a valid email address")
	}
	frontendURL := strings.TrimSpace(os.Getenv("FRONTEND_URL"))
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}
	parsedURL, err := url.Parse(frontendURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return errors.New("FRONTEND_URL must be an absolute URL")
	}

	passwordMailer = &smtpPasswordResetMailer{
		host:     values["SMTP_HOST"],
		port:     port,
		username: values["SMTP_USERNAME"],
		password: smtpPassword(values["SMTP_HOST"], values["SMTP_PASSWORD"]),
		from:     mail.Address{Name: values["SMTP_FROM_NAME"], Address: from.Address},
	}
	passwordResetFrontendURL = frontendURL
	return nil
}

// Google displays App Passwords in four groups for readability, but SMTP
// authentication uses the 16 characters without spaces.
func smtpPassword(host, password string) string {
	if strings.EqualFold(host, "smtp.gmail.com") {
		return strings.ReplaceAll(password, " ", "")
	}
	return password
}

func passwordResetURL(token string) (string, error) {
	parsedURL, err := url.Parse(passwordResetFrontendURL)
	if err != nil {
		return "", err
	}
	parsedURL.Path = strings.TrimRight(parsedURL.Path, "/") + "/reset-password"
	query := parsedURL.Query()
	query.Set("token", token)
	parsedURL.RawQuery = query.Encode()
	return parsedURL.String(), nil
}

func (m *smtpPasswordResetMailer) SendPasswordReset(_ context.Context, recipient, resetURL string) error {
	if _, err := mail.ParseAddress(recipient); err != nil {
		return fmt.Errorf("invalid recipient: %w", err)
	}
	message := buildPasswordResetEmail(m.from, recipient, resetURL)
	address := net.JoinHostPort(m.host, strconv.Itoa(m.port))
	connection, err := smtp.Dial(address)
	if err != nil {
		return fmt.Errorf("connect to SMTP server: %w", err)
	}
	defer connection.Quit()
	if ok, _ := connection.Extension("STARTTLS"); !ok {
		return errors.New("SMTP server does not support STARTTLS")
	}
	if err := connection.StartTLS(&tls.Config{ServerName: m.host, MinVersion: tls.VersionTLS12}); err != nil {
		return fmt.Errorf("start SMTP TLS: %w", err)
	}
	if err := connection.Auth(smtp.PlainAuth("", m.username, m.password, m.host)); err != nil {
		return fmt.Errorf("authenticate SMTP sender: %w", err)
	}
	if err := connection.Mail(m.from.Address); err != nil {
		return fmt.Errorf("set SMTP sender: %w", err)
	}
	if err := connection.Rcpt(recipient); err != nil {
		return fmt.Errorf("set SMTP recipient: %w", err)
	}
	writer, err := connection.Data()
	if err != nil {
		return fmt.Errorf("create SMTP message: %w", err)
	}
	if _, err := writer.Write([]byte(message)); err != nil {
		return fmt.Errorf("write SMTP message: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("send SMTP message: %w", err)
	}
	return nil
}

func buildPasswordResetEmail(from mail.Address, recipient, resetURL string) string {
	escapedURL := html.EscapeString(resetURL)
	plainText := "We received a request to reset your Garago password.\r\n\r\n" +
		"Reset your password: " + resetURL + "\r\n\r\n" +
		"This link expires in 15 minutes. If you did not request it, you can ignore this email.\r\n"
	htmlBody := "<p>We received a request to reset your Garago password.</p>" +
		"<p><a href=\"" + escapedURL + "\">Reset your password</a></p>" +
		"<p>This link expires in 15 minutes. If you did not request it, you can ignore this email.</p>"
	boundary := "garago-password-reset"
	return strings.Join([]string{
		"From: " + from.String(),
		"To: " + recipient,
		"Subject: Reset your Garago password",
		"MIME-Version: 1.0",
		"Content-Type: multipart/alternative; boundary=\"" + boundary + "\"",
		"",
		"--" + boundary,
		"Content-Type: text/plain; charset=UTF-8",
		"",
		plainText,
		"--" + boundary,
		"Content-Type: text/html; charset=UTF-8",
		"",
		htmlBody,
		"--" + boundary + "--",
		"",
	}, "\r\n")
}

func logPasswordResetDeliveryFailure(userID int32) {
	log.Printf("password reset email delivery failed for user_id=%d", userID)
}
