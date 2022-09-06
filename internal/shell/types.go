package shell

import goa "goa.design/goa/v3/pkg"

type EndpointMiddleware func(goa.Endpoint) goa.Endpoint // EndpointMiddleware are executed around the endpoint handler.

type key string

// RequestIdKey is the key used to store the request ID in the context or header.
const RequestIdKey key = "X-Request-Id"
