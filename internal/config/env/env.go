// Package env reads and parses the service's configuration from environment
// variables, exposing each setting as a typed, ready-to-use value.
package env

import (
	"os"

	"github.com/a-novel-kit/golib/config"
	"github.com/a-novel-kit/golib/downtime"
	"github.com/a-novel-kit/golib/grpcf"
	"github.com/a-novel-kit/golib/httpf"
)

// prefix is prepended to every configuration variable name that this package reads.
// Set SERVICE_TEMPLATE_ENV_PREFIX when embedding the service in another project,
// where unprefixed variable names could collide with the host project's own.
var prefix = os.Getenv("SERVICE_TEMPLATE_ENV_PREFIX")

func getEnv(name string) string {
	return os.Getenv(prefix + name)
}

// Default values applied when an environment variable is unset.
const (
	AppNameDefault = "service-template"

	// PostgresMaxOpenConnsDefault keeps the pool well under a stock PostgreSQL
	// max_connections of 100 once multiplied by a service's replica count, leaving
	// room for the migration job and a psql session. Go's own default is unlimited,
	// which turns a spike into connection refusals for everything on that database
	// rather than queueing inside this process.
	PostgresMaxOpenConnsDefault = 20
	// PostgresMaxIdleConnsDefault matches the open limit so a burst does not close
	// connections it is about to reopen.
	PostgresMaxIdleConnsDefault = 20
	PostgresPortDefault         = 5432
	PostgresTLSEnabledDefault   = true
)

// Raw values for environment variables.
var (
	postgresDsn          = getEnv("POSTGRES_DSN")
	postgresHost         = getEnv("POSTGRES_HOST")
	postgresPort         = getEnv("POSTGRES_PORT")
	postgresUser         = getEnv("POSTGRES_USER")
	postgresPassword     = getEnv("POSTGRES_PASSWORD")
	postgresDatabase     = getEnv("POSTGRES_DATABASE")
	postgresTLSEnabled   = getEnv("POSTGRES_TLS_ENABLED")
	postgresMaxOpenConns = getEnv("POSTGRES_MAX_OPEN_CONNS")
	postgresMaxIdleConns = getEnv("POSTGRES_MAX_IDLE_CONNS")

	appName = getEnv("APP_NAME")
	otel    = getEnv("OTEL")

	grpcUrl = getEnv("GRPC_URL")

	gcloudProjectId = getEnv("GCLOUD_PROJECT_ID")

	downtimeStart = getEnv("DOWNTIME_START")
)

var (
	// PostgresDsn is the legacy URL used when PostgresHost is empty, formatted as
	//	postgres://<user>:<password>@<host>:<port>/<database>
	//
	// Configure the discrete PostgreSQL fields for new deployments.
	PostgresDsn = postgresDsn
	// PostgresHost is the PostgreSQL server hostname or IP address.
	PostgresHost = postgresHost
	// PostgresPort is the PostgreSQL server port.
	PostgresPort = config.LoadEnv(postgresPort, PostgresPortDefault, config.IntParser)
	// PostgresUser is the PostgreSQL login role.
	PostgresUser = postgresUser
	// PostgresPassword is the PostgreSQL login credential.
	PostgresPassword = postgresPassword
	// PostgresDatabase is the PostgreSQL database name.
	PostgresDatabase = postgresDatabase
	// PostgresTLSEnabled controls transport encryption for the PostgreSQL connection.
	PostgresTLSEnabled = config.LoadEnv(postgresTLSEnabled, PostgresTLSEnabledDefault, config.BoolParser)

	// PostgresMaxOpenConns is the maximum number of open connections to the database.
	PostgresMaxOpenConns = config.LoadEnv(postgresMaxOpenConns, PostgresMaxOpenConnsDefault, config.IntParser)
	// PostgresMaxIdleConns is the maximum number of connections kept open while idle.
	PostgresMaxIdleConns = config.LoadEnv(postgresMaxIdleConns, PostgresMaxIdleConnsDefault, config.IntParser)

	// AppName is the name of the application, as it appears in logs and tracing.
	AppName = config.LoadEnv(appName, AppNameDefault, config.StringParser)
	// Otel enables OpenTelemetry instrumentation.
	//
	// See: https://opentelemetry.io/
	Otel = config.LoadEnv(otel, false, config.BoolParser)

	// GrpcUrl is the URL of the gRPC service, typically <host>:<port>.
	GrpcUrl = grpcUrl
	// Grpc configures the gRPC server from the GRPC_* variables.
	Grpc = grpcf.LoadServerConfig(getEnv, grpcf.ServerConfigDefault)
	// Rest configures the REST server from the REST_* variables.
	Rest = httpf.LoadServerConfig(getEnv, httpf.ServerConfigDefault)

	// GcloudProjectId names the Google Cloud project the service runs in. Setting
	// it switches logging and tracing from the local console to Google Cloud.
	//
	// See: https://docs.cloud.google.com/resource-manager/docs/creating-managing-projects
	GcloudProjectId = gcloudProjectId

	// DowntimeStart is when a planned downtime of this service starts, in RFC 3339. From then
	// until it is removed, the service refuses work. Nil when none is planned.
	DowntimeStart = config.LoadEnv(downtimeStart, nil, downtime.ParseStart)
)
