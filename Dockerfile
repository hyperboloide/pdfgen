FROM ubuntu
MAINTAINER Frederic Delbos <fred.delbos@gmail.com>

RUN apt-get -qq update && \
	apt-get -qqy install \
  libssl1.0.0 \
	fontconfig \
	libfontconfig1 \
	libfreetype6 \
	libjpeg-turbo8 \
	libx11-6 \
	libxext6 \
	libxrender1 \
	xfonts-base \
	xfonts-75dpi \
	xz-utils

RUN mkdir /tmp/wkhtmltox
WORKDIR /tmp/wkhtmltox

ADD https://github.com/wkhtmltopdf/wkhtmltopdf/releases/download/0.12.4/wkhtmltox-0.12.4_linux-generic-amd64.tar.xz .
RUN tar xf wkhtmltox-0.12.4_linux-generic-amd64.tar.xz &&\
	mv wkhtmltox/bin/wkhtmltopdf /usr/local/bin/

WORKDIR /
RUN rm -fr /tmp/wkhtmltox
COPY pdfgen /usr/local/bin/pdfgen

EXPOSE 8888
ENV PDFGEN_PORT=8888
ENV PDFGEN_ADDR=0.0.0.0

CMD ["pdfgen"]
