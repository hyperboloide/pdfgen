# pdfgen
HTTP service to generate PDF from Json requests

## How to use it ?
There is docker container :
```
docker run -d \
    -v ~/my_templates/:/templates \
    -e PDFGEN_TEMPLATES=/templates \
    hyperboloide/pdfgen
```


