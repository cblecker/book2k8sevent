.PHONY: default
default: run

.PHONY: run
run:
	go run .

bindata.go:
	go-bindata data/
