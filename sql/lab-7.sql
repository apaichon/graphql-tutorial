CREATE TABLE IF NOT EXISTS logs (
    log_id TEXT PRIMARY KEY,
    timestamp DATETIME,
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
    duration INTEGER, -- Using INTEGER to store duration in nanoseconds
    errors TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

delete from logs

select count(*) from logs

select * from logs 
