-- up
CREATE INDEX idx_recording_lowercase_stage_name ON recordings (LOWER(stage_name));
ALTER TABLE recordings ALTER stage_name TYPE TEXT;
CREATE INDEX idx_recording_text_search_stage_name ON recordings USING GIN (TO_TSVECTOR('english', stage_name));

-- down
DROP INDEX IF EXISTS idx_recording_lowercase_stage_name;
DROP INDEX IF EXISTS idx_recording_text_search_stage_name;
ALTER TABLE recordings ALTER stage_name TYPE VARCHAR(40);