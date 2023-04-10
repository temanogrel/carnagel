-- up
ALTER TABLE payment_transactions DROP plutus_uuid;
ALTER TABLE payment_transactions DROP latest_plutus_status_code;

ALTER TABLE payment_transactions RENAME received_amount TO received_amount_in_satoshis;

ALTER TABLE payment_transactions ADD expected_amount_in_satoshis BIGINT;
ALTER TABLE payment_transactions ADD conversion_rate FLOAT;
ALTER TABLE payment_transactions ADD webhook_uuid UUID UNIQUE;
ALTER TABLE payment_transactions ADD payment_address VARCHAR(100) UNIQUE;

-- down
ALTER TABLE payment_transactions DROP webhook_uuid;
ALTER TABLE payment_transactions DROP expected_amount_in_satoshis;
ALTER TABLE payment_transactions DROP conversion_rate;
ALTER TABLE payment_transactions DROP payment_address;

ALTER TABLE payment_transactions RENAME received_amount_in_satoshis TO received_amount;

ALTER TABLE payment_transactions ADD latest_plutus_status_code INTEGER;
ALTER TABLE payment_transactions ADD plutus_uuid UUID UNIQUE;
