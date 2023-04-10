-- up
DROP TABLE IF EXISTS user_performer_favorites;

CREATE TABLE user_recording_favorites (
  user_uuid      UUID NOT NULL,
  recording_uuid UUID NOT NULL
);

ALTER TABLE user_recording_favorites
  ADD CONSTRAINT fk_recording_uuid FOREIGN KEY (recording_uuid) REFERENCES recordings (uuid) ON DELETE CASCADE;

ALTER TABLE user_recording_favorites
  ADD CONSTRAINT fk_user_uuid FOREIGN KEY (user_uuid) REFERENCES users (uuid) ON DELETE RESTRICT;

CREATE INDEX idx_user_recording_favorite_keys
  ON user_recording_favorites (user_uuid, recording_uuid);

-- down

DROP INDEX IF EXISTS idx_user_recording_favorite_keys;

DROP TABLE user_recording_favorites;