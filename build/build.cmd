@echo off

go env || exit /b
go vet -v ./... || exit /b
rem go test -v ./... -cover || exit /b
go build -v -o bin\stocker-windows.exe cmd\stocker\stocker.go || exit /b
