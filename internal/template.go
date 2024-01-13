package golang

import "embed"

//go:embed templates2/*
//go:embed templates2/*/*
var templates embed.FS
