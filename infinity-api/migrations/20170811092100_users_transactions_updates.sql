-- up
ALTER TABLE users ADD payment_plan_ends_at TIMESTAMP NULL;

ALTER TABLE payment_transactions ADD unconfirmed_received_amount_in_satoshis BIGINT;
ALTER TABLE payment_transactions ADD confirmed_fully_paid BOOLEAN;

-- down
ALTER TABLE users DROP payment_plan_ends_at;

ALTER TABLE payment_transactions DROP unconfirmed_received_amount_in_satoshis;
ALTER TABLE payment_transactions DROP confirmed_fully_paid;