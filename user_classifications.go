package main

import (
	"github.com/gofiber/fiber"
)

func userClassifications()  {
	app.Post("/user/:user_id", func(ctx *fiber.Ctx) error {
		return nil
	})
}
