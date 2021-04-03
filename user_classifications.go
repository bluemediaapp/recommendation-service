package main

import (
	"github.com/gofiber/fiber/v2"
)

func userClassifications()  {
	app.Post("/user/:user_id", func(ctx *fiber.Ctx) error {
		return nil
	})
}
