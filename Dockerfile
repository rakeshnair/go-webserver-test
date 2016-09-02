FROM ubuntu:16.04
MAINTAINER rakesh@segment.com

ENV TERM="xterm"
RUN apt-get update --fix-missing && apt-get install -y \
  curl \
  logrotate \
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
RUN go get "gopkg.in/natefinch/lumberjack.v2" 

RUN mkdir -p /var/segment/log
RUN mkdir -p /var/segment/app

ADD contents/ /var/segment/app/
RUN ls -la /var/segment/app/*

# Place the logrotate file in the correct location
COPY contents/main_app_logrotate /etc/logrotate.d/
RUN chmod +x /var/segment/app/main_app_logrotate.sh

# Place the custom version for logrotate.conf. Reqd. since we dont have syslog
RUN rm /etc/logrotate.conf
COPY contents/logrotate.conf /etc/logrotate.conf

# Place custom crontab for rotating main app logs
ADD contents/main_app_crontab /etc/cron.d/main_app_crontab
RUN chmod 0644 /etc/cron.d/main_app_crontab

# Create cron log file
RUN touch /var/log/cron.log

WORKDIR /var/segment/app
CMD cron && go run main.go
