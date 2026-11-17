// Package handlers is the transport layer of the service, sitting above the core layer. Handlers
// translate incoming requests into core calls, and map the results and sentinel errors back to
// protocol responses and status codes.
//
// Two transports serve distinct audiences:
//
//   - gRPC handlers: the private API, for other services on the internal network. Wired by
//     cmd/grpc, which admits only the callers whose service token it accepts.
//
//   - REST handlers: the API for the same callers, behind the same service-token gate. A fork that
//     serves any HTTP client drops that gate in cmd/rest. Wired by cmd/rest.
package handlers
