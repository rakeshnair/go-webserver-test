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
    rm -rf $HEKA_FILE_NAME.tar.gz

# Place the Heka config file
RUN mkdir /usr/local/etc/heka
COPY contents/heka/hekad.toml /usr/local/etc/heka

# Create directory to store Heka file dumps
RUN mkdir /var/log/heka

COPY contents/startup.sh $GOPATH/src/github.com/rakeshnair/go-streaming-app/

# Create volume for dumping log files
RUN mkdir /data && chmod 744 /data
VOLUME /data

WORKDIR $GOPATH/src/github.com/rakeshnair/go-streaming-app
RUN chmod +x startup.sh
CMD ["bash", "-C", "startup.sh"]
