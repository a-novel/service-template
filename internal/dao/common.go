package dao

import "errors"

// ErrItemNotFound is returned when no item has the requested ID. A DAO joins it with
// [database/sql.ErrNoRows], so callers keep both identities.
var ErrItemNotFound = errors.New("item not found")
