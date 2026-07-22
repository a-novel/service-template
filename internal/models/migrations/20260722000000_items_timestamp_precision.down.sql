-- Narrowing back to second precision truncates any sub-second component already stored; the rows
-- survive but their exact write instants do not. That is inherent to reverting this change.
ALTER TABLE items
ALTER COLUMN created_at TYPE timestamp(0) with time zone,
ALTER COLUMN updated_at TYPE timestamp(0) with time zone,
ALTER COLUMN created_at
SET DEFAULT CURRENT_TIMESTAMP(0),
ALTER COLUMN updated_at
SET DEFAULT CURRENT_TIMESTAMP(0);
