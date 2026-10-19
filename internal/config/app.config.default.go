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

// OtelFlushTimeout bounds how long the OpenTelemetry exporter waits to flush pending spans on
// shutdown.
const OtelFlushTimeout = 2 * time.Second

// LoggerProdGrpc sends gRPC request logs to Google Cloud Logging.
var LoggerProdGrpc = loggingpresets.GRPCGcloud{
	Component: env.AppName,
	ProjectId: env.GcloudProjectId,
}

// LoggerDevGrpc prints gRPC request logs to the console for local development.
var LoggerDevGrpc = loggingpresets.GRPCLocal{}

// LoggerDev prints application and REST request logs to the console for local development.
var LoggerDev = &loggingpresets.LogLocal{
	Out: os.Stdout,
}

// LoggerProd sends application and REST request logs to Google Cloud Logging.
var LoggerProd = &loggingpresets.LogGcloud{
	ProjectId: env.GcloudProjectId,
}

// AppPresetDefault is the configuration the service starts with. It reads every value from the
// environment, and picks the Google Cloud logging and tracing backends once a project ID is set.
var AppPresetDefault = App{
	App: Main{
		Name:          env.AppName,
		DowntimeStart: env.DowntimeStart,
	},
	Grpc: env.Grpc,
	Rest: env.Rest,

	Otel: lo.If[otel.Config](!env.Otel, &otelpresets.Disabled{}).
		ElseIf(env.GcloudProjectId == "", &otelpresets.Local{
			FlushTimeout: OtelFlushTimeout,
		}).
		Else(&otelpresets.Gcloud{
			ProjectID:    env.GcloudProjectId,
			FlushTimeout: OtelFlushTimeout,
		}),
	Logger:     lo.Ternary[logging.Log](env.GcloudProjectId == "", LoggerDev, LoggerProd),
	GrpcLogger: lo.Ternary[logging.RPCConfig](env.GcloudProjectId == "", &LoggerDevGrpc, &LoggerProdGrpc),
	RestLogger: lo.Ternary[logging.HTTPConfig](
		env.GcloudProjectId == "",
		&loggingpresets.HTTPLocal{BaseLogger: LoggerDev},
		&loggingpresets.HTTPGcloud{BaseLogger: LoggerProd},
	),
	Postgres: PostgresPresetDefault,
}
