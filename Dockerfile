# first image used to build the sources
FROM golang:1.19-buster AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY . .
RUN go mod download

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o bin/feederd cmd/feederd/main.go


# Second image, running the tdexd executable
FROM debian:buster-slim

# $USER name, and data $DIR to be used in the `final` image
ARG USER=feeder
ARG DIR=/home/feeder

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends ca-certificates

# NOTE: Default GID == UID == 1000
RUN adduser --disabled-password \
            --home "$DIR/" \
            --gecos "" \
            "$USER"
USER $USER

COPY --from=builder /app/bin/* /usr/local/bin/

# Prevents `VOLUME $DIR/.tdex-feeder/` being created as owned by `root`
RUN mkdir -p "$DIR/.tdex-feeder/"

# Expose volume containing all data
VOLUME $DIR/.tdex-feeder/



ENTRYPOINT [ "feederd" ]