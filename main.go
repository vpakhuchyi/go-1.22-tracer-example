package main

import (
	"os"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/exp/trace"
)

func main() {
	fr := trace.NewFlightRecorder()
	if err := fr.Start(); err != nil {
		// handle error
	}

	defer func() {
		if err := fr.Stop(); err != nil {
			// handle error
		}
	}()

	app := fiber.New()

	app.Get("/trace", func(ctx fiber.Ctx) error {
		f, err := os.OpenFile("file.out", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
		if err != nil {
			// handle error
		}

		if _, err := fr.WriteTo(f); err != nil {
			// handle error
		}
		return ctx.SendStatus(fiber.StatusOK)
	})

	go runGame()

	if err := app.Listen(":8080"); err != nil {
		// handle error
		app.Shutdown()
	}
}
