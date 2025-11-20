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

URL = "https://billetterie.philharmoniedeparis.fr/content#"
EMAIL_PASSWORD = getenv("EMAIL_PASSWORD")
TICKET_ID = getenv("TICKET_ID")
NOTIFY_EMAIL = getenv("NOTIFY_EMAIL")
SELENIUM_URL = getenv("SELENIUM_URL")

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
    logger.info("Connection to Firefox Webdriver")
    driver = webdriver.Remote(
        command_executor=SELENIUM_URL,
        options=webdriver.FirefoxOptions(),
    )
    logger.info("Firefox Webdriver initialized")
    logger.info("Loading page ...")
    driver.get(URL)
    logger.info("Page loaded")
    try:
        logger.info("Waiting for element to be present ...")
        element = EC.presence_of_element_located(
            (By.CLASS_NAME, f"stx-ProductCard-{ticket_id}")
        )
        WebDriverWait(driver, 5).until(element)
        logger.info("Element found")
    except TimeoutException:
        print("Loading took too much time!")

    logger.info("Parsing page source ...")
    soup = BeautifulSoup(driver.page_source, "html.parser")
    logger.info("Page source parsed")
    logger.info("Checking for ticket availability ...")
    element = soup.find(class_=f"stx-ProductCard-{ticket_id}")
    sold_out_status = element.find(id=ticket_id).find("span").text.strip() == "Sold Out"
    logger.info(f"Ticket status: {'Available' if not sold_out_status else 'Unavailable'}")
    if not sold_out_status:
        logger.info("Sending notification email ...")
        send_email(
            subject="Tickets Available!",
            body=f"Tickets are now available at https://billetterie.philharmoniedeparis.fr/selection/event/date?productId={ticket_id}",
            to_email=NOTIFY_EMAIL,
        )
        logger.info("Notification email sent")


if __name__ == "__main__":
    main()
