#!/bin/sh

openssl req -new -key ./src/certs/cert.key -subj "/CN=$1" -sha256 | openssl x509 -req -days 3650 -CA ./src/certs/ca.crt -CAkey ./src/certs/ca.key -set_serial "$2"
