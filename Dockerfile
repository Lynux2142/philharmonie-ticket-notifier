FROM python:3.12-slim
WORKDIR /app
COPY . /app
RUN pip install .
ENTRYPOINT ["python", "main.py"]
