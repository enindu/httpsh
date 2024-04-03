#!/bin/bash

if ! command -v openssl ?> /dev/null; then
    echo "openssl command is not found"
    exit 1
fi

openssl req \
    -nodes \
    -sha256 \
    -x509 \
    -days 365 \
    -keyout server.key \
    -newkey rsa:4096 \
    -out server.crt \
