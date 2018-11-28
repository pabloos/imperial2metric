FROM golang:stretch as cli_builder

WORKDIR /go/src/imperial2metric/

CMD export GO111MODULE=on && \
    go mod init && \
    GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o bin/imperial2metric.exe   cmd/shell/main.go && \
    GOOS=linux   GOARCH=amd64 CGO_ENABLED=0 go build -o bin/imperial2metric.linux cmd/shell/main.go && \
    GOOS=darwin  GOARCH=amd64 CGO_ENABLED=0 go build -o bin/imperial2metric.darwin cmd/shell/main.go && \
    rm go.mod go.sum

FROM golang:stretch as lambda_builder

WORKDIR /go/src/imperial2metric/

CMD export GO111MODULE=on && \
    go mod init && \
    GOOS=linux   GOARCH=amd64 CGO_ENABLED=0 go build -o lambda-deploy/handler cmd/lambda/*