CREATE TABLE IF NOT EXISTS reception (
    id UUID PRIMARY KEY,
    date_time TIMESTAMP NOT NULL,
    pvz_id UUID NOT NULL REFERENCES pvz(id) ON DELETE CASCADE,
    status TEXT NOT NULL CHECK (status IN ('in_progress', 'close'))
)