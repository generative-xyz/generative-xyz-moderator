FROM golang:1.21.0-bullseye as deps

RUN apt-get -y update && apt-get -y upgrade && \
    apt-get -y install git && \
    apt-get -y install make

ARG ENV=dev

ENV ENV=${ENV} \
    CGO_ENABLED=1

WORKDIR /app

COPY go.mod go.sum Makefile ./

RUN make init

RUN go mod download

FROM deps as builder
COPY  . .

RUN echo "✅ Build for Linux"; make build

# Distribution
FROM debian:bullseye as runner
RUN apt-get -y update && apt upgrade -y
RUN apt-get -y install  software-properties-common && \
    apt-get -y install wget && \
    DEBIAN_FRONTEND=noninteractive TZ=Etc/UTC apt-get -y install tzdata

 RUN wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
 RUN apt-get -y install ./google-chrome-stable_current_amd64.deb

WORKDIR /app

COPY --from=builder /app/backend-api /app
