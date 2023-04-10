-- up
ALTER TABLE non_deletable_recordings RENAME url TO upstore_hash;
CREATE INDEX idx_non_deletable_recordings_upstore_hash ON non_deletable_recordings (upstore_hash);

-- down