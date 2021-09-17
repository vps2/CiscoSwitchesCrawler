.PHONY: build clean help

GOOS = $(shell go env GOOS)

## build: создать исполняемый файл в каталоге .bin/
build:
ifeq ($(GOOS),windows)
	go build -o .bin/cisco_crawler.exe -ldflags "-s -w" cmd/cisco_crawler/main.go
else
	go build -o .bin/cisco_crawler -ldflags "-s -w" cmd/cisco_crawler/main.go
endif

## clean: удалить содержимое папки .bin/
clean:
	rm -f .bin/*

help: Makefile
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
