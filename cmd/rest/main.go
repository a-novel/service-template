// Command rest serves the item API over HTTP. It is one of the service's
// entrypoints; cmd/grpc serves the same domain over gRPC and cmd/migrations
// applies the database schema.
package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
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

	daoItemCreate := dao.NewItemCreate()
	daoItemGet := dao.NewItemGet()
	daoItemList := dao.NewItemList()
	daoItemUpdate := dao.NewItemUpdate()
	daoItemDelete := dao.NewItemDelete()

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

	handlerPing := handlers.NewPing()
	handlerHealth := handlers.NewRestHealth()
	handlerItemCreate := handlers.NewItemCreatePublic(serviceItemCreate, cfg.Log)
	handlerItemGet := handlers.NewItemGetPublic(serviceItemGet, cfg.Log)
	handlerItemList := handlers.NewItemListPublic(serviceItemList, cfg.Log)
	handlerItemUpdate := handlers.NewItemUpdatePublic(serviceItemUpdate, cfg.Log)
	handlerItemDelete := handlers.NewItemDeletePublic(serviceItemDelete, cfg.Log)

	// =================================================================================================================
	// ROUTER
	// =================================================================================================================

	router := httpf.NewRouter(cfg.Rest, httpf.RouterOptions{
		Otel:          cfg.Otel,
		Logger:        cfg.HttpLogger,
		DowntimeStart: cfg.App.DowntimeStart,
		Liveness:      "/ping",
		Methods:       []string{http.MethodHead, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	})

	// Probes reach /ping without a budget.
	router.Get("/ping", handlerPing.ServeHTTP)

	router.Group(func(api chi.Router) {
		api.Use(httpf.NewLimits(cfg.Rest).Standard())

		api.Get("/healthcheck", handlerHealth.ServeHTTP)
		api.Post("/items", handlerItemCreate.ServeHTTP)
		api.Get("/items", handlerItemList.ServeHTTP)
		api.Get("/item", handlerItemGet.ServeHTTP)
		api.Put("/item", handlerItemUpdate.ServeHTTP)
		api.Delete("/item", handlerItemDelete.ServeHTTP)
	})

	// =================================================================================================================
	// RUN
	// =================================================================================================================

	httpServer := httpf.NewServer(ctx, cfg.Rest, router)

	log.Println("Starting REST server on " + httpServer.Addr)

	err := httpf.Serve(ctx, httpServer, cfg.Rest.Timeouts.Shutdown)

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
