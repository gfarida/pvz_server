CREATE INDEX IF NOT EXISTS idx_reception_pvz_status ON reception(pvz_id, status);

CREATE INDEX IF NOT EXISTS idx_reception_date_time ON reception(date_time);

CREATE INDEX IF NOT EXISTS idx_product_reception_id ON product(reception_id);
