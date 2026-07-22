-- id breaks ties on created_at. Without it the ordering is partial, and LIMIT/OFFSET over a partial
-- order lets Postgres return a tied group in a different sequence per query, so paging silently
-- repeats some rows and skips others.
SELECT
  *
FROM
  items
ORDER BY
  created_at DESC,
  id DESC
LIMIT
  ?0
OFFSET
  ?1;
