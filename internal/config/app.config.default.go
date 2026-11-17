package config

import (
	"github.com/a-novel/service-template/internal/config/env"
)

// AppPresetDefault is the configuration the service starts with. It reads every value from the
// environment. Once telemetry exports over OTLP, spans record each request and call, so the servers
// write no request log.
var AppPresetDefault = App{
	App: Main{
		Name:          env.AppName,
		DowntimeStart: env.DowntimeStart,
	},
	Grpc:         env.Grpc,
	Rest:         env.Rest,
	ServiceToken: env.ServiceToken,

	Otel:       env.Telemetry.Otel,
	Logger:     env.Telemetry.Logger,
	GrpcLogger: env.Telemetry.RPCLogger,
	RestLogger: env.Telemetry.HTTPLogger,
	Postgres:   PostgresPresetDefault,
}
