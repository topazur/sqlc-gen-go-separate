package golang

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"go/format"
	"strings"
	"text/template"

	"github.com/meta-programming/go-codegenutil/unusedimports"
	"github.com/sqlc-dev/plugin-sdk-go/metadata"
	"github.com/sqlc-dev/plugin-sdk-go/plugin"
	"github.com/sqlc-dev/plugin-sdk-go/sdk"
	"github.com/topazur/sqlc-gen-go-separate/internal/opts"
)

type tmplCtx struct {
	Q           string
	Package     string
	TypePackage string
	SQLDriver   SQLDriver
	Enums       []Enum
	Structs     []Struct
	GoQueries   []Query
	SqlcVersion string

	// TODO: Race conditions
	SourceName string

	EmitJSONTags              bool
	JsonTagsIDUppercase       bool
	EmitDBTags                bool
	EmitPreparedQueries       bool
	EmitInterface             bool
	EmitEmptySlices           bool
	EmitMethodsWithDBArgument bool
	EmitEnumValidMethod       bool
	EmitAllEnumValues         bool
	UsesCopyFrom              bool
	UsesBatch                 bool
	OmitSqlcVersion           bool
	BuildTags                 string
}

func (t *tmplCtx) OutputQuery(sourceName string) bool {
	return t.SourceName == sourceName
}

func (t *tmplCtx) codegenDbarg() string {
	if t.EmitMethodsWithDBArgument {
		return "db DBTX, "
	}
	return ""
}

// Called as a global method since subtemplate queryCodeStdExec does not have
// access to the toplevel tmplCtx
func (t *tmplCtx) codegenEmitPreparedQueries() bool {
	return t.EmitPreparedQueries
}

func (t *tmplCtx) codegenQueryMethod(q Query) string {
	db := "q.db"
	if t.EmitMethodsWithDBArgument {
		db = "db"
	}

	switch q.Cmd {
	case ":one":
		if t.EmitPreparedQueries {
			return "q.queryRow"
		}
		return db + ".QueryRowContext"

	case ":many":
		if t.EmitPreparedQueries {
			return "q.query"
		}
		return db + ".QueryContext"

	default:
		if t.EmitPreparedQueries {
			return "q.exec"
		}
		return db + ".ExecContext"
	}
}

func (t *tmplCtx) codegenQueryRetval(q Query) (string, error) {
	switch q.Cmd {
	case ":one":
		return "row :=", nil
	case ":many":
		return "rows, err :=", nil
	case ":exec":
		return "_, err :=", nil
	case ":execrows", ":execlastid":
		return "result, err :=", nil
	case ":execresult":
		return "return", nil
	default:
		return "", fmt.Errorf("unhandled q.Cmd case %q", q.Cmd)
	}
}

func Generate(ctx context.Context, req *plugin.GenerateRequest) (*plugin.GenerateResponse, error) {
	options, err := opts.Parse(req)
	if err != nil {
		return nil, err
	}

	if err := opts.ValidateOpts(options); err != nil {
		return nil, err
	}

	enums := buildEnums(req, options)
	structs := buildStructs(req, options)
	queries, err := buildQueries(req, options, structs)
	if err != nil {
		return nil, err
	}

	if options.OmitUnusedStructs {
		enums, structs = filterUnusedStructs(enums, structs, queries)
	}

	if err := validate(options, enums, structs, queries); err != nil {
		return nil, err
	}

	return generate(req, options, enums, structs, queries)
}

func validate(options *opts.Options, enums []Enum, structs []Struct, queries []Query) error {
	enumNames := make(map[string]struct{})
	for _, enum := range enums {
		enumNames[enum.Name] = struct{}{}
		enumNames["Null"+enum.Name] = struct{}{}
	}
	structNames := make(map[string]struct{})
	for _, struckt := range structs {
		if _, ok := enumNames[struckt.Name]; ok {
			return fmt.Errorf("struct name conflicts with enum name: %s", struckt.Name)
		}
		structNames[struckt.Name] = struct{}{}
	}
	if !options.EmitExportedQueries {
		return nil
	}
	for _, query := range queries {
		if _, ok := enumNames[query.ConstantName]; ok {
			return fmt.Errorf("query constant name conflicts with enum name: %s", query.ConstantName)
		}
		if _, ok := structNames[query.ConstantName]; ok {
			return fmt.Errorf("query constant name conflicts with struct name: %s", query.ConstantName)
		}
	}
	return nil
}

func generate(req *plugin.GenerateRequest, options *opts.Options, enums []Enum, structs []Struct, queries []Query) (*plugin.GenerateResponse, error) {
	i := &importer{
		Options: options,
		Queries: queries,
		Enums:   enums,
		Structs: structs,
	}

	tctx := tmplCtx{
		EmitInterface:             options.EmitInterface,
		EmitJSONTags:              options.EmitJsonTags,
		JsonTagsIDUppercase:       options.JsonTagsIdUppercase,
		EmitDBTags:                options.EmitDbTags,
		EmitPreparedQueries:       options.EmitPreparedQueries,
		EmitEmptySlices:           options.EmitEmptySlices,
		EmitMethodsWithDBArgument: options.EmitMethodsWithDbArgument,
		EmitEnumValidMethod:       options.EmitEnumValidMethod,
		EmitAllEnumValues:         options.EmitAllEnumValues,
		UsesCopyFrom:              usesCopyFrom(queries),
		UsesBatch:                 usesBatch(queries),
		SQLDriver:                 parseDriver(options.SqlPackage),
		Q:                         "`",
		Package:                   options.Package,
		TypePackage:               opts.DBTYPE_PACKAGE,
		Enums:                     enums,
		Structs:                   structs,
		SqlcVersion:               req.SqlcVersion,
		BuildTags:                 options.BuildTags,
		OmitSqlcVersion:           options.OmitSqlcVersion,
	}

	// 阻止非 pgx or mysql 驱动使用 :copyfrom 查询注释
	if tctx.UsesCopyFrom && !tctx.SQLDriver.IsPGX() && !SQLDriverGoSQLDriverMySQL.Equal(options.SqlDriver) {
		return nil, errors.New(":copyfrom is only supported by pgx and github.com/go-sql-driver/mysql")
	}

	// 当 mysql 驱动使用 :copyfrom 查询注释时，参数不支持 time.Time 类型
	if tctx.UsesCopyFrom && SQLDriverGoSQLDriverMySQL.Equal(options.SqlDriver) {
		if err := checkNoTimesForMySQLCopyFrom(queries); err != nil {
			return nil, err
		}
		tctx.SQLDriver = SQLDriverGoSQLDriverMySQL
	}

	// 仅 pgx 支持批量查询
	if tctx.UsesBatch && !tctx.SQLDriver.IsPGX() {
		return nil, errors.New(":batch* commands are only supported by pgx")
	}

	funcMap := template.FuncMap{
		"lowerTitle": sdk.LowerTitle,
		"comment":    sdk.DoubleSlashComment,
		"escape":     sdk.EscapeBacktick,
		"imports":    i.Imports,
		"hasPrefix":  strings.HasPrefix,

		// These methods are Go specific, they do not belong in the codegen package
		// (as that is language independent)
		"dbarg":               tctx.codegenDbarg,
		"emitPreparedQueries": tctx.codegenEmitPreparedQueries,
		"queryMethod":         tctx.codegenQueryMethod,
		"queryRetval":         tctx.codegenQueryRetval,
	}

	tmpl := template.Must(
		template.New("table").
			Funcs(funcMap).
			ParseFS(
				templates,
				"templates2/*.tmpl",
				"templates2/*/*.tmpl",
			),
	)

	output := map[string]string{}

	execute := func(name, templateName string) error {
		// 二维数组 - 当前文件导入的package信息
		imports := i.Imports(name)
		// 防止参数名(query.Arg.Name)与包名冲突
		replacedQueries := replaceConflictedArg(imports, queries)

		var b bytes.Buffer
		w := bufio.NewWriter(&b)
		tctx.SourceName = name
		tctx.GoQueries = replacedQueries
		err := tmpl.ExecuteTemplate(w, templateName, &tctx)
		w.Flush()
		if err != nil {
			return err
		}
		code, err := format.Source(b.Bytes())
		if err != nil {
			fmt.Println(b.String())
			return fmt.Errorf("source error: %w", err)
		}

		/*
			NOTICE: 当 "queries": ["queries/one.sql", "queries/two.sql"] 存在多个的时候，queryFile 和 typeFile 也会生成多个
			rootDir
			├── outDir
			│   ├── db.go
			│   ├── db_batch.go
			│   ├── db_querier.go
			│   └── ****.sql.sqlc.go # queryFile
			├── dbtype
			│   ├── db_models.go # 表定义
			│   └── ****.query.go # typeFile
			├── db.go # 聚合所有方法
			└── ****.go # 其他文件
		*/
		if templateName == "queryFile" && options.OutputFilesSuffix != "" {
			// eg: `****.sql.sqlc.go` => `outDir/****.sql.sqlc.go`
			name += options.OutputFilesSuffix
		}

		// NOTICE: 分离点 - queryFile's type
		if templateName == "typeFile" {
			// eg: `../dbtype/****.tsql.go` => `dbtype/****.tsql.go`，其中 dbtype 和 outDir 是同级的目录
			name = strings.Replace(name, ".sql", ".tsql", 1)
			name = fmt.Sprintf("../%s/%s", opts.DBTYPE_PACKAGE, name)
		}

		// NOTICE: 分离点 - models
		if templateName == "modelsFile" {
			// 利用 `"schema": "migration/*"` 和 `"omit_unused_structs": false` 两个配置，最终仅生成一个model，包含全量的表定义。无须防止模块间的类型冲突
			// eg: `../dbtype/db_models.go` => `dbtype/db_models.go`
			name = fmt.Sprintf("../%s/%s", opts.DBTYPE_PACKAGE, options.OutputModelsFileName)
		}

		if !strings.HasSuffix(name, ".go") {
			name += ".go"
		}

		// NOTICE: 删除多余的 import 语句，未使用的 package 会被移除
		codeStr, err := unusedimports.PruneUnparsed(name, string(code))
		if err != nil {
			return fmt.Errorf("unusedimports.PruneUnparsed error: %w", err)
		}

		output[name] = codeStr
		return nil
	}

	dbFileName := "db.go"
	if options.OutputDbFileName != "" {
		dbFileName = options.OutputDbFileName
	}
	modelsFileName := "models.go"
	if options.OutputModelsFileName != "" {
		modelsFileName = options.OutputModelsFileName
	}
	querierFileName := "querier.go"
	if options.OutputQuerierFileName != "" {
		querierFileName = options.OutputQuerierFileName
	}
	copyfromFileName := "copyfrom.go"
	if options.OutputCopyfromFileName != "" {
		copyfromFileName = options.OutputCopyfromFileName
	}

	batchFileName := "batch.go"
	if options.OutputBatchFileName != "" {
		batchFileName = options.OutputBatchFileName
	}

	if err := execute(dbFileName, "dbFile"); err != nil {
		return nil, err
	}
	if err := execute(modelsFileName, "modelsFile"); err != nil {
		return nil, err
	}
	if options.EmitInterface {
		if err := execute(querierFileName, "interfaceFile"); err != nil {
			return nil, err
		}
	}
	if tctx.UsesCopyFrom {
		if err := execute(copyfromFileName, "copyfromFile"); err != nil {
			return nil, err
		}
	}
	if tctx.UsesBatch {
		if err := execute(batchFileName, "batchFile"); err != nil {
			return nil, err
		}
	}

	files := map[string]struct{}{}
	for _, gq := range queries {
		files[gq.SourceName] = struct{}{}
	}

	for source := range files {
		if err := execute(source, "queryFile"); err != nil {
			return nil, err
		}

		// NOTICE: 由于 SQLGo 模式的类型都是生成在query文件里，所以这里无需对 copyfromFileName、batchFileName 两种文件额外处理
		if err := execute(source, "typeFile"); err != nil {
			return nil, err
		}
	}
	resp := plugin.GenerateResponse{}

	for filename, code := range output {
		resp.Files = append(resp.Files, &plugin.File{
			Name:     filename,
			Contents: []byte(code),
		})
	}

	return &resp, nil
}

func usesCopyFrom(queries []Query) bool {
	for _, q := range queries {
		if q.Cmd == metadata.CmdCopyFrom {
			return true
		}
	}
	return false
}

func usesBatch(queries []Query) bool {
	for _, q := range queries {
		for _, cmd := range []string{metadata.CmdBatchExec, metadata.CmdBatchMany, metadata.CmdBatchOne} {
			if q.Cmd == cmd {
				return true
			}
		}
	}
	return false
}

func checkNoTimesForMySQLCopyFrom(queries []Query) error {
	for _, q := range queries {
		if q.Cmd != metadata.CmdCopyFrom {
			continue
		}
		for _, f := range q.Arg.CopyFromMySQLFields() {
			if f.Type == "time.Time" {
				return fmt.Errorf("values with a timezone are not yet supported")
			}
		}
	}
	return nil
}

func filterUnusedStructs(enums []Enum, structs []Struct, queries []Query) ([]Enum, []Struct) {
	keepTypes := make(map[string]struct{})

	for _, query := range queries {
		if !query.Arg.isEmpty() {
			keepTypes[query.Arg.Type()] = struct{}{}
			if query.Arg.IsStruct() {
				for _, field := range query.Arg.Struct.Fields {
					keepTypes[field.Type] = struct{}{}
				}
			}
		}
		if query.hasRetType() {
			keepTypes[query.Ret.Type()] = struct{}{}
			if query.Ret.IsStruct() {
				for _, field := range query.Ret.Struct.Fields {
					keepTypes[field.Type] = struct{}{}
					for _, embedField := range field.EmbedFields {
						keepTypes[embedField.Type] = struct{}{}
					}
				}
			}
		}
	}

	keepEnums := make([]Enum, 0, len(enums))
	for _, enum := range enums {
		_, keep := keepTypes[enum.Name]
		_, keepNull := keepTypes["Null"+enum.Name]
		if keep || keepNull {
			keepEnums = append(keepEnums, enum)
		}
	}

	keepStructs := make([]Struct, 0, len(structs))
	for _, st := range structs {
		if _, ok := keepTypes[st.Name]; ok {
			keepStructs = append(keepStructs, st)
		}
	}

	return keepEnums, keepStructs
}
