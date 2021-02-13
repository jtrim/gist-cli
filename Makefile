.PHONY: build
build:
	go build -o gist .

.PHONY: install
install: build
	mv ./gist /usr/local/bin
