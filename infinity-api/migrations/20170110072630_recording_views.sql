-- up
CREATE TABLE recording_views (
  recording_uuid UUID        NOT NULL,
  user_uuid      UUID        NULL     DEFAULT NULL,
  user_hash      VARCHAR(50) NOT NULL,
  created_at     TIMESTAMP   NOT NULL DEFAULT now()
);

CREATE INDEX idx_recording_view_exists_as_guest
  ON recording_views (recording_uuid, user_hash, created_at);

CREATE INDEX idx_recording_view_exists_as_user
  ON recording_views (recording_uuid, user_uuid);

-- down
DROP INDEX IF EXISTS idx_recording_view_exists_as_guest;
DROP INDEX IF EXISTS idx_recording_view_exists_as_user;

DROP TABLE recording_views;
