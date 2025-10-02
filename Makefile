install-deps:
	go install tool

gen-mocks:
	mockery

injection:
	wire gen internal/injection/wire.go

build:
	make install-deps
	make gen-mocks
	make injection

lint:
	golangci-lint run
