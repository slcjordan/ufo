-- +goose Up
-- +goose StatementBegin
CREATE TABLE sighting (
  datetime timestamp without time zone,
  city text,
  country text,
  shape text,
  duration_in_seconds integer,
  duration text,
  posted timestamp without time zone,
  latitude numeric(9,7),
  longitude numeric(10,7),
  comments text,
  state text,
  id serial primary key
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
