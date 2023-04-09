-- up
ALTER TABLE files
  ADD upstore_hash VARCHAR(20) NULL;

ALTER TABLE files
  ADD pending_upload BOOL NOT NULL DEFAULT FALSE;

CREATE INDEX idx_files_pending_upload
  ON files (pending_upload);

CREATE INDEX idx_files_pending_deletion
  ON files (pending_deletion);

CREATE UNIQUE INDEX idx_files_upstore_hash_unique
  ON files (upstore_hash);

-- down
DROP INDEX IF EXISTS idx_files_upstore_hash_unique;
DROP INDEX IF EXISTS idx_files_pending_deletion;
DROP INDEX IF EXISTS idx_files_pending_upload;

ALTER TABLE files
  DROP upstore_hash;

ALTER TABLE files
  DROP pending_upload;
