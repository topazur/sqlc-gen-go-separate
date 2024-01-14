
# [go](https://docs.sqlc.dev/en/latest/reference/config.html#go)


- package
> 用于生成的代码的包名称。Defaults to `out` basename.

- out
> 生成的代码的输出目录。

- sql_package
> Either pgx/v4, pgx/v5 or database/sql. Defaults to database/sql.

- emit_db_tags
> 如果为true，则向生成的结构添加DB标记。默认为 false 。

- emit_prepared_queries
> 如果为true，则包括对预查询的支持。默认为 false 。

- emit_interface
> 如果为true，则在生成的包中输出一个 Querier 接口。默认为 false 。

- emit_exact_table_names
> 如果为true，结构名称与表名称相同。否则，sqlc 会尝试将复数表名单数化。默认为 false 。

- emit_empty_slices
> 如果为true，则 :many 查询返回的切片将为空切片，而不是 nil 。默认为 false 。

- emit_exported_queries
> 如果为true，则可以导出自动生成的SQL语句以供另一个包访问。

- emit_json_tags
> 如果为true，则向生成的结构添加JSON标记。默认为 false 。

- emit_result_struct_pointers
> 如果为true，则将查询结果作为指向结构的指针返回。返回多个结果的指针以指针切片的形式返回。默认为 false 。

- emit_params_struct_pointers
> 如果为true，则参数作为指向结构的指针传递。默认为 false 。

- emit_methods_with_db_argument
> 如果为true，则生成的方法将接受 DBTX 参数，而不是在 *Queries 结构体上存储 DBTX 。

- emit_pointers_for_null_types
> 如果为true，为 nullable 列生成的类型将作为指针发出 (ie. *string)，而不是 database/sql null types (ie. NullString)。目前只支持 PostgreSQL ，如果 sql 包是 pgx/v4 或 pgx/v5 ，以及 SQLite 。默认为 false 。

- emit_enum_valid_method
> 如果为true，则对枚举类型生成 Valid 方法，判断是否为有效的枚举值。

- emit_all_enum_values
> 如果为true，则为每个枚举类型发出一个返回所有有效枚举值的函数。

- emit_sql_as_comment
> 如果为true，则将SQL语句作为代码块注释发送到生成的函数上方，并附加到任何现有的注释。默认为 false 。

- build_tags
> 如果为true，在每个生成的 Go 文件的开头添加一个 `//go:build <build_tags>` 指令。

- json_tags_id_uppercase
> 如果为true，json标签中的 “Id” 将是 uppercase 格式。如果为 false，将是 camelcase 格式。

- json_tags_case_style
> `camel` for camelCase, `pascal` for PascalCase, `snake` for snake_case or `none` to use the column name in the DB. Defaults to `none`.

- omit_unused_structs
> 如果是 true ，sqlc将不会生成在给定包的查询中未使用的表和枚举结构。默认为 false 。

- output_batch_file_name
> 自定义批量查询文件的名称。默认为 `batch.go` 。

- output_db_file_name
> 自定义db文件的名称。默认为 `db.go` 。

- output_models_file_name
> 自定义模型文件的名称。默认为 `models.go` 。

- output_querier_file_name
> 自定义查询器文件的名称。默认为 `querier.go` 。

- output_copyfrom_file_name
> 自定义copyfrom文件的名称。默认为 `copyfrom.go` 。

- output_files_suffix
> 如果指定，后缀将添加到生成的文件的名称中。

- query_parameter_limit
> 将为 Go 函数生成的位置参数的数量。要始终发出参数结构体，请将此设置为 0。Defaults to 1.

- inflection_exclude_table_names
> 该字符串列表在表单单数化时会被排除，直接返回

- omit_sqlc_version
> 模板文件顶部是否添加版本号注释

- sql_driver
> ???

- rename
> 自定义生成的结构字段的名称。有关用法信息，请参阅重命名字段。

- overrides
> 它是定义的集合，指示使用哪些类型来映射数据库类型。


## 推荐配置

```json5
{
  "sql_package": "pgx/v5",
  "package": "db",
  "out": "todo_dir", // custom
  "module_name": "todo_path:github.com/***/dbtype", // custom

  /** 输出文件 */
  "output_db_file_name": "db.go",
  "output_querier_file_name": "db_querier.go",
  "output_batch_file_name": "db_batch.go",
  "output_models_file_name": "db_models.go",
  "output_copyfrom_file_name": "db_copyfrom.go",
  "output_files_suffix": ".sqlc",

  /** 结构体标签 */
  "emit_json_tags": true,
  "json_tags_case_style": "camel",
  "json_tags_id_uppercase": false,
  "emit_db_tags": true,

  /** 预查询 + interface + sql + enum + table_name */
  "emit_prepared_queries": false,
  "emit_interface": true,
  "emit_exported_queries": true,
  "emit_enum_valid_method": true,
  "emit_all_enum_values": true,
  "emit_exact_table_names": false,
  "inflection_exclude_table_names": [],
  "omit_unused_structs": false,

  /** 参数及返回值 */
  "emit_empty_slices": true,
  "emit_result_struct_pointers": false,
  "emit_params_struct_pointers": false,
  "emit_methods_with_db_argument": false,
  "emit_pointers_for_null_types": false,
  "query_parameter_limit": 0

  // build_tags
  // omit_sqlc_version
  // sql_driver
  // rename
  // overrides
}
```
