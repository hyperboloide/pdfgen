NAME = pdfgen
VERSION = 1.0.1
DOCKERID = hyperboloide

all: container

clean:
	rm -fr $(NAME)

fmt:
	go fmt ./...

test:
	go test ./...

container: clean
	GOOS=linux GOARCH=amd64 go build -a
	docker build -t $(DOCKERID)/$(NAME) .

push:
	docker push     $(DOCKERID)/$(NAME)

run:
	docker run \
		--name pdfgen \
		--rm \
		-it \
		-p 8888:8888 \
		--mount src=$(CURDIR)/templates/,target=/etc/pdfgen/templates,type=bind \
		$(DOCKERID)/$(NAME)

.PHONY: all clean fmt test container push run
