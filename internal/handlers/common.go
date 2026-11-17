package handlers

import "errors"

// ErrServiceTokenIssuerRequired refuses a server configured with no service-token issuer: it would
// take every call.
var ErrServiceTokenIssuerRequired = errors.New("service token issuer required")
