@echo off

set GOARCH=wasm
set GOOS=js
go build -o public/app.wasm

npm install
npm run build

:: set GOARCH=amd64
:: set GOOS=windows
:: go run ./server
