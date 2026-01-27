package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tebeka/selenium"
)

type WebDriver struct {
	driver selenium.WebDriver
	selenium_url string
}

func NewWebDriver(selenium_url string) (*WebDriver, error) {
	wd := &WebDriver{
		selenium_url: selenium_url,
	}
	caps := selenium.Capabilities{
		"browserName": "firefox",
	}
	driver, err := selenium.NewRemote(caps, wd.selenium_url)
	if err != nil {
		return nil, fmt.Errorf("connection to Selenium remote server failed: %w", err)
	}
	if err := driver.SetImplicitWaitTimeout(10 * time.Second); err != nil {
        return nil, fmt.Errorf("failed to set implicit wait: %w", err)
    }
	wd.driver = driver
	return wd, nil
}

func (wd *WebDriver) CheckTicketAvailability(website website, ticket_id string) (bool, error) {
	if err := wd.driver.Get(website.events_url); err != nil {
		return false, fmt.Errorf("Failed to load page: %w", err)
	}
	log.Printf("Page open: %s\n", website.events_url)

	prod, err := wd.driver.FindElement(selenium.ByID, fmt.Sprintf("prod_%s", ticket_id))
	if err != nil {
		return false, fmt.Errorf("Concert not found: %w", err)
	}
	log.Println("Concert found.")

	title, err := prod.FindElement(selenium.ByCSSSelector, ".title")
	if err != nil {
		log.Printf("Failed to find concert title: %s", err)
	} else {
		text, err := title.Text()
		if err != nil {
			log.Printf("Failed to get concert title text: %s", err)
		} else {
			log.Printf("Concert Title: %s", text)
		}
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
