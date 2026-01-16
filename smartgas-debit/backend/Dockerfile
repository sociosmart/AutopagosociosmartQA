FROM python:3.13-bullseye

# Set up environment variables for Python
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
ENV TZ=UTC

WORKDIR /usr/src/backend

COPY ./requirements.txt .

RUN pip install -r requirements.txt
RUN pip install waitress

COPY . .

EXPOSE 8005

CMD ["waitress-serve", "--host", "0.0.0.0", "--port", "8005", "--call",  "app:create_app"]
