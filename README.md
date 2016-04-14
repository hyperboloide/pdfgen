# pdfgen
HTTP service to generate PDF from Json requests

## Install

There is docker container :

```
docker run -d \
    -v ~/my_templates/:/templates \
    -e PDFGEN_TEMPLATES=/templates \
    hyperboloide/pdfgen
```

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

The response is a of type `application/pdf` and contains the resulting PDF.

Each PDF template should be in it's own directory under the root directory defined in `PDFGEN_TEMPLATES`.
The urls endpoints will be generated from these directories names. For example a template
in directory `invoice` will be a reachable at a url that look like that: `http://pdfgen:8888/invoice` 

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

Finally don't forget to set the `PDFGEN_TEMPLATES` env variable the path of your templates parent directory. 

## Adding fonts

You could just create a new container with your fonts and rebuild the
cache. Bellow an example Dockerfile.

```
FROM hyperboloide/pdfgen
COPY my_fonts /usr/local/share/fonts/
RUN  fc-cache -f -v
```

