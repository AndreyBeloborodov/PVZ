-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders
ADD COLUMN price INTEGER,
ADD COLUMN weight INTEGER,
ADD COLUMN name_package TEXT REFERENCES packages(name),
ADD COLUMN result_price INTEGER;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders
DROP COLUMN IF EXISTS price,
DROP COLUMN IF EXISTS weight,
DROP COLUMN IF EXISTS name_package,
DROP COLUMN IF EXISTS result_price;
-- +goose StatementEnd
