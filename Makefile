.SILENT: build

program = imperial2metric

build:
	mkdir -p bin
	docker build -t builder .
	docker run -it --rm --volume "$$PWD:/go/src/imperial2metric" builder

clean:
	rm -R bin

test:
	go test -v ./...

lambda-ready:
	GOOS=linux go build -o handler cmd/lambda/zipper.go
	zip handler.zip handler
