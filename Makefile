.PHONY: all tls check

tls:
	mkdir -p tls
	cd tls && go run /usr/local/go/src/crypto/tls/generate_cert.go \
		--rsa-bits=2048 \
		--host=localhost
fmt:
	gofmt -w .

vet:
	go vet ./...

check: fmt vet test

dev:
	@go run ./cmd/web/