FROM golang:1.12.5 as builder
WORKDIR /go/src/dynamic-http-gateway
COPY . .
RUN CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" -o dynamic-http-gateway

# Create a minimal container to run a Golang static binary
FROM ubuntu:14.04

LABEL authors="dequan.ma@56qq.com"

RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN apt-get update \
  && apt-get install -y wget \
  && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /go/src/dynamic-http-gateway/dynamic-http-gateway .
COPY scripts/start.sh /start.sh

ENTRYPOINT ["/start.sh"]
EXPOSE 8080 8081
