
-- Lab 4.1

CREATE TABLE IF NOT EXISTS contact (
  contact_id INTEGER PRIMARY KEY,
  name TEXT,
  first_name TEXT,
  last_name TEXT,
  gender_id INTEGER,
  dob DATE,
  email TEXT,
  phone TEXT,
  address TEXT,
  photo_path TEXT,
  created_at DATETIME,
  created_by TEXT
);