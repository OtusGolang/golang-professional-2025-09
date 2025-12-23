-- +goose Up
-- +goose StatementBegin
create table pages (
  id int not null,
  count int not null,
  created_at      timestamptz not null default now()
);

INSERT INTO pages (id, count)
VALUES
    (1, 100),
    (2, 200),
    (3, 300);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table pages;
-- +goose StatementEnd
