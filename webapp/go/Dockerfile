FROM golang:1.14

EXPOSE 1323

WORKDIR /go/src/isuumo

RUN apt-get update && apt-get install -y wget default-mysql-client

ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download
COPY . .
RUN make isuumo
