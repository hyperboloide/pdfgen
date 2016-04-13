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
	xfonts-75dpi

ADD https://bitbucket.org/wkhtmltopdf/wkhtmltopdf/downloads/wkhtmltox-0.13.0-alpha-7b36694_linux-trusty-amd64.deb /tmp/
RUN dpkg -i /tmp/wkhtmltox-0.13.0-alpha-7b36694_linux-trusty-amd64.deb && \
	rm -fr /tmp/*

COPY pdfgen /usr/local/bin/
