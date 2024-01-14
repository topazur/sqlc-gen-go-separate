-- Create "bar1" table
-- https://atlasgo.io/
CREATE TABLE bar1
(
  id   BIGSERIAL PRIMARY KEY
);

-- name: DeleteBar1RetExec :exec
DELETE FROM bar1 WHERE id = $1;

-- name: DeleteBar1RetExecresult :execresult
DELETE FROM bar1 WHERE id = $1;

-- name: DeleteBar1RetExecrows :execrows
DELETE FROM bar1 WHERE id = $1;

-- NOTICE: name: InsertBar1RetExeclastid :execlastid
-- This Query annotations belong mysql
-- INSERT INTO bar1 (id) VALUES ($1);

-- name: GetBar1RetMany :many
SELECT * FROM bar1 WHERE id = $1;

-- name: GetBar1RetOne :one
SELECT * FROM bar1 WHERE id = $1 LIMIT 1;

-- name: DeleteBar1RetBatchexec :batchexec
DELETE FROM bar1 WHERE id = $1;

-- name: GetBar1RetBatchmany :batchmany
SELECT * FROM bar1 WHERE id = $1;

-- name: InsertBar1RetBatchone :batchone
INSERT INTO bar1 (id) VALUES ($1)
RETURNING id;

-- name: InsertBar1RetCopyfrom :copyfrom
INSERT INTO bar1 (id) VALUES ($1);
