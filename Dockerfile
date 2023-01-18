FROM golang:1.18 as builder

RUN apt-get -y update && apt-get -y upgrade && \
    apt-get -y install git && \
    apt-get -y install make 

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
FROM ubuntu:20.04
RUN apt-get -y update && apt upgrade -y
RUN apt-get -y install  software-properties-common && \
    apt-get -y install wget && \
    DEBIAN_FRONTEND=noninteractive TZ=Etc/UTC apt-get -y install tzdata


RUN wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
RUN apt-get -y install ./google-chrome-stable_current_amd64.deb
RUN apt-get -y install xvfb xorg unzip dbus-x11 
RUN apt-get -y install unzip --fix-missing

RUN wget -N http://chromedriver.storage.googleapis.com/2.25/chromedriver_linux64.zip
RUN unzip chromedriver_linux64.zip
RUN chmod +x chromedriver
RUN mv -f chromedriver /usr/local/bin/chromedriver

ENV DISPLAY=:99
ENV XVFB_WHD=1280x720x16
WORKDIR /app 

COPY --from=builder /app/backend-api /app
COPY --from=builder /app/entrypoint.sh /app

RUN chmod +x /app/entrypoint.sh
CMD ["./app/entrypoint.sh"]