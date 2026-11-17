// Package servicetemplate provides a gRPC client for the template service.
//
// [NewClient] dials the service and returns a [Client] with one method per RPC. The service admits
// only the callers its service tokens allow: dial it with the options of golib's
// grpcf.ServiceTokenCredentialsProvider.
package servicetemplate
