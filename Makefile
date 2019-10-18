BINARY_NAME=bot

all:gofmt
	go build -o ${BINARY_NAME} .

clean:
	go clean -i

gofmt:
	gofmt -w .
	go vet .

help:
	@echo "make - compile the source code"
	@echo "make clean - remove binary file and vim swp files"
	@echo "make gofmt - run gofmt and go vet"
	@echo "make ca - generate ca files"

