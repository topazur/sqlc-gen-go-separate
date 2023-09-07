package main

import (
	"github.com/sqlc-dev/sqlc-go/codegen"

	golang "github.com/topazur/sqlc-gen-go-separate/internal"
)

func main() {
	codegen.Run(golang.Generate)
}
