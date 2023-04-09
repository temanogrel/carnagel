-- up
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

CREATE TABLE files (
  uuid             UUID         NOT NULL UNIQUE PRIMARY KEY,
  external_id      BIGINT       NULL     DEFAULT NULL,
  type             SMALLINT     NOT NULL,

  hostname         VARCHAR(50)  NOT NULL,
  path             VARCHAR(150) NOT NULL,

  size             BIGINT       NOT NULL,
  meta             JSONB        NULL,
  pending_deletion BOOLEAN      NOT NULL DEFAULT FALSE,

  created_at       TIMESTAMPTZ  NOT NULL,
  updated_at       TIMESTAMPTZ  NOT NULL
);

CREATE TABLE file_hits (
  uuid           UUID        NOT NULL UNIQUE PRIMARY KEY,
  file_uuid      UUID        NOT NULL,
  remote_address VARCHAR(50) NOT NULL,
  created_at     TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_files_recording
  ON files (external_id);

CREATE INDEX idx_files_recording_type
  ON files (external_id, type);

CREATE INDEX idx_file_recording_type
  ON files (type);

CREATE INDEX idx_storage_path
  ON files (hostname, path);

CREATE INDEX idx_file_hits_file_uuid
  ON file_hits (file_uuid);

ALTER TABLE files
  ADD UNIQUE (hostname, path);

ALTER TABLE file_hits
  ADD CONSTRAINT fk_file_hilts_file_uuid FOREIGN KEY (file_uuid) REFERENCES files (uuid) ON DELETE CASCADE;

-- down
DROP INDEX IF EXISTS idx_files_recording;
DROP INDEX IF EXISTS idx_file_recording_type;
DROP INDEX IF EXISTS idx_storage_path;

DROP TABLE IF EXISTS servers;
DROP TABLE IF EXISTS files;
