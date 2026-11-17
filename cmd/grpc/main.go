// Command grpc serves the private gRPC API, for other services on the internal network. It admits
// only the callers its service tokens allow, and refuses to start without a token issuer unless
// SERVICE_ISSUER_OPTIONAL opts a local stack out. cmd/rest serves the same domain over HTTP, and
// cmd/migrations applies the database schema.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/a-novel-kit/golib/downtime"
	"github.com/a-novel-kit/golib/grpcf"
	"github.com/a-novel-kit/golib/httpf"
	"github.com/a-novel-kit/golib/otel"
	"github.com/a-novel-kit/golib/postgres"

	"github.com/a-novel/service-template/internal/config"
	"github.com/a-novel/service-template/internal/core"
	"github.com/a-novel/service-template/internal/dao"
	"github.com/a-novel/service-template/internal/handlers"
	templatev0 "github.com/a-novel/service-template/internal/handlers/protogen/anovel/template/v0"
)

func main() {
	cfg := config.AppPresetDefault

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	otel.SetAppName(cfg.App.Name)

	lo.Must0(otel.Init(cfg.Otel))
	defer cfg.Otel.Flush()

	// A server started during a planned downtime refuses work, so it leaves the database alone.
	if !downtime.Started(cfg.App.DowntimeStart, time.Now()) {
		ctx = lo.Must(postgres.NewContext(ctx, cfg.Postgres))
		defer cfg.Postgres.Close()
	}

	server := lo.Must(newServer(ctx, cfg))

	log.Println("Starting gRPC server on :" + strconv.Itoa(cfg.Grpc.Port))

	lo.Must0(grpcf.ListenAndServe(ctx, cfg.Grpc, server))
}

// newServer builds the gRPC server. Each call carries the database of ctx, and passes the
// service-token gate first, but for the health service. The gate refuses an empty issuer unless cfg
// opts a local stack out, so the server never starts open by mistake.
func newServer(ctx context.Context, cfg config.App) (*grpc.Server, error) {
	serviceToken, err := grpcf.ServiceTokenGate(cfg.ServiceToken, httpf.NewPoolClient(httpf.PoolOptions{}))
	if err != nil {
		return nil, fmt.Errorf("service token gate: %w", err)
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

	handlerStatus := handlers.NewGrpcStatus()
	handlerItemCreate := handlers.NewGrpcItemCreate(serviceItemCreate)
	handlerItemGet := handlers.NewGrpcItemGet(serviceItemGet)
	handlerItemList := handlers.NewGrpcItemList(serviceItemList)
	handlerItemUpdate := handlers.NewGrpcItemUpdate(serviceItemUpdate)
	handlerItemDelete := handlers.NewGrpcItemDelete(serviceItemDelete)

	// =================================================================================================================
	// SERVER
	// =================================================================================================================

	server := grpcf.NewGRPCServer(ctx, cfg.GrpcLogger, cfg.App.DowntimeStart, nil, serviceToken...)

	grpcf.RegisterHealth(ctx, server)

	templatev0.RegisterStatusServiceServer(server, handlerStatus)
	templatev0.RegisterItemCreateServiceServer(server, handlerItemCreate)
	templatev0.RegisterItemGetServiceServer(server, handlerItemGet)
	templatev0.RegisterItemListServiceServer(server, handlerItemList)
	templatev0.RegisterItemUpdateServiceServer(server, handlerItemUpdate)
	templatev0.RegisterItemDeleteServiceServer(server, handlerItemDelete)

	reflection.Register(server)

	return server, nil
}
