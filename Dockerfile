FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./

COPY go.sum ./

# install curl for healthcheck
RUN apk add curl

RUN go mod download

COPY *.go ./

RUN go build -o /main

EXPOSE 8080

CMD [ "/main" ]
