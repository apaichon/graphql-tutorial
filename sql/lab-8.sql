{
  "log_id": "018f5bce-a120-7c2d-97d7-5efe0c30ee45",
  "timestamp": "2024-05-09 12:23:39.424767+07:00",
  "user_id": 2,
  "action": "contactMutations.updateContact",
  "resource": "GraphQLApi",
  "status": "OK",
  "client_ip": "127.0.0.1:55083",
  "client_device": "",
  "client_os": "",
  "client_os_ver": "",
  "client_browser": "https:",
  "client_browser_ver": "/www.thunderclient.com",
  "duration": 743833,
  "errors": "Message: unauthorized: missing contactMutations.updateContact permission | Locations: Line 3, Column 2 | Path: contactMutations.updateContact",
  "created_at": "2024-05-09 05:25:06"
}


CREATE TABLE IF NOT EXISTS logs (
    log_id varchar(50),
    timestamp TIMESTAMPTZ NOT NULL,
    user_id INTEGER,
    action varchar(100),
    resource varchar(50),
    status varchar(50),
    client_ip varchar(50),
    client_device varchar(50),
    client_os varchar(50),
    client_os_ver varchar(50),
    client_browser varchar(50),
    client_browser_ver varchar(50),
    duration INTERVAL, -- Using INTERVAL to store duration
    errors varchar(50),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create a time-series hypertable partitioned by timestamp
SELECT create_hypertable('logs', 'timestamp');


CREATE VIEW analytics_view AS
SELECT
    DATE(timestamp) AS date,
    COUNT(*) AS total_requests,
    SUM(CASE WHEN status = 'ERROR' THEN 1 ELSE 0 END) AS total_errors,
    AVG(duration) AS avg_duration,
    MIN(duration) AS min_duration,
    MAX(duration) AS max_duration,
    FROM_UNIXTIME(UNIX_TIMESTAMP(MIN(timestamp) + INTERVAL duration SECOND), '%H:%i:%s') AS peak_time,
    client_os AS favorite_os,
    client_device AS favorite_device
FROM
    logs
GROUP BY
    DATE(timestamp);


    SELECT
    DATE(timestamp) AS date,
    COUNT(*) AS total_requests,
    SUM(CASE WHEN status = 'ERROR' THEN 1 ELSE 0 END) AS total_errors,
    AVG(duration) AS avg_duration,
    MIN(duration) AS min_duration,
    MAX(duration) AS max_duration,
    strftime('%H:%M:%S', timestamp + (duration / 1000000000), 'unixepoch') AS peak_time,
    client_os AS favorite_os,
    client_device AS favorite_device
FROM
    logs
GROUP BY
    DATE(timestamp);


    CREATE VIEW analytics_view AS
SELECT
    strftime('%Y-%m-%d %H:', timestamp) || (CAST(strftime('%M', timestamp) AS INTEGER) / 5) * 5 AS interval,
    COUNT(*) AS total_requests,
    SUM(CASE WHEN status = 'ERROR' THEN 1 ELSE 0 END) AS total_errors,
    AVG(duration) AS avg_duration,
    MIN(duration) AS min_duration,
    MAX(duration) AS max_duration,
    strftime('%H:%M:%S', MIN(timestamp) + (AVG(duration) / 1000000000), 'unixepoch') AS peak_time,
    client_os AS favorite_os,
    client_device AS favorite_device
FROM
    logs
    where timestamp between '' and ''
GROUP BY
    strftime('%Y-%m-%d %H:', timestamp) || (CAST(strftime('%M', timestamp) AS INTEGER) / 5) * 5,
    client_os,
    client_device;


    select 5002390875/1_000_000_000


    select * from logs