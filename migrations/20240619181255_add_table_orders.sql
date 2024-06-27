-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
                                         id BIGSERIAL PRIMARY KEY,
                                         user_id BIGSERIAL NOT NULL,
                                         time_end TIMESTAMP NOT NULL,
                                         is_given BOOLEAN,
                                         time_given TIMESTAMP,
                                         is_returned BOOLEAN
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
