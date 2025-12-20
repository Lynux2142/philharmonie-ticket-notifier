package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type website struct {
	events_url string
	purchase_url string
}

var website_list = [...]website{
	{
		events_url: "https://billetterie.philharmoniedeparis.fr/list/events",
		purchase_url: "https://billetterie.philharmoniedeparis.fr/selection/event/date?productId=%s",
	},
	{
		events_url: "https://bourseauxbillets.philharmoniedeparis.fr/list/events",
		purchase_url: "https://bourseauxbillets.philharmoniedeparis.fr/selection/event/date?productId=%s",
	},
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	selenium_url := os.Getenv("SELENIUM_URL")
	ticket_id := os.Getenv("TICKET_ID")
	sender_email := "guillerotlucas@gmail.com"
	email_password := os.Getenv("EMAIL_PASSWORD")
	notify_email := os.Getenv("NOTIFY_EMAIL")

	driver, err := NewWebDriver(selenium_url)
	if err != nil {
		log.Printf("Failed to initialize WebDriver: %s", err)
		return
	}
	defer driver.Quit()

	smtp_client, err := NewSmtpClient(
		"smtp.gmail.com",
		sender_email,
		email_password,
	)
	if err != nil {
		log.Printf("Failed to initialize SMTP client: %s", err)
		return
	}
	defer smtp_client.Quit()

	for _, website := range website_list {
		available, err := driver.CheckTicketAvailability(website, ticket_id)
		if err != nil {
			log.Printf("Error checking ticket availability: %s", err)
			continue
		}
		if !available {
			log.Println("Concert sold out.")
			continue
		}
		log.Println("Concert available")

		msg := []byte("To: " + notify_email + "\r\n" +
			"Subject: Concert Ticket Available!\r\n" +
			"\r\n" +
			"The concert ticket you were waiting for is now available! Visit the following link to book your ticket:\r\n" +
			fmt.Sprintf(website.purchase_url, ticket_id) +
			"\r\n",
		)

		err = smtp_client.SendMail(
			strings.Split(notify_email, ","),
			msg,
		)
		if err != nil {
			log.Printf("Failed to send notification email: %s", err)
			return
		}
		log.Println("Notification email sent successfully.")
	}
}
