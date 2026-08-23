set shell := ["cmd.exe", "/c"]

dev:
    go tool air -c .air.toml

templ *FLAGS:
    go tool templ generate {{FLAGS}}

tls-unix:
    mkdir -p tls
    cd tls && go run /usr/local/go/src/crypto/tls/generate_cert.go \
        --rsa-bits=2048 \
        --host=localhost

tls-windows:
    if not exist tls mkdir tls
    go run "$(go env GOROOT)/src/crypto/tls/generate_cert.go" --rsa-bits=2048 --host=localhost

fmt:
    go fmt ./...

vet:
    go vet ./...

test:
    go test ./...

getair:
    go get -tool github.com/air-verse/air@latest

check: fmt vet test