
build:
	go build

bench:
	cd benchmarks; go test -bench . -benchmem

test:
	go test ./... -count=1
	cd to; go test ./... -count=1
	cd tests; go test ./... -count=1

format:
	gofmt -w ./

lint:
	gofmt -d ./
	test -z $(shell gofmt -l ./)
