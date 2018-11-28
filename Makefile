.SILENT: build-cli build-lambda

### VARS
program=imperial2metric
cvgfile=coverage.out

travis: build-cli build-lambda test

#build the command line util
build-cli:
	mkdir -p bin
	docker build -t builder .
	docker run -it --rm --volume "$$PWD:/go/src/$(program)" builder

#build the lambda handler and zip
build-lambda:
	GOOS=linux go build -o handler cmd/lambda/*
	zip handler.zip handler
	rm handler

clean:
	yes | rm -R bin

test:
	go test -v ./... -coverprofile=$(cvgfile) && go tool cover -func=$(cvgfile) && rm $(cvgfile)