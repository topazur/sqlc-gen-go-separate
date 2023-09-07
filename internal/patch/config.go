package patch

import (
	"encoding/json"

	"buf.build/gen/go/sqlc/sqlc/protocolbuffers/go/protos/plugin"
)

type Config struct {
	// 无法享受到 pluginGoCode 的shim，手动补充相关逻辑
	Go plugin.GoCode `json:"go,omitempty"`

	// NOTICE: 插件方式无法享受到 pluginOverride 的shim，所以Override的格式与SQLGo中的不一样
	Overrides []plugin.Override `json:"overrides,omitempty" yaml:"overrides"`
	Rename    map[string]string `json:"rename,omitempty" yaml:"rename"`

	TypePackage string `json:"type_package"` // 生成 type 包的报名
	TypeOut     string `json:"type_out"`     // 生成 type 包的输出目录名，与out是同级目录
	ModuleName  string `json:"module_name"`  // 当前工程的moudle名称
}

func New(pluginOptions []byte) (*Config, error) {
	var conf Config

	// 没有配置时，使用字段对应的零值
	if len(pluginOptions) <= 0 {
		return &conf, nil
	}

	err := json.Unmarshal(pluginOptions, &conf)
	if err != nil {
		return &conf, err
	}

	return &conf, nil
}
