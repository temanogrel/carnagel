-- up
ALTER TABLE files
  ADD checksum VARCHAR(80) DEFAULT NULL;

-- down
ALTER TABLE files
  DROP checksum
