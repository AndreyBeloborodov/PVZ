-- +goose Up
-- +goose StatementBegin
ALTER TABLE packages
ADD CONSTRAINT packages_name_unique UNIQUE (name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE packages
DROP CONSTRAINT IF EXISTS packages_name_unique;
-- +goose StatementEnd
