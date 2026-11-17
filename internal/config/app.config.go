package config

import (
	"time"

	"github.com/a-novel-kit/golib/grpcf"
	"github.com/a-novel-kit/golib/httpf"
	"github.com/a-novel-kit/golib/logging"
	"github.com/a-novel-kit/golib/otel"
	"github.com/a-novel-kit/golib/postgres"
	"github.com/a-novel-kit/golib/servicetoken"
)

// Main holds the top-level application settings.
type Main struct {
	// Name of the application, as it appears in logs and tracing.
	Name string
	// DowntimeStart is when a planned downtime starts; nil when none is planned. From then until
	// it is removed, the service refuses work and leaves its database alone.
	DowntimeStart *time.Time
}

// App is the complete configuration the service reads at startup.
type App struct {
	// App holds the top-level application settings.
	App Main
	// Grpc configures the gRPC server.
	Grpc grpcf.ServerConfig
	// Rest configures the REST server.
	Rest httpf.ServerConfig
	// ServiceToken names the service tokens both servers accept. Without an issuer, the servers refuse
	// to start unless ServiceIssuerOptional is set.
	ServiceToken servicetoken.Config
	// ServiceIssuerOptional lets the servers start without a service-token issuer, and then take every
	// call. Only a local or test stack sets it.
	ServiceIssuerOptional bool

	// Otel configures the OpenTelemetry exporter.
	Otel otel.Config
	// Logger is the application logger the REST handlers report their errors to.
	Logger logging.Log
	// GrpcLogger configures the gRPC request logging middleware.
	GrpcLogger logging.RPCConfig
	// RestLogger configures the REST request logging middleware.
	RestLogger logging.HTTPConfig
	// Postgres configures the PostgreSQL connection.
	Postgres postgres.Config
}
