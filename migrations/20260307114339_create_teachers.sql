-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS teachers 
(
    user_id integer unique not null references users(id) on delete cascade,

    department text not null,

    created_at timestamptz default now(),
    updated_at timestamptz default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS teachers;
-- +goose StatementEnd
