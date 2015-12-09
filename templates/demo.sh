#!/bin/bash

#curl -H "Content-Type: application/json" -X POST -d @demo.json http://localhost:8888/invoice > demo.pdf

#     -d @demo.json \

curl \
    -H "Content-Type: application/json" \
    -X POST \
    -d @demo.json \
    http://localhost:8888/invoice \
    > demo.pdf
