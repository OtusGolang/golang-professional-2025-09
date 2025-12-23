-- +goose Up
-- +goose StatementBegin
create index idx_books_description
    on books (description);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index idx_books_description;
-- +goose StatementEnd
