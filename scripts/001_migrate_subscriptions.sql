CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE subscriptions (
  id          UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
  email       TEXT        NOT NULL,
  city        TEXT        NOT NULL,
  frequency   SMALLINT    NOT NULL DEFAULT 0 CHECK (frequency IN (0,1)),
  token       TEXT        NOT NULL UNIQUE,
  confirmed   BOOLEAN     NOT NULL DEFAULT FALSE,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);