CREATE TABLE IF NOT EXISTS taunts (
  id   integer PRIMARY KEY,
  name text NOT NULL,
  file bytea NOT NULL
);
