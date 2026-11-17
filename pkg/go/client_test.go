package servicetemplate_test

import (
	"context"
	"net"
	"os/exec"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	"github.com/a-novel-kit/golib/grpcf"
	golibproto "github.com/a-novel-kit/golib/grpcf/proto/gen"

	"github.com/a-novel/service-template/internal/config/configtest"
	"github.com/a-novel/service-template/pkg/go"
)

func TestClient(t *testing.T) {
	t.Parallel()

	t.Run("Success/RemoteCalls", func(t *testing.T) {
		t.Parallel()

		client, err := servicetemplate.NewClient(configtest.GrpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)

		defer client.Close()

		ctx := t.Context()

		_, err = client.UnaryEcho(ctx, &golibproto.UnaryEchoRequest{})
		require.NoError(t, err)

		_, err = client.Status(ctx, &servicetemplate.StatusRequest{})
		require.NoError(t, err)

		// One item goes through every RPC, which proves each one is registered on the running server.
		created, err := client.ItemCreate(ctx, &servicetemplate.ItemCreateRequest{Name: "client item"})
		require.NoError(t, err)

		id := created.GetItem().GetId()

		read, err := client.ItemGet(ctx, &servicetemplate.ItemGetRequest{Id: id})
		require.NoError(t, err)
		require.Equal(t, "client item", read.GetItem().GetName())

		updated, err := client.ItemUpdate(ctx, &servicetemplate.ItemUpdateRequest{Id: id, Name: "renamed item"})
		require.NoError(t, err)
		require.Equal(t, "renamed item", updated.GetItem().GetName())

		page, err := client.ItemList(ctx, &servicetemplate.ItemListRequest{Limit: 100})
		require.NoError(t, err)
		require.True(t, slices.ContainsFunc(page.GetItems(), func(item *servicetemplate.Item) bool {
			return item.GetId() == id
		}))

		deleted, err := client.ItemDelete(ctx, &servicetemplate.ItemDeleteRequest{Id: id})
		require.NoError(t, err)
		require.Equal(t, id, deleted.GetItem().GetId())

		_, err = client.ItemGet(ctx, &servicetemplate.ItemGetRequest{Id: id})
		require.Equal(t, codes.NotFound, status.Code(err))
	})

	t.Run("Success/DependencyBoundary", func(t *testing.T) {
		t.Parallel()

		output, err := exec.CommandContext(t.Context(), "go", "list", "-deps", ".").Output()
		require.NoError(t, err)

		// The client links the generated protobuf code of this module and nothing else of it: neither
		// the configuration, which reads the server's environment, nor the core or DAO layers.
		const internalPath = "github.com/a-novel/service-template/internal/"

		for dependency := range strings.SplitSeq(strings.TrimSpace(string(output)), "\n") {
			if strings.HasPrefix(dependency, internalPath) {
				require.True(t, strings.HasPrefix(dependency, internalPath+"handlers/protogen/"), dependency)
			}
		}
	})

	t.Run("Success/PropagatesTraceContext", func(t *testing.T) {
		t.Parallel()

		// golib's otel.Init installs this propagator in every service; the default one sends nothing.
		otel.SetTextMapPropagator(propagation.TraceContext{})

		received := make(chan []string, 1)
		server := grpc.NewServer(grpc.UnaryInterceptor(func(
			ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
		) (any, error) {
			md, _ := metadata.FromIncomingContext(ctx)
			received <- md.Get("traceparent")

			return handler(ctx, req)
		}))
		grpcf.RegisterHealth(t.Context(), server)

		listener := bufconn.Listen(1 << 16)
		go func() { _ = server.Serve(listener) }()

		defer server.Stop()

		client, err := servicetemplate.NewClient(
			"passthrough:///bufconn",
			grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
				return listener.DialContext(ctx)
			}),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		require.NoError(t, err)

		defer client.Close()

		callerConfig := trace.SpanContextConfig{
			TraceID:    trace.TraceID{1},
			SpanID:     trace.SpanID{1},
			TraceFlags: trace.FlagsSampled,
		}
		caller := trace.NewSpanContext(callerConfig)

		_, err = client.UnaryEcho(
			trace.ContextWithSpanContext(t.Context(), caller), &golibproto.UnaryEchoRequest{},
		)
		require.NoError(t, err)

		traceparent := <-received
		require.Len(t, traceparent, 1)
		require.Contains(t, traceparent[0], caller.TraceID().String())
	})
}
