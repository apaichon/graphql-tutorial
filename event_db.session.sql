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

drop table event

CREATE TABLE IF NOT EXISTS event (
  event_id INTEGER PRIMARY KEY,
  parent_event_id INTEGER,
  name TEXT,
  description TEXT,
  start_date DATE,
  end_date DATE,
  location_id INTEGER,
  created_at DATETIME,
  created_by TEXT,
  FOREIGN KEY (location_id) REFERENCES location (location_id)
);

CREATE TABLE IF NOT EXISTS ticket (
  ticket_id INTEGER PRIMARY KEY,
  type TEXT,
  price REAL,
  event_id INTEGER,
  created_at DATETIME,
  created_by TEXT,
  FOREIGN KEY (event_id) REFERENCES event (event_id)
);


-- Sample data for the 'event' table
INSERT INTO event (event_id, parent_event_id, name, description, start_date, end_date, location_id, created_at, created_by)
VALUES 
  (1, NULL, 'Music Concert', 'An evening of music and entertainment', '2024-07-01', '2024-07-01', 1, '2024-06-01 10:00:00', 'admin'),
  (2, NULL, 'Art Exhibition', 'A showcase of contemporary art', '2024-07-15', '2024-07-20', 2, '2024-06-10 11:00:00', 'admin'),
  (3, 1, 'Charity Concert', 'A concert for charity', '2024-08-01', '2024-08-01', 1, '2024-06-20 12:00:00', 'admin'),
  (4, NULL, 'Theater Play', 'A classic play performed live', '2024-09-01', '2024-09-01', 3, '2024-07-01 13:00:00', 'admin');

-- Sample data for the 'ticket' table
INSERT INTO ticket (ticket_id, type, price, event_id, created_at, created_by)
VALUES 
  (1, 'VIP', 150.00, 1, '2024-06-05 09:00:00', 'admin'),
  (2, 'Standard', 75.00, 1, '2024-06-05 09:30:00', 'admin'),
  (3, 'Early Bird', 60.00, 2, '2024-06-10 12:00:00', 'admin'),
  (4, 'Regular', 50.00, 3, '2024-06-20 14:00:00', 'admin'),
  (5, 'VIP', 200.00, 4, '2024-07-01 14:00:00', 'admin');


  select * from ticket 

   SELECT * FROM ticket
             Where type like '%%'
            LIMIT 10 OFFSET 0


      SELECT * FROM event
             Where event_id in (1,1,2,3,4)

      --delete from biding

      select * from biding

