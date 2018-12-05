.SILENT: build-cli build-lambda

############
### VARS ###
############
program=imperial2metric
docker=docker run -it --rm --volume "$$PWD:/go/src/$(program)" gopher /bin/bash -c
lmbdir=lambda-deploy

#this needs to be the first rule in the makefile because will be the one invoked by travis CI
travis: build-cli build-lambda test

builder:
	docker build --target gopher -t gopher .

#build the command line util
build-cli: builder
	mkdir -p bin
	$(docker) ./x-compile.sh

#build the lambda handler and zip
build-lambda: builder
	yes | rm -R lambda-deploy || true
	mkdir -p lambda-deploy
	$(docker) "go mod vendor && GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -mod=vendor -race -o lambda-deploy/handler imperial2metric/cmd/lambda"
	chmod 755 lambda-deploy/handler

deploy-lambda: build-lambda
	cd $(lmbdir) && zip handler.zip handler && mv handler.zip ../

test: builder
	$(docker) "go test -v -mod=vendor -coverpkg=./... -coverprofile=coverage.out ./... && go tool cover -func=coverage.out && rm coverage.out"

clean-vendor:
	yes | rm -R vendor