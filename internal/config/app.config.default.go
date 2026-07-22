package config

import (
	"os"
	"time"

	"github.com/samber/lo"

	"github.com/a-novel-kit/golib/logging"
	loggingpresets "github.com/a-novel-kit/golib/logging/presets"
	"github.com/a-novel-kit/golib/otel"
	otelpresets "github.com/a-novel-kit/golib/otel/presets"

	"github.com/a-novel/service-template/internal/config/env"
)

const (
	// OtelFlushTimeout bounds how long the OpenTelemetry exporter waits to flush
	// pending spans on shutdown before giving up.
	OtelFlushTimeout = 2 * time.Second
)

// LoggerProd ships production logs to Google Cloud Logging.
var LoggerProd = loggingpresets.GRPCGcloud{
	Component: env.GcloudProjectId,
}

// LoggerDev pretty-prints logs to the console for local development.
var LoggerDev = loggingpresets.GRPCLocal{}

// LoggerDevHttp pretty-prints HTTP-level logs to the console for local development.
var LoggerDevHttp = &loggingpresets.LogLocal{
	Out: os.Stdout,
}

// LoggerProdHttp ships production HTTP-level logs to Google Cloud Logging.
var LoggerProdHttp = &loggingpresets.LogGcloud{
	ProjectId: env.GcloudProjectId,
}

// AppPresetDefault is the configuration the service starts with. It reads every
// value from the environment, and picks the Google Cloud logging and tracing
// backends once a project ID is set.
var AppPresetDefault = App{
	App: Main{
		Name: env.AppName,
	},
	Grpc: Grpc{
		Port: env.GrpcPort,
		Ping: env.GrpcPing,
	},
	Rest: Rest{
		Port: env.RestPort,
		Timeouts: RestTimeouts{
			Read:       env.RestTimeoutRead,
			ReadHeader: env.RestTimeoutReadHeader,
			Write:      env.RestTimeoutWrite,
			Idle:       env.RestTimeoutIdle,
			Request:    env.RestTimeoutRequest,
		},
		MaxRequestSize: env.RestMaxRequestSize,
		Cors: RestCors{
			AllowedOrigins:   env.CorsAllowedOrigins,
			AllowedHeaders:   env.CorsAllowedHeaders,
			AllowCredentials: env.CorsAllowCredentials,
			MaxAge:           env.CorsMaxAge,
		},
	},

	Otel: lo.If[otel.Config](!env.Otel, &otelpresets.Disabled{}).
		ElseIf(env.GcloudProjectId == "", &otelpresets.Local{
			FlushTimeout: OtelFlushTimeout,
		}).
		Else(&otelpresets.Gcloud{
			ProjectID:    env.GcloudProjectId,
			FlushTimeout: OtelFlushTimeout,
		}),
	Log:    lo.Ternary[logging.Log](env.GcloudProjectId == "", LoggerDevHttp, LoggerProdHttp),
	Logger: lo.Ternary[logging.RPCConfig](env.GcloudProjectId == "", &LoggerDev, &LoggerProd),
	HttpLogger: lo.Ternary[logging.HTTPConfig](
		env.GcloudProjectId == "",
		&loggingpresets.HTTPLocal{BaseLogger: LoggerDevHttp},
		&loggingpresets.HTTPGcloud{BaseLogger: LoggerProdHttp},
	),
	Postgres: PostgresPresetDefault,
}
