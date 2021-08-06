FROM golang:alpine as builder

LABEL "com.github.actions.name"="Notion Card Updater"
LABEL "com.github.actions.description"="Updates a Notion card based on events and inputs using the Notion API"
LABEL "com.github.actions.icon"="align-justify"
LABEL "com.github.actions.color"="blue"

LABEL "repository"="https://github.com/zant/notion-cards-action/"
LABEL "homepage"="https://github.com/zant/notion-cards-action/README.md"
LABEL "maintainer"="zant <yo@zant.xyz>"

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

ENTRYPOINT ["/app/main"]