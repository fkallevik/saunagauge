CREATE TABLE dht_data (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  temperature REAL DEFAULT 0,
  humidity REAL DEFAULT 0,
  created_at TEXT DEFAULT NULL
);