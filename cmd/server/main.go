package main

import (
	"log"

	"github.com/DangVTNhan/goacl/internal/app"
)

func main() {
	application := app.New()
	if err := application.Run(); err != nil {
		log.Fatalf("Failed to run application: %v", err)
	}
}
