# Base Image
FROM golang:1.23-alpine AS base

WORKDIR /usr/src/backend

ENV GOOS linux
ENV CGO_ENABLED 0

# Go binaries
RUN go install github.com/air-verse/air@latest && \
    go install github.com/google/wire/cmd/wire@latest && \
    go install github.com/swaggo/swag/cmd/swag@latest && \
    go install github.com/spf13/cobra-cli@latest && \
    go install github.com/vektra/mockery/v2@v2.20.0

RUN apk update && apk upgrade && apk add \
                                        curl \
                                        tzdata


# Builder
FROM base AS builder

COPY go.mod .
COPY go.sum .

COPY . .

RUN go mod download

RUN go build -o /go/bin/app .

# Runtime
FROM base AS runtime

COPY --from=builder /go/bin/app /go/bin/app

CMD ["/go/bin/app", "serve"]
