#!/bin/bash
go run main.go worker --id=1 --http=:8071 --grpc=:9081
go run main.go worker --id=2 --http=:8072 --grpc=:9082