-- +goose Up
-- +goose StatementBegin
CREATE TABLE urls (
  id text primary key,
  original_url text not null,
  created_at timestamp not null default (current_timestamp),
  visit_count integer not null default 0,
  last_visited_at timestamp null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE urls;
-- +goose StatementEnd
