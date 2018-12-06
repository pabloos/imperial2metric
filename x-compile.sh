#!/bin/bash

PROGRAM_NAME=${PWD##*/} 

OSs=( windows linux darwin )

declare -A extensions=(
    [windows]=.exe
    [linux]=.linux
    [darwin]=.osx
)

go mod vendor

for os in "${OSs[@]}" ; do
    GOOS=$os GOARCH=amd64 CGO_ENABLED=0 go build -mod=vendor -o bin/$PROGRAM_NAME${extensions[$os]} cmd/shell/main.go
done