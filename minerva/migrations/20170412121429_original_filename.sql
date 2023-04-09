-- up
ALTER TABLE files
  ADD original_filename VARCHAR(70) DEFAULT NULL;

-- down
ALTER TABLE files
  DROP original_filename;