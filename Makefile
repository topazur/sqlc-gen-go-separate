all: sqlc-gen-golang sqlc-gen-golang.wasm

sqlc-gen-golang:
	cd plugin && go build -o ~/bin/sqlc-gen-golang ./main.go

sqlc-gen-golang.wasm:
	cd plugin && GOOS=wasip1 GOARCH=wasm go build -o sqlc-gen-golang.wasm main.go
	openssl sha256 plugin/sqlc-gen-golang.wasm
