NAME = pdfgen
VERSION = 1.0.1
DOCKERID = hyperboloide


all:
	go build -v

clean:
	rm -fr $(NAME)

fmt:
	go fmt ./...

test:
	go test ./...

container: clean
	GOOS=linux GOARCH=amd64 go build -a
	docker build -t $(DOCKERID)/$(NAME) .

dockerhub: container
	docker push     $(DOCKERID)/$(NAME)

.PHONY: all clean fmt test
