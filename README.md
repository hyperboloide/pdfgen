# pdfgen
[![Build Status](https://travis-ci.org/hyperboloide/pdfgen.svg?branch=master)](https://travis-ci.org/hyperboloide/pdfgen)
[![Go Report Card](https://goreportcard.com/badge/github.com/hyperboloide/pdfgen)](https://goreportcard.com/report/github.com/hyperboloide/pdfgen)


HTTP service to generate PDF from Json requests

## Install and run

The recommended method is to use the docker container by mounting your template
directory (here with the provided example `template` directory):

```
docker run --rm -it -p 8888:8888 \
  --mount src=my_templates/,target=/etc/pdfgen/templates,type=bind \
  hyperboloide/pdfgen
```

If you rather not using Docker, you need to install [wkhtmltopdf](https://wkhtmltopdf.org/downloads.html)
first, then run:
```
go install github.com/hyperboloide/pdfgen
PDFGEN_TEMPLATES=./templates pdfgen
```

Once installed you can test with something like this:
```
curl -H "Content-Type: application/json" -X POST -d @my_json_file.json \
  http://localhost:8888/invoice > result.pdf
```

Note that the rendering may differ depending on your os (especially OSX) and installed fonts,
that's why it is recommended to test and develop on the Docker environment to
get the same result in production.

## Templates

The PDF are generated from HTML templates. These templates closely ressemble [Django Templates](https://docs.djangoproject.com/en/1.9/ref/templates/language/).

the following template:
```html
<h1>Hello, {{ user }}</h1>
```

can be generated with a `application/json` POST request:

```json
{"user": "fred"}
```

The response is of type `application/pdf` and contains the resulting PDF.

Each PDF template should be in it's own directory under the root directory
defined in the `PDFGEN_TEMPLATES` environment variable.

The urls endpoints will be generated from these directories names. For example a template
in the directory `invoice` will be a reachable at a url that look like that: `http://host:port/invoice`

The template directory must contain an `index.html` file and optionnaly
a `footer.html` file. Other assets like images and CSS should be in
that directory too.
Note that each PDF is generated in isolation and so
your templates should use absolutes paths.
For example if you use bower and have a path like that:
`invoices/bower_components/`
you should have:

```html
<link rel="stylesheet" href="/bower_components/bootstrap/dist/css/bootstrap.min.css" media='screen,print'>
```

Finally don't forget to set the `PDFGEN_TEMPLATES` env variable to the path of
your templates parent directory
Alternatively you copy your templates to either :
`/etc/pdfgen/templates` or `~/.templates`.

## Adding fonts

You could just create a new container with your fonts and rebuild the
cache. Bellow an example Dockerfile.

```
FROM hyperboloide/pdfgen
COPY my_fonts /usr/local/share/fonts/
RUN  fc-cache -f -v
```
