-- up
CREATE INDEX idx_performer_lowercase_stage_name ON performers (lower(stage_name));

-- down
DROP INDEX IF EXISTS idx_performer_lowercase_stage_name;