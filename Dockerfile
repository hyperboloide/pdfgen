FROM ubuntu:trusty
MAINTAINER Frederic Delbos <fred@hyperboloide.com>

RUN apt-get -qq update && \
	apt-get -qqy install \
    libssl1.0.0 \
	fontconfig \
	libfontconfig1 \
	libfreetype6 \
	libjpeg-turbo8 \
	libicu52 \
	libx11-6 \
	libxext6 \
	libxrender1 \
	xfonts-base \
	xfonts-75dpi \
	xz-utils \
	curl


WORKDIR /tmp/

ADD http://download.gna.org/wkhtmltopdf/0.12/0.12.3/wkhtmltox-0.12.3_linux-generic-amd64.tar.xz .

RUN tar xf wkhtmltox-0.12.3_linux-generic-amd64.tar.xz &&\
	mv wkhtmltox/bin/wkhtmltopdf /usr/local/bin/ && \
	rm -fr /tmp

COPY pdfgen /pdfgen

EXPOSE 8888
ENV PDFGEN_PORT=8888
ENV PDFGEN_ADDR=0.0.0.0

CMD ["/pdfgen"]
