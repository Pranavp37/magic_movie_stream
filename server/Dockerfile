FROM golang:1.25.0-alpine

# install git and Air
RUN apk update && apk add --no-cache git
RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .


EXPOSE 8080


CMD ["air"]
