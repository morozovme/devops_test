FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./

COPY go.sum ./

# install curl for healthcheck
RUN apk add curl

# need for deps fetch
RUN apk add git

RUN go get github.com/prometheus/client_golang/prometheus/promhttp

RUN go mod download

COPY *.go ./

RUN go build -o /main

EXPOSE 8080

CMD [ "/main" ]
