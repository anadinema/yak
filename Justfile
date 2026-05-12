BINARY := "yak"
CMD_PATH := "./cmd/yak"
INSTALL := env_var("HOME") + "/.local/bin"

default:
    @just --list

build:
    go build -o bin/{{BINARY}} {{CMD_PATH}}

install: build
    mkdir -p {{INSTALL}}
    cp bin/{{BINARY}} {{INSTALL}}/{{BINARY}}

clean:
    rm -rf bin/

test:
    go test ./...

lint:
    golangci-lint run ./...

tidy:
    go mod tidy

fmt:
    golangci-lint fmt ./...

fmt-check:
    golangci-lint fmt --diff ./...
