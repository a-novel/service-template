package handlers

import (
	"context"
	"net/http"

	"github.com/samber/lo"
	"github.com/uptrace/bun"

	"github.com/a-novel-kit/golib/httpf"
	"github.com/a-novel-kit/golib/otel"
	"github.com/a-novel-kit/golib/postgres"
)

const (
	// RestHealthStatusUp marks a dependency the service can currently reach.
	RestHealthStatusUp = "up"
	// RestHealthStatusDown marks a dependency the service failed to reach.
	RestHealthStatusDown = "down"
)

// RestHealthStatus reports the reachability of a single backing dependency in
// the health response. Err carries the failure message when the status is down.
type RestHealthStatus struct {
	Status string `json:"status"`
	Err    string `json:"err,omitempty"`
}

// NewRestHealthStatus builds a status from a dependency probe: up when err is
// nil, down with err's message otherwise.
func NewRestHealthStatus(err error) *RestHealthStatus {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	return &RestHealthStatus{
		Status: lo.Ternary(err == nil, RestHealthStatusUp, RestHealthStatusDown),
		Err:    errMsg,
	}
}

// RestHealth is the HTTP handler for the health endpoint. It probes each
// backing dependency and reports their individual status.
type RestHealth struct{}

func NewRestHealth() *RestHealth {
	return &RestHealth{}
}

func (handler *RestHealth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer().Start(r.Context(), "rest.Health")
	defer span.End()

	httpf.SendJSON(ctx, w, span, map[string]any{
		"client:postgres": NewRestHealthStatus(handler.reportPostgres(ctx)),
	})
}

func (handler *RestHealth) reportPostgres(ctx context.Context) error {
	ctx, span := otel.Tracer().Start(ctx, "rest.Health(reportPostgres)")
	defer span.End()

	pg, err := postgres.GetContext(ctx)
	if err != nil {
		return otel.ReportError(span, err)
	}

	pgdb, ok := pg.(*bun.DB)
	if !ok {
		// In transaction mode the pooled handle is a transaction rather than a
		// *bun.DB and exposes no Ping; treat the dependency as healthy.
		return nil
	}

	err = pgdb.Ping()
	if err != nil {
		return otel.ReportError(span, err)
	}

	otel.ReportSuccessNoContent(span)

	return nil
}
