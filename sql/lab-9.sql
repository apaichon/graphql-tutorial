CREATE TABLE biding_room (
    room_id INTEGER PRIMARY KEY AUTOINCREMENT,
    start_date DATETIME NOT NULL,
    end_date DATETIME NOT NULL,
    product_name TEXT NOT NULL,
    floor_price REAL NOT NULL,
    product_image TEXT
);



CREATE TABLE biding (
    bid_id INTEGER PRIMARY KEY AUTOINCREMENT,
    room_id INTEGER,
    bidder TEXT NOT NULL,
    bid_price REAL NOT NULL,
    bid_time DATETIME NOT NULL,
    FOREIGN KEY (room_id) REFERENCES biding_roomRoom (room_id)
);

--drop table biding

// BidingRoom represents the biding_room table
type BidingRoom struct {
    RoomID      int       `db:"room_id"`
    StartDate   time.Time `db:"start_date"`
    EndDate     time.Time `db:"end_date"`
    ProductName string    `db:"product_name"`
    FloorPrice  float64   `db:"floor_price"`
    ProductImage string   `db:"product_image"`
}

// Biding represents the biding table
type Biding struct {
    BidID    int       `db:"bid_id"`
    RoomID   int       `db:"room_id"`
    Bidder   string    `db:"bidder"`
    BidPrice float64   `db:"bid_price"`
    BidTime  time.Time `db:"bid_time"`
}

INSERT INTO biding_room (start_date, end_date, product_name, floor_price, product_image)
VALUES ('2024-05-23', '2024-06-20', 'LaFerrari 2013', 20000000.00, 'static/images/ferrari-laferrari-2013-chassis.png');

select * from biding_room

SELECT * FROM biding_room WHERE room_id = 1