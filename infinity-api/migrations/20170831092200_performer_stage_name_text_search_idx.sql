-- up
ALTER TABLE performers ALTER stage_name TYPE TEXT;
CREATE INDEX idx_performers_text_search_stage_name ON performers USING GIN (TO_TSVECTOR('english', stage_name));

-- down
DROP INDEX IF EXISTS idx_performers_text_search_stage_name;
ALTER TABLE performers ALTER stage_name TYPE VARCHAR(40);
