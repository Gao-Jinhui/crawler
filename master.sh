#!/bin/bash
go run main.go master --id=1 --http=:8081 --grpc=:9091
go run main.go master --id=2 --http=:8082 --grpc=:9092