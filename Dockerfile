FROM golang:1.21-alpine3.17

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN rm -rf /var/cache/apk/* && rm -rf /tmp/*
RUN apk --no-cache add curl

RUN go install github.com/cosmtrek/air@latest
RUN go mod download

COPY ./ ./
COPY *.json ./

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]
