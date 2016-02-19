FROM golang:1.6.0-wheezy
MAINTAINER Qlean

# Install system dependencies
RUN apt-get update -qq && \
    apt-get install -qq -y geoip-bin libgeoip-dev pkg-config build-essential

RUN mkdir -p /app
WORKDIR /app
COPY . /app/
ENV GOPATH /go/
RUN go get -d -v
RUN go build
EXPOSE 8080

CMD ["./app"]
