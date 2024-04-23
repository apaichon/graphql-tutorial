-- Lab3-1
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

INSERT INTO contact (contact_id, name, first_name, last_name, gender_id, dob, email, phone, address, photo_path, created_at, created_by)
VALUES 
(1, 'John Doe', 'John', 'Doe', 1, '1990-01-01', 'john@example.com', '123-456-7890', '123 Main St', 'path/to/photo.jpg', '2024-04-16 12:00:00', 'Admin'),
(2, 'Jane Smith', 'Jane', 'Smith', 2, '1992-05-15', 'jane@example.com', '987-654-3210', '456 Oak St', 'path/to/photo2.jpg', '2024-04-16 12:30:00', 'Admin');
