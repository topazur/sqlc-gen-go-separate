-- Create "bar2" table
CREATE TABLE bar2
(
  id   BIGSERIAL PRIMARY KEY,
  name text NOT NULL,
  bio  text
);

-- name: GetBar2RetOne :one
SELECT * FROM bar2 WHERE id = $1 LIMIT 1;
