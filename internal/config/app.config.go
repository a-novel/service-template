// Package config assembles the runtime configuration for the service: the typed
// structs the application reads, and the default preset that populates them from
// the environment.
package config

import (
	"time"

	"github.com/a-novel-kit/golib/logging"
	"github.com/a-novel-kit/golib/otel"
	"github.com/a-novel-kit/golib/postgres"
)

// RestCors holds CORS configuration for the REST server.
type RestCors struct {
	AllowedOrigins   []string `json:"allowedOrigins"   yaml:"allowedOrigins"`
	AllowedHeaders   []string `json:"allowedHeaders"   yaml:"allowedHeaders"`
	AllowCredentials bool     `json:"allowCredentials" yaml:"allowCredentials"`
	MaxAge           int      `json:"maxAge"           yaml:"maxAge"`
}

// Main holds the top-level application settings.
type Main struct {
	// Name of the application, as it appears in logs and tracing.
	Name string `json:"name" yaml:"name"`
}

// Grpc holds the gRPC server configuration.
type Grpc struct {
	// Port on which the gRPC server listens for incoming requests.
	Port int `json:"port" yaml:"port"`
	// Ping is the refresh interval for the gRPC server's internal health check.
	Ping time.Duration `json:"ping" yaml:"ping"`
}

// RestTimeouts holds timeout configuration for the REST server.
type RestTimeouts struct {
	Read       time.Duration `json:"read"       yaml:"read"`
	ReadHeader time.Duration `json:"readHeader" yaml:"readHeader"`
	Write      time.Duration `json:"write"      yaml:"write"`
	Idle       time.Duration `json:"idle"       yaml:"idle"`
	Request    time.Duration `json:"request"    yaml:"request"`
}

// Rest holds the REST server configuration.
type Rest struct {
	// Port on which the REST server listens for incoming requests.
	Port int `json:"port" yaml:"port"`
	// Timeouts bounds the lifecycle of a REST request.
	Timeouts RestTimeouts `json:"timeouts" yaml:"timeouts"`
	// MaxRequestSize is the maximum size of an incoming request body, in bytes.
	MaxRequestSize int64 `json:"maxRequestSize" yaml:"maxRequestSize"`
	// Cors holds the CORS configuration.
	Cors RestCors `json:"cors" yaml:"cors"`
}

// App is the complete configuration consumed by the service at startup, grouping
// the server, observability, logging, and database settings.
type App struct {
	App  Main `json:"app"  yaml:"app"`
	Grpc Grpc `json:"grpc" yaml:"grpc"`
	Rest Rest `json:"rest" yaml:"rest"`

	Otel       otel.Config        `json:"otel"       yaml:"otel"`
	Log        logging.Log        `json:"log"        yaml:"log"`
	Logger     logging.RPCConfig  `json:"logger"     yaml:"logger"`
	HttpLogger logging.HTTPConfig `json:"httpLogger" yaml:"httpLogger"`
	Postgres   postgres.Config    `json:"postgres"   yaml:"postgres"`
}
