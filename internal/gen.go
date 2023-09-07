package golang

import (
	"context"

	"buf.build/gen/go/sqlc/sqlc/protocolbuffers/go/protos/plugin"
)

func Generate(_ context.Context, req *plugin.CodeGenRequest) (*plugin.CodeGenResponse, error) {

	resp := plugin.CodeGenResponse{}

	return &resp, nil
}
