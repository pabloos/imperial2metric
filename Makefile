.SILENT: build

build: clean
	mkdir -p bin
	docker build -t builder .
	docker run -it --rm --volume "$$PWD:/go/src/imperial2metric" builder

clean:
	rm -R bin