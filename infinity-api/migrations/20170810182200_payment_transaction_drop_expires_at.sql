-- up
ALTER TABLE payment_transactions DROP expires_at;

--down
ALTER TABLE payment_transactions ADD expires_at TIMESTAMP;