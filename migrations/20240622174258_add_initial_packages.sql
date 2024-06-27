-- +goose Up
-- +goose StatementBegin
INSERT INTO packages (name, max_weight, price) VALUES
                                                   ('package', 10, 5),
                                                   ('box', 30, 20),
                                                   ('film', -1, 1);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM packages WHERE name IN ('package', 'box', 'film');
-- +goose StatementEnd
