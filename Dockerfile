FROM golang:1.24.1-alpine

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

ARG ENVIRONMENT=${ENVIRONMENT}
ARG TARGETARCH

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

ENTRYPOINT []
