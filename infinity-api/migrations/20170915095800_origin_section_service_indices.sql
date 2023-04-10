-- up
CREATE INDEX idx_performers_origin_section ON performers (origin_section);
CREATE INDEX idx_performers_origin_service ON performers (origin_service);

-- down
DROP INDEX IF EXISTS idx_performers_origin_section;
DROP INDEX IF EXISTS idx_performers_origin_service;