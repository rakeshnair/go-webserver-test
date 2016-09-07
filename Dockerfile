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

# Install Heka
ENV HEKA_FILE_NAME heka-0_10_0-linux-amd64
ENV HEKA_VERSION 0.10.0
ENV HEKA_DOWNLOAD_URL https://github.com/mozilla-services/heka/releases/download/v$HEKA_VERSION/$HEKA_FILE_NAME.tar.gz
ENV HEKA_MD5 89ff62fe2ccad3d462c9951de0c15e38

RUN cd /usr/local && \
    curl -LO $HEKA_DOWNLOAD_URL && \
    echo "$HEKA_MD5  $HEKA_FILE_NAME.tar.gz" | md5sum --check && \
    echo "$HEKA_FILE_NAME.tar.gz" | xargs tar -zxf && \
    mv $HEKA_FILE_NAME heka && \
    echo "$HEKA_FILE_NAME.tar.gz" | xargs rm -rf

WORKDIR $GOPATH/src/github.com/rakeshnair/go-streaming-app
CMD go run main.go
