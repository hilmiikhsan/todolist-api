-- +goose Up
-- +goose StatementBegin
CREATE TABLE todos
(
    todo_id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    activity_group_id INTEGER NOT NULL,
    title VARCHAR(100) NOT NULL,
    priority VARCHAR(100) NOT NULL,
    is_active BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE todos;
-- +goose StatementEnd
