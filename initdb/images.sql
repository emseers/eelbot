CREATE TABLE IF NOT EXISTS images (
  id   integer PRIMARY KEY,
  name text NOT NULL,
  file bytea NOT NULL
);
