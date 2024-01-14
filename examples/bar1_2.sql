-- Create "bar1_2" table
CREATE TABLE bar1_2
(
  id   BIGSERIAL PRIMARY KEY
);

-- name: DeleteBar1_2RetExec :exec
DELETE FROM bar1 WHERE id = $1;

