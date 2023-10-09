package patch

import (
	"fmt"
	"strings"
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
	// eg: "dt", "github.com/xxx/xxx/dt"
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
	// "omit_unused_structs": false
	// 如果为true，sqlc将不会生成在给定包的查询中不使用的表和枚举结构。默认为false。
	// 	设置成false，会生成全量的结构，包括不使用的表和枚举结构。无须防止模块间的类型冲突
	if file == "models" {
		return fmt.Sprintf("../%s/db_models", typeOut)
	}

	// 防止 `"out": "../internal/dao/xxx",` 中的/被解析成目录, 我们只需要最后一项xxx即可
	parts := strings.Split(codeGenOut, "/")
	lastItem := parts[len(parts)-1]
	return fmt.Sprintf("../%s/%s_%s", typeOut, lastItem, file)
}
