package patch

import (
	"go/format"

	"github.com/meta-programming/go-codegenutil/unusedimports"
)

type FormatTmpl string

// String return string type
func (f FormatTmpl) String() string {
	return string(f)
}

// SliceByte return []byte type
func (f FormatTmpl) SliceByte() []byte {
	return []byte(f)
}

// Source is go/format fo file
func (f FormatTmpl) Source() (FormatTmpl, error) {
	code, err := format.Source(f.SliceByte())
	return FormatTmpl(code), err
}

// Unusedimports is remove unused imports in go file
func (f FormatTmpl) Unusedimports(filename string) (FormatTmpl, error) {
	code, err := unusedimports.PruneUnparsed(filename, f.String())
	return FormatTmpl(code), err
}
