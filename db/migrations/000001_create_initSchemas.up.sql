-- CREATE TYPE currency AS ENUM (
--   'USD',
--   'IDR'
-- );

CREATE TABLE accounts (
  id bigserial PRIMARY KEY,
  owner varchar NOT NULL,
  balance bigint NOT NULL DEFAULT 0,
  currency varchar(5) NOT NULL,
  created_at timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE entries (
  id bigserial PRIMARY KEY,
  account_id bigint NOT NULL,
  amount bigint NOT NULL,
  description varchar,
  created_at timestamptz NOT NULL DEFAULT (now()),

  CONSTRAINT transaction_account_id_fk FOREIGN KEY (account_id) REFERENCES accounts (id)
);

CREATE TABLE transfers (
  id bigserial PRIMARY KEY,
  from_account_id bigint NOT NULL,
  to_account_id bigint NOT NULL,
  amount bigint NOT NULL,
  created_at timestamptz NOT NULL DEFAULT (now()),

  CONSTRAINT transfer_from_account_id_fk FOREIGN KEY (from_account_id) REFERENCES accounts(id),
  CONSTRAINT transfer_to_account_id_fk FOREIGN KEY (to_account_id) REFERENCES accounts(id)
);

CREATE INDEX transaction_account_idx ON entries (account_id);
CREATE INDEX account_owner_idx ON accounts (owner);
CREATE INDEX transfer_from_idx ON transfers (from_account_id);
CREATE INDEX transfer_to_idx ON transfers (to_account_id);
CREATE INDEX transfer_from_to_idx ON transfers (from_account_id, to_account_id);

COMMENT ON COLUMN "entries"."amount" IS 'can be negative or positive';
COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';
