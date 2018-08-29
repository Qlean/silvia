FROM golang:1.11.0
# as build

# Install system dependencies

ARG SRC=/go/src/github.com/Qlean/silvia
ENV GOPATH /go
WORKDIR /app

RUN apt-get update -qq && \
    apt-get install -qq -y geoip-bin libgeoip-dev pkg-config build-essential

COPY Gopkg.* $SRC/
COPY cmd/ $SRC/cmd/
COPY silvia/ $SRC/silvia/

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && \
    chmod +x /usr/local/bin/dep && \
    cd $SRC && \
    dep ensure -vendor-only && rm /usr/local/bin/dep &&\
    curl http://geolite.maxmind.com/download/geoip/database/GeoLiteCity.dat.gz | gunzip > /app/GeoLiteCity.dat && \
    cd ./silvia && \
    ln -s /app/GeoLiteCity.dat GeoLiteCity.dat && \
    go test


RUN cd $SRC && \
    go build  -o /app/silvia ./cmd && \
    chmod +x /app/silvia && \
    rm -rf $GOPATH

# # FROM scratch
# # COPY --from=build /go/src/github.com/Qlean/silvia/bin/silvia /
CMD ["/app/silvia"]

LABEL MAINTAINER=Qlean
