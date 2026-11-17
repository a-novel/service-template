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
	// itemID names itemOne, the oldest item insertItems seeds.
	itemID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	// newItemID is the ID a write offers a new item.
	newItemID = uuid.MustParse("00000000-0000-0000-0000-000000000004")
	// absentID names no item.
	absentID = uuid.MustParse("00000000-0000-0000-0000-0000000000ff")

	createdAt = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	now       = time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC)
)

// itemOne, itemTwo and itemThree are the items insertItems seeds, created a day apart.
var (
	itemOne = dao.Item{
		ID: itemID, Name: "item one", Description: "description one", CreatedAt: createdAt, UpdatedAt: createdAt,
	}
	itemTwo = dao.Item{
		ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		Name:      "item two",
		CreatedAt: createdAt.AddDate(0, 0, 1),
		UpdatedAt: createdAt.AddDate(0, 0, 1),
	}
	itemThree = dao.Item{
		ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
		Name:      "item three",
		CreatedAt: createdAt.AddDate(0, 0, 2),
		UpdatedAt: createdAt.AddDate(0, 0, 2),
	}
)

// insertItems seeds the database of ctx with copies of itemOne, itemTwo and itemThree: the insert
// scans the rows back into what it is given, which parallel tests must not share.
func insertItems(ctx context.Context, t *testing.T) {
	t.Helper()

	db, err := postgres.GetContext(ctx)
	require.NoError(t, err)

	fixtures := []*dao.Item{new(itemOne), new(itemTwo), new(itemThree)}

	_, err = db.NewInsert().Model(&fixtures).Exec(ctx)
	require.NoError(t, err)
}
