.PHONY: build test

build:
	go build ./...

test: bin/sqlc-gen-golang.wasm
	go test ./...

all: bin/sqlc-gen-golang bin/sqlc-gen-golang.wasm

bin:
	@# 创建文件夹
	mkdir -p bin

bin/sqlc-gen-golang: bin go.mod go.sum $(wildcard **/*.go)
	@# 生成二进制文件
	cd plugin && go build -o ../bin/sqlc-gen-golang ./main.go

bin/sqlc-gen-golang.wasm: bin/sqlc-gen-golang
	@# 生成 wasm 文件
	cd plugin && GOOS=wasip1 GOARCH=wasm go build -o ../bin/sqlc-gen-golang.wasm ./main.go

bin/openssl: bin/sqlc-gen-golang.wasm
	@# 输出 wasm 文件的 SHA256 哈希值
	cd plugin && openssl sha256 ../bin/sqlc-gen-golang.wasm
