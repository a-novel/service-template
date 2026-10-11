// Package config assembles the runtime configuration for the service: the typed
// structs the application reads, and the default preset that populates them from
// the environment.
package config

import (
	"time"

	"github.com/a-novel-kit/golib/grpcf"
	"github.com/a-novel-kit/golib/httpf"
	"github.com/a-novel-kit/golib/logging"
	"github.com/a-novel-kit/golib/otel"
	"github.com/a-novel-kit/golib/postgres"
)

// Main holds the top-level application settings.
type Main struct {
	// Name of the application, as it appears in logs and tracing.
	Name string `json:"name" yaml:"name"`
	// DowntimeStart is when a planned downtime starts; nil when none is planned. From then until
	// it is removed, the service refuses work and leaves its database alone.
	DowntimeStart *time.Time `json:"downtimeStart" yaml:"downtimeStart"`
}

// App is the complete configuration consumed by the service at startup, grouping
// the server, observability, logging, and database settings.
type App struct {
	App  Main               `json:"app"  yaml:"app"`
	Grpc grpcf.ServerConfig `json:"grpc" yaml:"grpc"`
	Rest httpf.ServerConfig `json:"rest" yaml:"rest"`

	Otel       otel.Config        `json:"otel"       yaml:"otel"`
	Log        logging.Log        `json:"log"        yaml:"log"`
	Logger     logging.RPCConfig  `json:"logger"     yaml:"logger"`
	HttpLogger logging.HTTPConfig `json:"httpLogger" yaml:"httpLogger"`
	Postgres   postgres.Config    `json:"postgres"   yaml:"postgres"`
}
