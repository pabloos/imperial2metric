FROM golang:stretch as gopher

ENV GO111MODULE=on

WORKDIR /go/src/imperial2metric/