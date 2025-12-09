from os import getenv
from dotenv import load_dotenv
import smtplib
from selenium import webdriver
from selenium.common.exceptions import TimeoutException
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.by import By
from bs4 import BeautifulSoup
from email.message import EmailMessage
import logging

load_dotenv()

URLs = [
    "https://billetterie.philharmoniedeparis.fr/list/events",
    "https://bourseauxbillets.philharmoniedeparis.fr/list/events",
]
EMAIL_PASSWORD = getenv("EMAIL_PASSWORD")
TICKET_ID = getenv("TICKET_ID")
NOTIFY_EMAIL = getenv("NOTIFY_EMAIL")
REMOTE_WEB_DRIVER = getenv("REMOTE_WEB_DRIVER", "false").lower() == "true"
SELENIUM_URL = getenv("SELENIUM_URL")
WAIT_ELEMENT_TIMEOUT = int(getenv("WAIT_ELEMENT_TIMEOUT", "5"))

logger = logging.getLogger(__name__)
formatter = logging.Formatter(
    fmt="%(asctime)s [%(levelname)-8s] - %(message)s",
    datefmt="%Y-%m-%d %H:%M:%S",
)
stream_handler = logging.StreamHandler()
stream_handler.setFormatter(formatter)
logger.setLevel(logging.INFO)
logger.addHandler(stream_handler)
logger.propagate = False

def send_email(subject, body, to_email):
    msg = EmailMessage()
    msg["Subject"] = subject
    msg["From"] = "guillerotlucas@gmail.com"
    msg["To"] = to_email
    msg.set_content(body)

    # --- SMTP configuration ---
    smtp_server = "smtp.gmail.com"
    smtp_port = 587
    smtp_user = "guillerotlucas@gmail.com"
    smtp_password = EMAIL_PASSWORD

    # --- send email ---
    with smtplib.SMTP(smtp_server, smtp_port) as server:
        server.starttls()
        server.login(smtp_user, smtp_password)
        server.send_message(msg)

def main():
    ticket_id = TICKET_ID
    options = webdriver.FirefoxOptions()
    options.add_argument("--headless")
    logger.info("Connection to Firefox Webdriver")
    if REMOTE_WEB_DRIVER:
        driver = webdriver.Remote(
            command_executor=SELENIUM_URL,
            options=options,
        )
    else:
        driver = webdriver.Firefox(options=options)
    logger.info("Firefox Webdriver initialized")
    is_available = False
    for url in URLs:
        logger.info("Loading page {}".format(url))
        driver.get(url)
        logger.info("Page loaded")
        try:
            logger.info("Waiting for element to be present ...")
            element = EC.presence_of_element_located(
                (By.ID, f"prod_{ticket_id}")
            )
            WebDriverWait(driver, WAIT_ELEMENT_TIMEOUT).until(element)
            logger.info("Element found")
        except TimeoutException:
            logger.error("Loading took too much time!")

        logger.info("Parsing page source ...")
        soup = BeautifulSoup(driver.page_source, "html.parser")
        logger.info("Checking for concert with id {}".format(ticket_id))
        element = soup.find(id=f"prod_{ticket_id}")
        title = element.find(class_="title").get_text(strip=True)
        logger.info("Concert found: {}".format(title))
        is_available = element.find(id=f"product_{ticket_id}_book") is not None
        logger.info(f"Ticket status: {'Available' if is_available else 'Unavailable'}")
        if is_available:
            logger.info("Sending notification email to: {}".format(NOTIFY_EMAIL))
            send_email(
                subject="Tickets Available!",
                body=f"Tickets are now available at https://billetterie.philharmoniedeparis.fr/selection/event/date?productId={ticket_id}",
                to_email=NOTIFY_EMAIL,
            )
            logger.info("Notification email sent")
            break


if __name__ == "__main__":
    main()
