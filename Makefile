.PHONY: bin
bin:
	@# 创建文件夹
	mkdir -p bin

.PHONY: bin/sqlc-gen-golang
bin/sqlc-gen-golang: bin go.mod go.sum $(wildcard **/*.go)
	@# 生成二进制文件
	cd plugin && go build -o ../bin/sqlc-gen-golang ./main.go

.PHONY: bin/sqlc-gen-golang.wasm
bin/sqlc-gen-golang.wasm: bin/sqlc-gen-golang
	@# 生成 wasm 文件
	cd plugin && GOOS=wasip1 GOARCH=wasm go build -o ../bin/sqlc-gen-golang.wasm ./main.go

.PHONY: bin/openssl
bin/openssl: bin/sqlc-gen-golang.wasm
	@# 输出 wasm 文件的 SHA256 哈希值
	openssl sha256 ./bin/sqlc-gen-golang.wasm

.PHONY: test
test: bin/sqlc-gen-golang.wasm
	go test ./...
