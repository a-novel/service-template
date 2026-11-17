// Package handlers is the transport layer of the service, sitting above the core layer. Handlers
// translate incoming requests into core calls, and map the results and sentinel errors back to
// protocol responses and status codes.
//
// Two transports serve distinct audiences:
//
//   - gRPC handlers: the private API, for other services on the internal network. [GrpcServiceToken]
//     admits only the callers whose service token it accepts. Wired by cmd/grpc.
//
//   - REST handlers: the public API. [RestServiceToken] restricts it to the same callers until a
//     fork opens it to any HTTP client. Wired by cmd/rest.
//
// Both gates fail closed: a server without a service-token issuer refuses to start, unless
// SERVICE_ISSUER_OPTIONAL opts a local or test stack out.
package handlers
