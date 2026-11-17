package dao_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/a-novel-kit/golib/postgres"

	"github.com/a-novel/service-template/internal/dao"
)

var (
	// itemID names the oldest item [insertItems] seeds.
	itemID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	// newItemID is the ID a write offers a new item.
	newItemID = uuid.MustParse("00000000-0000-0000-0000-000000000004")
	// absentID names no item.
	absentID = uuid.MustParse("00000000-0000-0000-0000-0000000000ff")

	createdAt = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	now       = time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC)
)

// itemFixtures returns the items every DAO test can seed: three items, created a day apart, oldest
// first. Each call builds new values, so an insert scanning rows back into them leaves the other
// parallel tests' copies alone.
func itemFixtures() []*dao.Item {
	return []*dao.Item{
		{ID: itemID, Name: "item one", Description: "description one", CreatedAt: createdAt, UpdatedAt: createdAt},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			Name:      "item two",
			CreatedAt: createdAt.AddDate(0, 0, 1),
			UpdatedAt: createdAt.AddDate(0, 0, 1),
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			Name:      "item three",
			CreatedAt: createdAt.AddDate(0, 0, 2),
			UpdatedAt: createdAt.AddDate(0, 0, 2),
		},
	}
}

// insertItems seeds the database of ctx with [itemFixtures] and returns them.
func insertItems(ctx context.Context, t *testing.T) []*dao.Item {
	t.Helper()

	db, err := postgres.GetContext(ctx)
	require.NoError(t, err)

	fixtures := itemFixtures()

	_, err = db.NewInsert().Model(&fixtures).Exec(ctx)
	require.NoError(t, err)

	return fixtures
}
