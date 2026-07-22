CREATE TABLE items (
  id uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
  name text NOT NULL CHECK (name <> ''),
  description text,
  -- Full precision, not timestamp(0): second precision makes rows written in the same second
  -- compare equal, which leaves created_at unusable as a sort key on its own and stops updated_at
  -- from distinguishing two updates within a second.
  created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);
