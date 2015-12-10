FROM debian:jessie
MAINTAINER Frederic Delbos <fred@hyperboloide.com>

ADD http://download.gna.org/wkhtmltopdf/0.12/0.12.2.1/wkhtmltox-0.12.2.1_linux-jessie-amd64.deb /tmp/
RUN apt-get -qq update && \
	apt-get -qqy install \
    libssl1.0.0 \
	fontconfig \
	libfontconfig1 \
	libfreetype6 \
	libjpeg62-turbo \
	libx11-6 \
	libxext6 \
	libxrender1 \
	xfonts-base \
	xfonts-75dpi && \
	dpkg -i /tmp/wkhtmltox-0.12.2.1_linux-jessie-amd64.deb

RUN mkdir /usr/local/templates
ENV PDFGEN_TEMPLATES /usr/local/templates


COPY pdfgen /usr/local/bin/
