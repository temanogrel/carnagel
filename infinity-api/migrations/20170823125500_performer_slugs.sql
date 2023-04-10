-- up
ALTER TABLE performers ADD slug VARCHAR(255) DEFAULT NULL;

CREATE UNIQUE INDEX idx_performer_unique_slug ON performers (slug);

-- down
ALTER TABLE performers DROP slug;

DROP INDEX IF EXISTS idx_performer_unique_slug;