package patch

import (
	"fmt"
	"unicode"
)

var (
	typePackage = "dt"
	typeOut     = "dt"
	moduleName  = "github.com/xxx/xxx"

	// `sql.Plugin.Out` absolute existed
	// https://github.com/topazur/sqlc/blob/5e81d02d80eae1cc01e11592c318877eae5b14ff/internal/cmd/generate.go#L379
	codeGenOut = ""
)

// GetTypeImportParams 生成被分离出去的type包的import链接
func GetTypeImportParams() (string, string) {
	return typePackage, moduleName + "/" + typeOut
}

// QueryExportedType 使用type包中的类型时加上包名
func QueryExportedType(t string) string {
	// 判断首字母是否为小写字母 <大写是导入>
	isLower := unicode.IsLower([]rune(t)[0])
	if isLower {
		return t
	}

	return typePackage + "." + t
}

// GetTypePackage 生成type包的包名
func GetTypePackage() string {
	return typePackage
}

// GetTypeOutput 生成type文件的输出路径
// `../${typeOut}` 是为了和 req.Settings.Go.Out 保持同级目录
// `${out}_${file}` 文件名统一加上 req.Settings.Go.Out 前缀，用以区分
func GetTypeOutput(file string) string {
	return fmt.Sprintf("../%s/%s_%s", typeOut, codeGenOut, file)
}
