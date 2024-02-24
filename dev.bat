@echo off

:: Build wasm app
set GOARCH=wasm
set GOOS=js
go build -o public/app.wasm

:: Run development server
npm start