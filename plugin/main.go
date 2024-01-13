package main

import (
	"github.com/sqlc-dev/plugin-sdk-go/codegen"

	golang "github.com/topazur/sqlc-gen-go-separate/internal"
)

func main() {
	codegen.Run(golang.Generate)
}
