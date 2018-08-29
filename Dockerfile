FROM golang:1.11.0
# as build

# Install system dependencies
ENV GOPATH /go

RUN apt-get update -qq && \
    apt-get install -qq -y geoip-bin libgeoip-dev pkg-config build-essential && \
    mkdir -p /go/src/github.com/Qlean/silvia
WORKDIR /go/src/github.com/Qlean/silvia

# COPY .  ./
COPY Gopkg.* ./
COPY cmd/ ./cmd/
COPY silvia/ ./silvia/

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && chmod +x /usr/local/bin/dep && \
    dep ensure -vendor-only && \
    curl http://geolite.maxmind.com/download/geoip/database/GeoLiteCity.dat.gz | gunzip > GeoLiteCity.dat
RUN cd ./silvia && ln -s ../GeoLiteCity.dat GeoLiteCity.dat && go test
RUN go build  -o /app/silvia ./cmd \
    && chmod +x /app/silvia && \
    ln -s $(pwd)/GeoLiteCity.dat /app/GeoLiteCity.dat


# FROM scratch
# COPY --from=build /go/src/github.com/Qlean/silvia/bin/silvia /
CMD ["./bin/silvia"]

LABEL MAINTAINER=Qlean
