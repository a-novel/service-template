package configtest

import "os"

// GrpcURL is the address of the running gRPC server the Go client tests dial, in the form
// <host>:<port>.
var GrpcURL = os.Getenv("GRPC_URL")
