-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS packages (
                                      id BIGSERIAL PRIMARY KEY,
                                      name TEXT NOT NULL,
                                      max_weight INTEGER NOT NULL,
                                      price INTEGER NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS packages;
-- +goose StatementEnd
