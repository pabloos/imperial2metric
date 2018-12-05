
PROGRAM_NAME=${PWD##*/} 

echo $PROGRAM_NAME

GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -mod=vendor -o bin/imperial2metric.exe cmd/shell/main.go && \
GOOS=linux   GOARCH=amd64 CGO_ENABLED=0 go build -mod=vendor -o bin/imperial2metric.linux cmd/shell/main.go && \
GOOS=darwin  GOARCH=amd64 CGO_ENABLED=0 go build -mod=vendor -o bin/imperial2metric.darwin cmd/shell/main.go;