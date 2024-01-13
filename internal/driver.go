package golang

type SQLDriver string

const (
	SQLPackagePGXV4    string = "pgx/v4"
	SQLPackagePGXV5    string = "pgx/v5"
	SQLPackageStandard string = "database/sql"
)

// only the first constant in this group has an explicit type (SA9004)
// 全指定类型 或 全不指定类型
const (
	SQLDriverPGXV4            SQLDriver = "github.com/jackc/pgx/v4"
	SQLDriverPGXV5            SQLDriver = "github.com/jackc/pgx/v5"
	SQLDriverLibPQ            SQLDriver = "github.com/lib/pq"
	SQLDriverGoSQLDriverMySQL SQLDriver = "github.com/go-sql-driver/mysql"
)

func parseDriver(sqlPackage string) SQLDriver {
	switch sqlPackage {
	case SQLPackagePGXV4:
		return SQLDriverPGXV4
	case SQLPackagePGXV5:
		return SQLDriverPGXV5
	default:
		return SQLDriverLibPQ
	}
}

func (d SQLDriver) IsPGX() bool {
	return d == SQLDriverPGXV4 || d == SQLDriverPGXV5
}

func (d SQLDriver) IsGoSQLDriverMySQL() bool {
	return d == SQLDriverGoSQLDriverMySQL
}

func (d SQLDriver) Package() string {
	switch d {
	case SQLDriverPGXV4:
		return SQLPackagePGXV4
	case SQLDriverPGXV5:
		return SQLPackagePGXV5
	default:
		return SQLPackageStandard
	}
}

func (d SQLDriver) Equal(sqlDriver string) bool {
	return sqlDriver == string(d)
}
