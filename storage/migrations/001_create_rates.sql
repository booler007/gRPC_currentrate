CREATE TABLE IF NOT EXISTS rates (
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp INTEGER NOT NULL,
    best_ask  TEXT    NOT NULL,
    best_bid  TEXT    NOT NULL
);
