-- up
CREATE TABLE payment_transactions (
  uuid                      UUID      NOT NULL PRIMARY KEY,
  plutus_uuid               UUID      NOT NULL UNIQUE,
  state                     SMALLINT  NOT NULL,
  latest_plutus_status_code INTEGER,
  received_amount           BIGINT,
  expires_at                TIMESTAMP NOT NULL,
  updated_at                TIMESTAMP NOT NULL DEFAULT now(),
  created_at                TIMESTAMP NOT NULL DEFAULT now(),
  payment_plan_uuid         UUID      NOT NULL,
  user_uuid                 UUID      NOT NULL
);

-- add foreign key to users table
ALTER TABLE payment_transactions
  ADD CONSTRAINT fk_user FOREIGN KEY (user_uuid) REFERENCES users (uuid) ON DELETE RESTRICT;

-- add foreign key to payment_plans table
ALTER TABLE payment_transactions
  ADD CONSTRAINT fk_payment_plan FOREIGN KEY (payment_plan_uuid) REFERENCES payment_plans (uuid) ON DELETE RESTRICT;

-- down
ALTER TABLE payment_transactions
  DROP CONSTRAINT fk_user;

ALTER TABLE payment_transactions
  DROP CONSTRAINT fk_payment_plan;

DROP TABLE payment_transactions;