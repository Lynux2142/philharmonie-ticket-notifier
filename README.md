# Philharmonie Ticket Notifier
# this project is a data extraction tool for the Philharmonie de Paris website to be notified by email if tickets for a concert are available.
## Overview
This project is designed to help users monitor ticket availability for concerts at the Philharmonie de Paris. By scraping the website for specific concert information, users will receive email notifications when tickets become available.
## Features
- Web scraping to check ticket availability
- Email notifications for ticket availability
- Configurable concert ID and email settings
## Requirements
- Python >= 3.12
- beautifulsoup4 >= 4.12.2
- python-dotenv >= 1.2.1
- selenium >= 4.38.0
## Installation
1. Clone the repository:
    ```bash
    git clone https://github.com/Lynux2142/philharmonie-ticket-notifier.git
    cd philharmonie-ticket-notifier
    ```
2. Install the required packages:
    ```bash
    uv sync
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
## Usage
1. From the python script:
    - Start a Selenium WebDriver server (e.g., using Docker):
        ```bash
        docker run -d -p 4444:4444 --name selenium-firefox selenium/standalone-chrome
        ```
    - Run the script:
        ```bash
        uv run main.py
        ```
2. From Docker:
    - Create a docker network:
        ```bash
        docker network create selenium-net
        ```
    - Start a Selenium WebDriver server in the network:
        ```bash
        docker run -d --net selenium-net -p 4444:4444 --name selenium-firefox selenium/standalone-chrome
        ```
    - Build the Docker image:
        ```bash
        docker build -t philharmonie-ticket-notifier .
        ```
    - Run the Docker container:
        ```bash
        docker run --net selenium-net --env-file .env philharmonie-ticket-notifier
        ```
