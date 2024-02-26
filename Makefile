build:
	go build -o bin/main.exe src/api/main.go

runbuild:
	bin/main.exe

run:
	go run src/api/main.go

execute: build runbuild