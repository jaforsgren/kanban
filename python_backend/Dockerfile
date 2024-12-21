FROM python:3.10-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .

# Gunicorn as the production server
CMD ["gunicorn", "-w", "4", "-b", "0.0.0.0:9082", "run:app"]