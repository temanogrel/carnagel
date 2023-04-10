-- up
CREATE INDEX idx_payment_transaction_user_uuid ON payment_transactions (user_uuid);

-- down
DROP INDEX IF EXISTS idx_payment_transaction_user_uuid;