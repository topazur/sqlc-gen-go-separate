package patch

import (
	"fmt"
	"unicode"
)

var (
	typePackage = "dt"
	typeOut     = "dt"
	moduleName  = "github.com/xxx/xxx"
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

// GetTypeOutput 生成type文件的输出路径
// `../${typeOut}` 是为了和 req.Settings.Go.Out 保持同级目录
// `${out}_${file}` 文件名统一加上 req.Settings.Go.Out 前缀，用以区分
func GetTypeOutput(out, file string) string {
	return fmt.Sprintf("../%s/%s_%s", typeOut, out, file)
}
