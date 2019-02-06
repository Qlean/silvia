FROM golang:1.11.0
# as build
ARG BUILD_SOURCE
ARG REPO_NAME
ARG REPO_OWNER

LABEL org.opencontainers.image.vendor=${REPO_OWNER} \
      org.opencontainers.image.title=${REPO_NAME} \
      org.opencontainers.image.source=${BUILD_SOURCE}

# Install system dependencies

ARG SRC=/go/src/github.com/Qlean/silvia
ENV GOPATH /go
WORKDIR /app

RUN apt-get update -qq && \
    apt-get install -qq -y geoip-bin libgeoip-dev pkg-config build-essential

COPY GeoLiteCity.dat /app/
COPY Gopkg.* $SRC/
COPY cmd/ $SRC/cmd/
COPY silvia/ $SRC/silvia/

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && \
    chmod +x /usr/local/bin/dep && \
    cd $SRC && \
    dep ensure -vendor-only && rm /usr/local/bin/dep &&\
    cd ./silvia && \
    ln -s /app/GeoLiteCity.dat GeoLiteCity.dat && \
    cd $SRC && \
    go build  -o /app/silvia ./cmd && \
    chmod +x /app/silvia

# # FROM scratch
# # COPY --from=build /go/src/github.com/Qlean/silvia/bin/silvia /
CMD ["/app/silvia"]

LABEL MAINTAINER=Qlean
