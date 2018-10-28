FROM golang:stretch

WORKDIR /go/src/imperial2metric

CMD export GO111MODULE=on && \
    go mod init && \
    GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o bin/imperial2metric.exe   main.go && \
    GOOS=linux   GOARCH=amd64 CGO_ENABLED=0 go build -o bin/imperial2metric.linux main.go && \
    rm go.mod go.sum