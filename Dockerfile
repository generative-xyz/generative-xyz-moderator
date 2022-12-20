FROM amd64/golang:1.17.0-alpine as builder

RUN apk update && apk upgrade && \
    apk --update add git make gcc g++ gnutls gnutls-dev gnutls-c++

ARG ENV=dev
ARG NETRC_USER=user
ARG NETRC_TOKEN=123

ENV ENV=${ENV} \
    CGO_ENABLED=1

WORKDIR /app

COPY go.mod go.sum Makefile ./

RUN make init

RUN echo machine gitlab.com login ${NETRC_USER} password ${NETRC_TOKEN} > $HOME/.netrc

RUN cat  $HOME/.netrc

RUN go mod download

COPY  . .

RUN  go mod tidy -compat=1.17

RUN echo "✅ Build for Linux"; make build

# Distribution
FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    mkdir /app 

WORKDIR /app 

EXPOSE 9090

COPY --from=builder /app/backend-api /app