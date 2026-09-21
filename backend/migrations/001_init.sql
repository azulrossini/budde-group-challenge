-- +goose Up
CREATE TABLE bills (
  id           BIGSERIAL PRIMARY KEY,
  description  TEXT   NOT NULL CHECK (length(trim(description)) > 0),
  total_cents  BIGINT NOT NULL CHECK (total_cents >= 0),
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE shares (
  id             BIGSERIAL PRIMARY KEY,
  bill_id        BIGINT  NOT NULL REFERENCES bills(id) ON DELETE CASCADE,
  person_name    TEXT    NOT NULL CHECK (length(trim(person_name)) > 0),
  percentage_bp  INTEGER NOT NULL CHECK (percentage_bp > 0 AND percentage_bp <= 10000)
);

CREATE UNIQUE INDEX shares_bill_person_unique ON shares (bill_id, lower(person_name));

-- +goose Down
DROP TABLE shares;
DROP TABLE bills;
