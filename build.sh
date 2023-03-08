#!/bin/sh

# Compila la función createOrder en un archivo binario
GOOS=linux GOARCH=amd64 go build -o bin/createOrder src/createOrder/main.go

# Compila la función processPayment en un archivo binario
GOOS=linux GOARCH=amd64 go build -o bin/processPayment src/processPayment/main.go
