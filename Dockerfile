FROM golang:alpine as builder

RUN apk update && apk upgrade && \
                apk add --no-cache git

RUN mkdir /app
WORKDIR /app

COPY . .

RUN go mod download
RUN GOOS=linux go build main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /app .

ENTRYPOINT ["/main"]