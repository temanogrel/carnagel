-- up
ALTER TABLE recordings ADD last_validation_at TIMESTAMP DEFAULT NULL;

-- down
ALTER TABLE recordings DROP last_validation_at;