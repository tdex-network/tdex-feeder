FROM golang:1.15.5-buster AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /tdex-feeder

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o feederd-linux cmd/feederd/main.go

WORKDIR /build

RUN cp /tdex-feeder/feederd-linux .

FROM debian:buster

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

COPY --from=builder /build/ /

CMD ["/feederd-linux"]