# Philharmonie Ticket Notifier
## Overview
This script send email notifications when tickets for a specified concert at the Philharmonie de Paris become available. It uses Selenium to automate browser interactions and BeautifulSoup to parse HTML content.
## Requirements
- go 1.25 
- github.com/joho/godotenv v1.5.1
- github.com/tebeka/selenium v0.9.9
## Installation
1. Clone the repository:
    ```bash
    git clone https://github.com/Lynux2142/philharmonie-ticket-notifier.git
    cd philharmonie-ticket-notifier
    ```
## Configuration
1. Create a `.env` file in the project root directory with the following content:
    ```env
    TICKET_ID=your_concert_id_here
    EMAIL_PASSWORD=gmail_app_password_here
    NOTIFY_EMAIL=recipient_email_here,another_email_here
    SELENIUM_URL=http://localhost:4444/wd/hub
    ```
- `TICKET_ID`: The ID of the concert you want to monitor.
- `EMAIL_PASSWORD`: The password for the email account used to send notifications (preferably an app password).
- `NOTIFY_EMAIL`: Comma-separated list of email addresses to notify.
- `SELENIUM_URL`: The URL of the Selenium WebDriver server.
- `WAIT_TIMEOUT`: Timeout in seconds for waiting for elements to load (default `5`).
## Usage
1. From the sources:
    - Start a Selenium WebDriver server (e.g., using Docker):
        ```bash
        docker run -d -p 4444:4444 --shm-size 2g --name selenium-firefox selenium/standalone-firefox
        ```
    - Install the required packages:
        ```bash
        go mod tidy
        ```
    - Run the script:
        ```bash
        go run ./main.go
        ```
2. From Docker:
    - Create a docker network:
        ```bash
        docker network create selenium-net
        ```
    - Start a Selenium WebDriver server in the network:
        ```bash
        docker run -d --net selenium-net -p 4444:4444 --shm-size 2g --name selenium-firefox selenium/standalone-firefox
        ```
    - Build the Docker image:
        ```bash
        docker build -t philharmonie-ticket-notifier .
        ```
    - Run the Docker container:
        ```bash
        docker run --net selenium-net --env-file .env philharmonie-ticket-notifier
        ```
