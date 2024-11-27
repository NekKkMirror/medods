package main

import (
	"fmt"
	"log"
	"os"

	"github.com/NekKkMirror/medods-tz.git/internal/app"
)

func main() {
	r, err := app.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize the application: %v", err)
	}

	port := os.Getenv("PORT")
	address := fmt.Sprintf("%s:%s", "0.0.0.0", port)

	if err := r.Run(address); err != nil {
		log.Fatalf("Failed to run the server: %v", err)
	}
}
