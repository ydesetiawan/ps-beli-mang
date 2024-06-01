CREATE INDEX IF NOT EXISTS merchants_idx_category ON merchants (merchant_category);
CREATE INDEX IF NOT EXISTS merchants_idx_name ON merchants (lower(name));
CREATE INDEX IF NOT EXISTS merchants_idx_created_at ON merchants (created_at);
CREATE INDEX IF NOT EXISTS merchant_items_idx_merchant_id ON merchant_items (merchant_id);
CREATE INDEX IF NOT EXISTS merchant_items_idx_name ON merchant_items (lower(name));
CREATE INDEX IF NOT EXISTS merchant_items_idx_category ON merchant_items (category);
CREATE INDEX IF NOT EXISTS merchant_items_idx_created_at ON merchant_items (created_at);