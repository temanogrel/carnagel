-- up
CREATE INDEX idx_recording_updated_at ON recordings (updated_at);

CREATE INDEX idx_recording_view_recording_uuid ON recording_views (recording_uuid);
CREATE INDEX idx_recording_view_created_at ON recording_views (created_at);

CREATE INDEX idx_user_recording_like_created_at on user_recording_likes (created_at);

-- down
DROP INDEX IF EXISTS idx_recording_updated_at;

DROP INDEX IF EXISTS idx_recording_view_recording_uuid;
DROP INDEX IF EXISTS idx_recording_view_created_at;

DROP INDEX IF EXISTS idx_user_recording_like_created_at;