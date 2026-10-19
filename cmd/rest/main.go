// Command rest serves the item API over HTTP. It is one of the service's
// entrypoints; cmd/grpc serves the same domain over gRPC and cmd/migrations
// applies the database schema.
package main

import (
	"context"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/samber/lo"

	"github.com/a-novel-kit/golib/downtime"
	"github.com/a-novel-kit/golib/httpf"
	"github.com/a-novel-kit/golib/otel"
	"github.com/a-novel-kit/golib/postgres"

	"github.com/a-novel/service-template/internal/config"
	"github.com/a-novel/service-template/internal/config/env"
	"github.com/a-novel/service-template/internal/core"
	"github.com/a-novel/service-template/internal/dao"
	"github.com/a-novel/service-template/internal/handlers"
)

func main() {
	cfg := config.AppPresetDefault

	processCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ctx := processCtx

	otel.SetAppName(cfg.App.Name)

	lo.Must0(otel.Init(cfg.Otel))
	defer cfg.Otel.Flush()

	if env.GcloudProjectId == "" {
		log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	}

	// A server started during a planned downtime refuses work, so it leaves the database alone.
	if !downtime.Started(cfg.App.DowntimeStart, time.Now()) {
		ctx = lo.Must(postgres.NewContext(ctx, cfg.Postgres))

		database := lo.Must(cfg.Postgres.DB(ctx))
		defer closeDatabase(database)
	}

	// =================================================================================================================
	// DAO
	// =================================================================================================================

	daoItemCreate := dao.NewPgItemCreate()
	daoItemGet := dao.NewPgItemGet()
	daoItemList := dao.NewPgItemList()
	daoItemUpdate := dao.NewPgItemUpdate()
	daoItemDelete := dao.NewPgItemDelete()

	// =================================================================================================================
	// SERVICES
	// =================================================================================================================

	serviceItemCreate := core.NewItemCreate(daoItemCreate)
	serviceItemGet := core.NewItemGet(daoItemGet)
	serviceItemList := core.NewItemList(daoItemList)
	serviceItemUpdate := core.NewItemUpdate(daoItemUpdate)
	serviceItemDelete := core.NewItemDelete(daoItemDelete)

	// =================================================================================================================
	// HANDLERS
	// =================================================================================================================

	handlerPing := handlers.NewRestPing()
	handlerHealth := handlers.NewRestHealth()
	handlerItemCreate := handlers.NewRestItemCreate(serviceItemCreate, cfg.Logger)
	handlerItemGet := handlers.NewRestItemGet(serviceItemGet, cfg.Logger)
	handlerItemList := handlers.NewRestItemList(serviceItemList, cfg.Logger)
	handlerItemUpdate := handlers.NewRestItemUpdate(serviceItemUpdate, cfg.Logger)
	handlerItemDelete := handlers.NewRestItemDelete(serviceItemDelete, cfg.Logger)

	// =================================================================================================================
	// ROUTER
	// =================================================================================================================

	router := httpf.NewRouter(cfg.Rest, cfg.Otel, cfg.RestLogger, cfg.App.DowntimeStart)

	router.Get("/ping", handlerPing.ServeHTTP)
	router.Get("/healthcheck", handlerHealth.ServeHTTP)
	router.Post("/items", handlerItemCreate.ServeHTTP)
	router.Get("/items", handlerItemList.ServeHTTP)
	router.Get("/item", handlerItemGet.ServeHTTP)
	router.Put("/item", handlerItemUpdate.ServeHTTP)
	router.Delete("/item", handlerItemDelete.ServeHTTP)

	// =================================================================================================================
	// RUN
	// =================================================================================================================

	log.Println("Starting REST server on :" + strconv.Itoa(cfg.Rest.Port))

	err := httpf.ListenAndServe(ctx, cfg.Rest, router)

	stop()

	if err != nil {
		panic(err)
	}
}

func closeDatabase(database io.Closer) {
	err := database.Close()
	if err != nil {
		log.Println("Close Postgres: " + err.Error())
	}
}
