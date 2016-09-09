FROM segment/base:v4

RUN apt-get install -y vim git

# Install Go
ENV GOPATH="/root/dev"
ENV EXEC_DIR $GOPATH/src/github.com/rakeshnair/go-streaming-app

RUN cd /usr/local && \
    curl -L# https://storage.googleapis.com/golang/go1.6.2.linux-amd64.tar.gz | \
    tar -zx && \
    cp -Rp /usr/local/go/bin/* /usr/local/bin

# Install packages reqd by the application
RUN go get github.com/tools/godep
RUN go get github.com/rakeshnair/go-streaming-app

COPY include/startup.sh $EXEC_DIR/

WORKDIR $EXEC_DIR
RUN chmod +x startup.sh
CMD ["bash", "-C", "startup.sh"]
