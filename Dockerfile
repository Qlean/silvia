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
RUN curl http://geolite.maxmind.com/download/geoip/database/GeoLiteCity.dat.gz | gunzip > GeoLiteCity.dat
RUN cd /app/silvia && go test
RUN go build

CMD ["./app"]
