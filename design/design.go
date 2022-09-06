package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = API("Rotabot", func() {
	Title("Rotabot Service")
	Description("SlackApp that makes team rotations easy")
	Server("web", func() {
		Host("localhost", func() {
			URI("http://localhost:8080")
		})
	})
})

var _ = Service("Rotabot", func() {
	Description("Service responsible for handling commands to create, update, and manage rotas")

	Method("Healthcheck", func() {

		HTTP(func() {
			Response(StatusOK, func() {
				ContentType("application/json")
			})
			GET("/healthcheck")
		})

		Result(Empty)

	})

})
