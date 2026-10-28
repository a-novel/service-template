-- Item IDs are UUIDv7, which sort by creation time, so the primary key serves this order.
SELECT
  *,
  count(*) OVER () AS total
FROM
  items
ORDER BY
  id DESC
LIMIT
  ?0
OFFSET
  ?1;
