package patch

import "buf.build/gen/go/sqlc/sqlc/protocolbuffers/go/protos/plugin"

func PluginTypeCode(typePackageArg, typeOutArg, moduleNameArg string) {
	if typePackageArg != "" {
		typePackage = typePackageArg
	}

	if typeOutArg != "" {
		typeOut = typeOutArg
	}

	if moduleNameArg != "" {
		moduleName = moduleNameArg
	}
}

// PluginGoCode sqlc中的逻辑是转换类型：`config.SQLGo` => `plugin.GoCode`
// https://github.com/sqc-dev/sqlc/blob/3c9ef73dd379613ff682326a58d402f0695f3242/internal/cmd/shim.go#L78
func PluginGoCode(s *plugin.GoCode) *plugin.GoCode {
	if s.QueryParameterLimit == nil {
		s.QueryParameterLimit = new(int32) // 0
		*s.QueryParameterLimit = 1
	}

	// 默认生成全量的结构
	s.OmitUnusedStructs = false

	return s
}

// PluginRenameCode 合并map
func PluginRenameCode(global, local map[string]string) map[string]string {
	merged := make(map[string]string)

	// local 覆盖 global
	for k, v := range global {
		merged[k] = v
	}
	for k, v := range local {
		merged[k] = v
	}

	return merged
}

// PluginOverride 合并切片
func PluginOverride(global []*plugin.Override, local []plugin.Override) []*plugin.Override {
	merged := make([]*plugin.Override, len(global)+len(local))

	// 将 global 切片的元素复制到 merged 中
	copy(merged, global)

	// 将 local 切片的元素添加到 merged 中
	// 在Go语言中，sync.Mutex类型不能被拷贝；使用索引来访问切片元素可以保证不会发生sync.Mutex的拷贝操作
	for i := range local {
		merged = append(merged, &local[i])
	}

	return merged
}
