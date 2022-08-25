// Code generated by goa v3.8.3, DO NOT EDIT.
//
// Rotabot service
//
// Command:
// $ goa gen github.com/kevinrobayna/rotabot/design

package rotabot

import (
	"context"
)

// Service responsible for handling commands to create, update, and manage rotas
type Service interface {
	// Healthcheck implements Healthcheck.
	Healthcheck(context.Context) (res string, err error)
	// Home implements Home.
	Home(context.Context) (res string, err error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "Rotabot"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [2]string{"Healthcheck", "Home"}
