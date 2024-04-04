#!/bin/bash

if ! command -v openssl ?> /dev/null; then
    echo "error: openssl command is not found"
    exit 1
fi

if [[ $# -ne 6 ]]; then
    echo "usage: sh generate.sh <ca_bits> <ca_days> <server_bits> <server_days> <client_bits> <client_days>"
    echo "ex: sh generate.sh 2048 365 2048 365 2048 365"
    exit 1
fi

CA_KEY=ca_certs/ca.key
CA_CRT=ca_certs/ca.crt
CA_BITS=$1
CA_DAYS=$2
SERVER_KEY=server_certs/server.key
SERVER_CSR=server_certs/server.csr
SERVER_CRT=server_certs/server.crt
SERVER_BITS=$3
SERVER_DAYS=$4
CLIENT_KEY=client_certs/client.key
CLIENT_CSR=client_certs/client.csr
CLIENT_CRT=client_certs/client.crt
CLIENT_BITS=$5
CLIENT_DAYS=$6

echo "=========================="
echo "Generating CA certificates"
echo "=========================="
echo ""

openssl genrsa -out $CA_KEY $CA_BITS
sleep 1
openssl req -new -x509 -days $CA_DAYS -key $CA_KEY -out $CA_CRT -noenc
sleep 1

echo ""
echo "=============================="
echo "Generating server certificates"
echo "=============================="
echo ""

openssl genrsa -out $SERVER_KEY $SERVER_BITS
sleep 1
openssl req -new -key $SERVER_KEY -out $SERVER_CSR
sleep 1
openssl x509 -in $SERVER_CSR -req -out $SERVER_CRT -days $SERVER_DAYS -CA $CA_CRT -CAkey $CA_KEY -CAcreateserial
sleep 1
rm -f $SERVER_CSR
sleep 1

echo ""
echo "=============================="
echo "Generating client certificates"
echo "=============================="
echo ""

openssl genrsa -out $CLIENT_KEY $CLIENT_BITS
sleep 1
openssl req -new -key $CLIENT_KEY -out $CLIENT_CSR
sleep 1
openssl x509 -in $CLIENT_CSR -req -out $CLIENT_CRT -days $CLIENT_DAYS -CA $CA_CRT -CAkey $CA_KEY -CAcreateserial
sleep 1
rm -f $CLIENT_CSR
