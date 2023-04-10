-- up
CREATE TABLE non_deletable_recordings (
  url VARCHAR(255) NOT NULL,
  filename VARCHAR(255) NOT NULL,
  PRIMARY KEY (url, filename)
);

-- down