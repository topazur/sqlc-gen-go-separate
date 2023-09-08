## Directory Structure

```bash
./examples
├── bar1 # sql[0]
│   ├── db.go
│   ├── querier.go
│   └── query.sql.go
├── bar2 # sql[1]
│   ├── db.go
│   ├── querier.go
│   └── query.sql.go
├── dbtype # sql[0] and sql[1]
│   ├── bar1_models.go
│   ├── bar1_query.go
│   ├── bar2_models.go
│   └── bar2_query.go
├── db.go.md # 聚合 sql[0] and sql[1] 查询 (防止运行报错，文件为md格式)
└── ...
```

## Try

```bash
$ cd examples
$ sqlc generate -f sqlc.json # Successful!
```

## Usage

```go
package demo

import (
	"context"
	"fmt"
  "testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
  db "github.com/topazur/sqlc-gen-go-separate/examples/db"
  dbtype "github.com/topazur/sqlc-gen-go-separate/examples/dbtype"
)

func TestDemo(t *testing.T) {
  urlExample := "postgres://username:password@localhost:5432/database_name"

  conn, err := pgx.Connect(context.Background(), urlExample)
  if err != nil {
		t.Fatal(err)
	}
  connPool, err := pgxpool.New(context.Background(), urlExample)
  if err != nil {
		t.Fatal(err)
	}

  store := db.New(conn || connPool)

  arg := dbtype.GetBar2RetOneParams{ID: 1}

  var oBar2 dbtype.Bar2
  oBar2, err = store.Bar2.GetBar2RetOne(context.Background(), arg)
  
  if err != nil {
		t.Fatal(err)
	}
  require.NoError(t, err)
	require.NotEmpty(t, oBar2)
}
```
