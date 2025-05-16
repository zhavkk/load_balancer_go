-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS client_limits (
    client_id       TEXT    PRIMARY KEY,
    req_per_sec     INTEGER NOT NULL,
    burst_capacity  INTEGER NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS client_limits;
-- +goose StatementEnd
