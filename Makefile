.PHONY: default
default: go-build

.PHONY: go-build
go-build:
	go build -o _output/ .

bindata.go:
	go-bindata data/
