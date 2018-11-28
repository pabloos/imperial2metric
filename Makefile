.SILENT: build-cli build-lambda

### VARS
program=imperial2metric
cvgfile=coverage.out

travis: build-cli build-lambda test

#build the command line util
build-cli:
	mkdir -p bin
	docker build --target cli_builder -t cli_builder .
	docker run -it --rm --volume "$$PWD:/go/src/$(program)" cli_builder

#build the lambda handler and zip
build-lambda:
	mkdir -p lambda-deploy
	docker build --target lambda_builder -t lambda_builder .
	docker run -it --rm --volume "$$PWD:/go/src/$(program)" lambda_builder
	yes | rm go.sum go.mod

	# GOOS=linux go build -o handler cmd/lambda/*
	# zip handler.zip handler
	# rm handler

test:
	docker build --target tester -t tester .
	docker run -it --rm --volume "$$PWD:/go/src/$(program)" tester
