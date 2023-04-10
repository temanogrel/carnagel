-- up
CREATE INDEX idx_recording_last_validation_at ON recordings (last_validation_at);

-- down
DROP INDEX IF EXISTS idx_recording_last_validation_at;