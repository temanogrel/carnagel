-- up
DROP TABLE servers;

-- down
CREATE TABLE servers (
  uuid              UUID        NOT NULL UNIQUE PRIMARY KEY,
  switch            VARCHAR(50) NOT NULL,
  type              SMALLINT    NOT NULL,
  enabled           BOOLEAN     NOT NULL DEFAULT TRUE,

  internal_hostname VARCHAR(50) NOT NULL,
  external_hostname VARCHAR(50) NOT NULL,

  created_at        TIMESTAMPTZ NOT NULL,
  updated_at        TIMESTAMPTZ NOT NULL
);
