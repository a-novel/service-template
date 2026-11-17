package core_test

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/a-novel/service-template/internal/core"
	"github.com/a-novel/service-template/internal/dao"
)

var (
	// errFoo is the unexpected error a DAO mock returns, which every service passes on.
	errFoo = errors.New("foo")

	itemID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

	createdAt = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedAt = time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)

	// daoItem is the item the DAO mocks return.
	daoItem = &dao.Item{
		ID: itemID, Name: "test item", Description: "test description", CreatedAt: createdAt, UpdatedAt: updatedAt,
	}
	// coreItem is daoItem as the services return it.
	coreItem = &core.Item{
		ID: itemID, Name: "test item", Description: "test description", CreatedAt: createdAt, UpdatedAt: updatedAt,
	}
)

// daoMock is what the DAO mock of a service test case returns. A nil one means the service must not
// call the DAO.
type daoMock[Resp any] struct {
	resp Resp
	err  error
}
