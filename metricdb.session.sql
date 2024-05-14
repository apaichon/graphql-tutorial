

CREATE TABLE IF NOT EXISTS logs (
    log_id TEXT,
    timestamp TIMESTAMPTZ NOT NULL,
    user_id INTEGER,
    action TEXT,
    resource TEXT,
    status TEXT,
    client_ip TEXT,
    client_device TEXT,
    client_os TEXT,
    client_os_ver TEXT,
    client_browser TEXT,
    client_browser_ver TEXT,
    duration INTERVAL, -- Using INTERVAL to store duration
    errors TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);


select * from logs


-- Create a time-series hypertable partitioned by timestamp
SELECT create_hypertable('logs', 'timestamp');


CREATE OR REPLACE VIEW vw_log_status_summary
AS 
SELECT
timestamp,
    action,
    CASE WHEN status = 'OK' and errors = '' THEN 1 ELSE  0 END AS success,
    CASE WHEN errors IS NOT NULL AND errors != '' THEN 1 ELSE 0 END error
FROM
    logs


   select * from  vw_log_status_summary
   order by timestamp desc



SELECT 
 DATE_TRUNC('minute', timestamp) + (((EXTRACT(MINUTE FROM timestamp)::int / 5) * 5) || ' minutes')::interval AS timestamp
, SUM(success) as total_successes, SUM(error) as total_errors FROM vw_log_status_summary GROUP BY 
 DATE_TRUNC('minute', timestamp) + (((EXTRACT(MINUTE FROM timestamp)::int / 5) * 5) || ' minutes')::interval
order by timestamp desc
LIMIT 50 