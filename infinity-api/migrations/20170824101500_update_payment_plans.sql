-- up
UPDATE payment_plans
SET
  bandwidth   = 50000000,
  name        = 'Guest',
  devices     = 1,
  description = 'As a guest you can view 50 MB each day for free'
WHERE uuid = '2f367b7d-f497-4c35-934b-53ac0cbf1c83';

UPDATE payment_plans
SET
  bandwidth   = 200000000,
  price       = 0,
  devices     = 1,
  description = 'As a Basic user you can view 250 MB each day for free and build your own collection of your favourite performers. When upgrading to a premium plan, your collection remains valid.'
WHERE uuid = 'c37552b6-3365-478d-b328-9724aee4ae5b';

UPDATE payment_plans
SET
  bandwidth   = 15000000000,
  name        = 'Premium 30',
  price       = 15,
  devices     = 1,
  description = 'As a premium user you can view 15 GB each day and build your own collection of your favourite performers. When upgrading to a premium plan, your collection remains valid.'
WHERE uuid = '2f367b7d-f497-4c35-934b-53ac0cbf1c23';

INSERT INTO payment_plans (uuid, name, description, bandwidth, devices, duration, price, updated_at, created_at)
VALUES
  (
    uuid_generate_v4(),
    'Premium 90',
    'As a premium user you can view 15 GB each day and build your own collection of your favourite performers. When upgrading to a premium plan, your collection remains valid.',
    15000000000,
    1,
    90,
    35,
    '2017-06-16 14:40:38.086000',
    '2017-06-16 14:40:39.142000'
  ),
  (
    uuid_generate_v4(),
    'Premium 180',
    'As a premium user you can view 15 GB each day and build your own collection of your favourite performers. When upgrading to a premium plan, your collection remains valid.',
    15000000000,
    1,
    180,
    60,
    '2017-06-16 14:40:38.086000',
    '2017-06-16 14:40:39.142000'
  ),
  (
    uuid_generate_v4(),
    'Premium 360',
    'As a premium user you can view 15 GB each day and build your own collection of your favourite performers. When upgrading to a premium plan, your collection remains valid.',
    15000000000,
    1,
    360,
    105,
    '2017-06-16 14:40:38.086000',
    '2017-06-16 14:40:39.142000'
  ),
  (
    uuid_generate_v4(),
    'Premium 720',
    'As a premium user you can view 15 GB each day and build your own collection of your favourite performers. When upgrading to a premium plan, your collection remains valid.',
    15000000000,
    1,
    720,
    190,
    '2017-06-16 14:40:38.086000',
    '2017-06-16 14:40:39.142000'
  );

-- down