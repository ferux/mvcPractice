@echo off
SET CGO_ENABLED=0
SET GOOS=linux
ECHO Starting to build go app
go build -a -installsuffix cgo -o myStoreBackend .
ECHO Done