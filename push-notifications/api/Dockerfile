FROM python:3.12-bullseye

# Set up environment variables for Python
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
ENV TZ=UTC

WORKDIR /usr/src/api

COPY ./requirements.txt .

RUN pip install -r requirements.txt
RUN pip install waitress

COPY . .

EXPOSE 5000

CMD ["waitress-serve", "--host", "0.0.0.0", "--port", "5000", "main:app"]
