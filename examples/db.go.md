```go
package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	bar1 "github.com/topazur/sqlc-gen-go-separate/examples/bar1"
	bar2 "github.com/topazur/sqlc-gen-go-separate/examples/bar2"
)

type DBTX interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	SendBatch(context.Context, *pgx.Batch) pgx.BatchResults

	Begin(ctx context.Context) (pgx.Tx, error)
}

// Queries 聚合各个 sqlc 生成的Queries
// 1️⃣ 不细分 sqlc 代码导致难以阅读与维护 - 放弃❌
// 2️⃣ 细分 sqlc 代码后，使用起来难以管理导入 - 先细分再聚合✅
// 该文件是模仿 `[sqlc_out]/db.go` 生成的。
type Queries struct {
	db DBTX

	Bar1 *bar1.Queries
	Bar2 *bar2.Queries
}

func New(db DBTX) *Queries {
	return &Queries{
		db: db,

		Bar1: bar1.New(db),
		Bar2: bar2.New(db),
	}
}

func (q *Queries) withTx(tx pgx.Tx) *Queries {
	return &Queries{
		db: tx,

		Bar1: bar1.New(tx),
		Bar2: bar2.New(tx),
	}
}

func (q *Queries) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := q.db.Begin(ctx)
	if err != nil {
		return err
	}

	// NOTICE: WithTx 得到专用于事务的 Queries 实例
	txQueries := q.withTx(tx)

	err = fn(txQueries)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
```
