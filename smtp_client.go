package main

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

type SmtpClient struct {
	client *smtp.Client
	smtp_server string
	sender string
	password string
}

func (c *SmtpClient) Init() error {
	var err error
	c.client, err = smtp.Dial(c.smtp_server + ":587")
	if err != nil {
		return fmt.Errorf("failed to dial SMTP server: %w", err)
	}
	
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         c.smtp_server,
	}
	if err = c.client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}
	
	auth := smtp.PlainAuth("", c.sender, c.password, c.smtp_server)
	if err = c.client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate to SMTP server: %w", err)
	}
	return nil
}

func (c *SmtpClient) SendMail(to []string, msg []byte) error {
	if err := c.client.Mail(c.sender); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	for _, addr := range to {
		if err := c.client.Rcpt(addr); err != nil {
			return fmt.Errorf("failed to set recipient: %w", err)
		}
	}
	w, err := c.client.Data()
	if err != nil {
		return fmt.Errorf("failed to get SMTP data writer: %w", err)
	}
	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("failed to write email message: %w", err)
	}
	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close SMTP data writer: %w", err)
	}
	return nil
}

func (c *SmtpClient) Quit () { 
	c.client.Quit()
}
