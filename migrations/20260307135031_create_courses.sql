-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS courses
(
    id integer primary key generated always as identity,

    teacher_id integer not null references teachers(user_id) on delete cascade,
    
    title text not null,
    description text not null,

    created_at timestamptz default now(),
    updated_at timestamptz default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS courses;
-- +goose StatementEnd
