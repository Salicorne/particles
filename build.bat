@echo off

set GOARCH=wasm
set GOOS=js
go build -o web/app.wasm

set GOARCH=amd64
set GOOS=windows
go run ./server