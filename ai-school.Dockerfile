FROM golang:1.18-bullseye as deps

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

FROM python:3.9.16-bullseye as python-env

RUN DEBIAN_FRONTEND=noninteractive \
    apt-get -y update \
    && apt-get install -y unzip

WORKDIR /app
EXPOSE 8001

ADD https://github.com/generative-xyz/perceptron-training/archive/refs/heads/main.zip /app/perceptron-training-main.zip
RUN unzip /app/perceptron-training-main.zip
RUN cd /app/perceptron-training-main \
    && pip3 install -r requirements.txt

FROM python-env as runner

WORKDIR /app
EXPOSE 8001

COPY --from=builder /app/backend-api /app

CMD ["/app/backend-api"]
