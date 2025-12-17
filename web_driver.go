package main

import (
	"fmt"
	"log"

	"github.com/tebeka/selenium"
)

type WebDriver struct {
	driver selenium.WebDriver
	selenium_url string
}

func (wd *WebDriver) Init() error {
	caps := selenium.Capabilities{
		"browserName": "firefox",
	}
	driver, err := selenium.NewRemote(caps, wd.selenium_url)
	if err != nil {
		return fmt.Errorf("connection to Selenium remote server failed: %w", err)
	}
	wd.driver = driver
	return nil
}

func (wd *WebDriver) CheckTicketAvailability(website website, ticket_id string) (bool, error) {
	if err := wd.driver.Get(website.events_url); err != nil {
		log.Printf("Failed to load page: %s", err)
		return false, err
	}
	log.Printf("Page open: %s\n", website.events_url)

	prod, err := wd.driver.FindElement(selenium.ByID, fmt.Sprintf("prod_%s", ticket_id))
	if err != nil {
		log.Printf("Concert not found: %s", err)
		return false, err
	}
	log.Println("Concert found.")

	title, err := prod.FindElement(selenium.ByCSSSelector, ".title")
	if err != nil {
		log.Printf("Failed to find concert title: %s", err)
	} else {
		text, err := title.Text()
		if err != nil {
			log.Printf("Failed to get concert title text: %s", err)
			return false, err
		}
		log.Printf("Concert Title: %s", text)
	}

	_, err = prod.FindElement(selenium.ByID, fmt.Sprintf("product_%s_book", ticket_id))
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (wd *WebDriver) Quit() {
	wd.driver.Quit()
}
