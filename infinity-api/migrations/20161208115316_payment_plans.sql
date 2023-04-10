-- up
CREATE TABLE payment_plans (
  uuid                    UUID        NOT NULL PRIMARY KEY,
  name                    VARCHAR(30) NOT NULL UNIQUE,
  description             TEXT        NOT NULL,
  bandwidth               BIGINT      NOT NULL DEFAULT 0 CHECK (bandwidth > 0),
  per_recording_bandwidth BIGINT      NOT NULL DEFAULT 0 CHECK (per_recording_bandwidth <= bandwidth),
  devices                 SMALLINT    NOT NULL DEFAULT 1 CHECK (devices >= 1),
  duration                SMALLINT    NOT NULL,
  price                   FLOAT       NOT NULL,
  updated_at              TIMESTAMP   NOT NULL DEFAULT now(),
  created_at              TIMESTAMP   NOT NULL DEFAULT now()
);

ALTER TABLE users
  ADD COLUMN payment_plan_uuid UUID NOT NULL;

ALTER TABLE users
  ADD COLUMN payment_plan_subscribed_at TIMESTAMP NULL;

-- restrict the deletion of the payment plan until all users have been moved off it.
ALTER TABLE users
  ADD CONSTRAINT fk_payment_plan FOREIGN KEY (payment_plan_uuid) REFERENCES payment_plans (uuid) ON DELETE RESTRICT;

-- add some dummy data
INSERT INTO public.payment_plans (uuid, name, description, bandwidth, devices, duration, price, updated_at, created_at)
VALUES
  ('2f367b7d-f497-4c35-934b-53ac0cbf1c83', 'Guest plan',
   'As a guest on our platform As a guest on our platform As a guest on our platform As a guest on our platform As a guest on our platform As a guest on our platform As a guest on our platform As a guest on our platform',
   2147483648, 1, 0, 0, '2017-06-16 14:40:38.086000', '2017-06-16 14:40:39.142000'),
  ('2f367b7d-f497-4c35-934b-53ac0cbf1c23', 'Premium',
   'As a guest on our platform As a guest on our platform As a guest on our platform As a guest on our platform As a guest on our platform As a guest on our platform As a guest on our platform As a guest on our platform',
   8247483648, 4, 30, 500, '2017-06-16 15:21:59.217000', '2017-06-16 15:21:59.954000'),
  ('c37552b6-3365-478d-b328-9724aee4ae5b', 'Basic',
   'As a guest on our platform As a guest on our platform As a guest on our platform As a guest on our platform As a guest on our platform As a guest on our platform As a guest on our platform As a guest on our platform',
   4247483648, 2, 30, 100, '2017-06-16 14:41:34.957000', '2017-06-16 14:41:35.876000');

-- down
ALTER TABLE users
  DROP CONSTRAINT fk_payment_plan;

ALTER TABLE users
  DROP COLUMN payment_plan_uuid;

DROP TABLE payment_plans;