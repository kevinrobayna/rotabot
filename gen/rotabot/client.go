// Code generated by goa v3.8.3, DO NOT EDIT.
//
// Rotabot client
//
// Command:
// $ goa gen github.com/kevinrobayna/rotabot/design

package rotabot

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "Rotabot" service client.
type Client struct {
	HealthcheckEndpoint goa.Endpoint
}

// NewClient initializes a "Rotabot" service client given the endpoints.
func NewClient(healthcheck goa.Endpoint) *Client {
	return &Client{
		HealthcheckEndpoint: healthcheck,
	}
}

// Healthcheck calls the "Healthcheck" endpoint of the "Rotabot" service.
func (c *Client) Healthcheck(ctx context.Context) (err error) {
	_, err = c.HealthcheckEndpoint(ctx, nil)
	return
}
