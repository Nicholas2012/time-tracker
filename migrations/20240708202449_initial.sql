-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION pg_trgm;
CREATE TABLE users (
                    id SERIAL PRIMARY KEY,
                    name VARCHAR NOT NULL,
                    surname VARCHAR NOT NULL,
                    patronymic VARCHAR NOT NULL,
                    passport_serie INT NOT NULL,
                    passport_number INT NOT NULL
);
CREATE TABLE tasks (
                    id SERIAL PRIMARY KEY,
		            user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                    start_time timestamptz NOT NULL,
                    end_time timestamptz NOT NULL,
                    minutes INT NOT NULL
);
--CREATE INDEX users_names ON users USING gin ((name || ' ' || surname || ' ' || patronymic) gin_trgm_ops);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP TABLE tasks;
-- +goose StatementEnd
