FROM ubuntu:16.04
MAINTAINER rakesh@segment.com

ENV TERM="xterm"
RUN apt-get update --fix-missing && apt-get install -y \
  curl \
  sudo \
  git \
  vim

# Install Go
ENV GOPATH="/root/dev"
RUN cd /usr/local && \
  curl -L# https://storage.googleapis.com/golang/go1.6.2.linux-amd64.tar.gz | \
  tar -zx && \
  cp -Rp /usr/local/go/bin/* /usr/local/bin

# Install packages reqd by the application
RUN go get github.com/tools/godep
RUN go get github.com/rakeshnair/go-streaming-app

WORKDIR /root/dev/src/github.com/rakeshnair/go-streaming-app
CMD go run streaming.go
