
PROGRAM_NAME=${PWD##*/} 

go mod vendor
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -mod=vendor -o bin/$PROGRAM_NAME.exe cmd/shell/main.go && \
GOOS=linux   GOARCH=amd64 CGO_ENABLED=0 go build -mod=vendor -o bin/$PROGRAM_NAME.linux cmd/shell/main.go && \
GOOS=darwin  GOARCH=amd64 CGO_ENABLED=0 go build -mod=vendor -o bin/$PROGRAM_NAME.darwin cmd/shell/main.go;