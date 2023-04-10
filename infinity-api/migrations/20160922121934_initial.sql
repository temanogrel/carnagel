-- up
CREATE TABLE performers (
  uuid            UUID PRIMARY KEY,
  external_id     INTEGER        NOT NULL UNIQUE,
  origin_service  VARCHAR(15)    NOT NULL,
  origin_section  VARCHAR(20)    NULL,
  stage_name      VARCHAR(40)    NOT NULL,
  aliases         VARCHAR(40) [] NOT NULL,
  recording_count INTEGER        NOT NULL DEFAULT 0,
  created_at      TIMESTAMP      NOT NULL,
  updated_at      TIMESTAMP      NOT NULL
);

CREATE TABLE recordings (
  uuid           UUID PRIMARY KEY,
  performer_uuid UUID        NOT NULL,
  external_id    INTEGER     NOT NULL UNIQUE,
  stage_name     VARCHAR(40) NOT NULL,
  duration       INTEGER     NOT NULL,
  video_uuid     UUID        NOT NULL,
  collage_uuid   UUID        NOT NULL,
  video_manifest TEXT        NOT NULL,
  sprites        UUID []     NOT NULL,
  images         UUID []     NOT NULL,
  view_count     INTEGER     NOT NULL DEFAULT 0,
  like_count     INTEGER     NOT NULL DEFAULT 0,
  created_at     TIMESTAMP   NOT NULL,
  updated_at     TIMESTAMP   NOT NULL
);

CREATE TABLE users (
  uuid       UUID PRIMARY KEY,
  email      VARCHAR(255) NOT NULL,
  username   VARCHAR(40)  NOT NULL,
  password   VARCHAR(100) NOT NULL,
  role       VARCHAR(10)  NOT NULL,
  created_at TIMESTAMP    NOT NULL,
  updated_at TIMESTAMP    NOT NULL
);

CREATE TABLE user_recording_likes (
  user_uuid      UUID      NOT NULL,
  recording_uuid UUID      NOT NULL,
  created_at     TIMESTAMP NOT NULL
);

CREATE TABLE user_performer_favorites (
  user_uuid      UUID NOT NULL,
  performer_uuid UUID NOT NULL
);

CREATE TABLE bandwidth_consumption (
  recording_uuid UUID      NOT NULL,
  user_uuid      UUID DEFAULT NULL,
  session_uuid   UUID      NOT NULL,
  bytes          BIGINT    NOT NULL,
  created_at     TIMESTAMP NOT NULL
);

ALTER TABLE recordings
  ADD CONSTRAINT fk_performer_uuid FOREIGN KEY (performer_uuid) REFERENCES performers (uuid) ON DELETE CASCADE;

ALTER TABLE user_recording_likes
  ADD CONSTRAINT fk_recording_uuid FOREIGN KEY (recording_uuid) REFERENCES recordings (uuid) ON DELETE CASCADE;

ALTER TABLE user_recording_likes
  ADD CONSTRAINT fk_user_uuid FOREIGN KEY (user_uuid) REFERENCES users (uuid) ON DELETE RESTRICT;

ALTER TABLE user_performer_favorites
  ADD CONSTRAINT fk_performer_uuid FOREIGN KEY (performer_uuid) REFERENCES performers (uuid) ON DELETE CASCADE;

ALTER TABLE user_performer_favorites
  ADD CONSTRAINT fk_user_uuid FOREIGN KEY (user_uuid) REFERENCES users (uuid) ON DELETE CASCADE;

ALTER TABLE bandwidth_consumption
  ADD CONSTRAINT fk_recording_uuid FOREIGN KEY (recording_uuid) REFERENCES recordings (uuid) ON DELETE CASCADE;

ALTER TABLE bandwidth_consumption
  ADD CONSTRAINT fk_user_uuid FOREIGN KEY (user_uuid) REFERENCES users (uuid) ON DELETE CASCADE;

-- Create a unique index on the lowercase version of the username since driver in case sensitive
CREATE UNIQUE INDEX idx_users_unique_username
  ON users (lower(username));

CREATE UNIQUE INDEX idx_users_unique_email
  ON users (lower(email));

CREATE INDEX idx_performers_external_id
  ON performers (external_id);

CREATE INDEX idx_recording_external_id
  ON recordings (external_id);

CREATE INDEX idx_recording_created_at
  ON recordings (created_at);

CREATE INDEX idx_bandwidth_consumption_user_uuid
  ON bandwidth_consumption (user_uuid);

CREATE INDEX idx_bandwidth_consumption_session_uuid
  ON bandwidth_consumption (session_uuid);

CREATE INDEX idx_bandwidth_consumption_for_user_on_recording
  ON bandwidth_consumption (user_uuid, recording_uuid);

CREATE UNIQUE INDEX idx_bandwidth_consumption_daily_for_session_recording
  ON bandwidth_consumption (session_uuid, recording_uuid, CAST(created_at AS DATE));

-- down
DROP INDEX IF EXISTS idx_recording_external_id;
DROP INDEX IF EXISTS idx_performers_external_id;
DROP INDEX IF EXISTS idx_users_unique_username;

ALTER TABLE recordings
  DROP CONSTRAINT fk_performer_uuid;
ALTER TABLE user_recording_likes
  DROP CONSTRAINT fk_user_uuid;
ALTER TABLE user_recording_likes
  DROP CONSTRAINT fk_recording_uuid;
ALTER TABLE user_performer_favorites
  DROP CONSTRAINT fk_user_uuid;
ALTER TABLE user_performer_favorites
  DROP CONSTRAINT fk_performer_uuid;
ALTER TABLE bandwidth_consumption
  DROP CONSTRAINT fk_user_uuid;
ALTER TABLE bandwidth_consumption
  DROP CONSTRAINT fk_recording_uuid;

DROP TABLE recordings;
DROP TABLE performers;
DROP TABLE user_performer_favorites;
DROP TABLE user_recording_likes;
DROP TABLE users;
DROP TABLE bandwidth_consumption;