-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS enrollments 
(
    student_id integer not null references students(user_id) on delete cascade,
    course_id integer not null references courses(id) on delete cascade,
    primary key (student_id, course_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS enrollments;
-- +goose StatementEnd
