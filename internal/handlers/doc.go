// Package handlers is the transport layer of the service, sitting above the core layer. Handlers
// translate incoming requests into core calls, and map the results and sentinel errors back to
// protocol responses and status codes.
//
// Two transports serve distinct audiences:
//
//   - gRPC handlers: the private API, for other services on the internal network. Once configured,
//     [GrpcServiceToken] admits only the callers whose service token it accepts. Otherwise the
//     network around it restricts access. Wired by cmd/grpc.
//
//   - REST handlers: the public API, for any HTTP client, unless [RestServiceToken] restricts it to
//     the same callers. Wired by cmd/rest.
package handlers
