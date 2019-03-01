FROM golang:alpine as builder

RUN  mkdir /build

ADD ./src /build/

WORKDIR /build

RUN apk add --no-cache git

RUN go get -d -v ./...

RUN go build -o server .

FROM alpine

RUN mkdir /snowstorm

RUN adduser -S -D -H -h /snowstorm stormy

USER stormy

COPY --from=builder /build/server /snowstorm

WORKDIR /snowstorm

ENTRYPOINT ["./server"]
