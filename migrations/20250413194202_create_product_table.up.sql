CREATE TABLE IF NOT EXISTS product (
    id UUID PRIMARY KEY,
    date_time TIMESTAMP NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('электроника', 'одежда', 'обувь')),
    reception_id UUID NOT NULL REFERENCES reception(id) ON DELETE CASCADE
);