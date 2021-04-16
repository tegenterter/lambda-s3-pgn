FROM golang:1-alpine3.12 AS build

ARG PGN_EXTRACT_DOWNLOAD_URL

RUN apk add --no-cache wget unzip make gcc musl-dev
RUN wget ${PGN_EXTRACT_DOWNLOAD_URL} -O /var/tmp/pgn-extract.zip
RUN cd /var/tmp && unzip pgn-extract.zip
RUN make -C /var/tmp/pgn-extract

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN go build -tags lambda.norpc -ldflags="-s -w" main.go

FROM alpine:3.12

WORKDIR /app

COPY --from=build /app/main ./
COPY --from=build /var/tmp/pgn-extract/pgn-extract /bin

ENTRYPOINT ["/app/main"]
