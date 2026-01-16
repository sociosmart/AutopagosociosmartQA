FROM python:3.13.1-alpine AS base
WORKDIR /usr/src/app

# Dependencies installation
RUN apk add --update &&\
  apk add --upgrade && \
  apk add tzdata


# Setting timezone
RUN ln -s /usr/share/zoneinfo/America/Mazatlan /etc/localtime

RUN python -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"

COPY ./requirements/base.txt ./requirements/base.txt

RUN pip install --upgrade pip
RUN pip install -r ./requirements/base.txt

FROM base AS dev
ENV PYTHONUNBUFFERED=1
ENV PYTHONDONTWRITEBYTECODE=1
ENV TZ=America/Mazatlan

WORKDIR /usr/src/app

COPY --from=base /opt/venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"

COPY ./requirements/dev.txt ./requirements/dev.txt
RUN pip install -r ./requirements/dev.txt

COPY . .

CMD ["python", "main.py"]

#RUNTIME
FROM base AS prod
ENV PYTHONUNBUFFERED=1
ENV PYTHONDONTWRITEBYTECODE=1
ENV TZ=America/Mazatlan
WORKDIR /usr/src/app

COPY --from=base /opt/venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"

COPY ./requirements/prod.txt ./requirements/prod.txt
RUN pip install -r ./requirements/prod.txt

COPY . .

CMD ["python", "main.py"]
