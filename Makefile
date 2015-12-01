NAME = pdfgen
VERSION = 1.0.1

all:
	go build -v

clean:
	rm -fr $(NAME)

fmt:
	go fmt ./...

test:
	go test ./...

.PHONY: all clean fmt test
