@echo off
:x
go get -fix -u -v github.com/NatoBoram/Discord-Phone
cls
go run main.go
goto x