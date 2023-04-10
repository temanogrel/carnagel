-- up
ALTER TABLE recordings ADD slug VARCHAR(255) DEFAULT NULL;

CREATE UNIQUE INDEX idx_recording_unique_slug ON recordings (slug);

-- down
ALTER TABLE recordings DROP slug;

DROP INDEX IF EXISTS idx_recording_unique_slug;