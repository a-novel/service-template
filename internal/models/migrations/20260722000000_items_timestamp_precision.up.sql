-- Widen the item timestamps from second to full precision.
--
-- timestamp(0) truncates to whole seconds, so any two rows written in the same second compare equal.
-- That made created_at unusable as a sort key on its own, and left updated_at unable to distinguish
-- two updates within one second — which forecloses using it for optimistic concurrency later.
ALTER TABLE items
  ALTER COLUMN created_at TYPE timestamp with time zone,
  ALTER COLUMN updated_at TYPE timestamp with time zone,
  ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP,
  ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP;
