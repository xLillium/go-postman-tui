package main

import (
    "log"

    "github.com/xlillium/go-postman-tui/internal/app"
)

func main() {
    if err := app.Run(); err != nil {
        log.Fatalf("Application failed: %v", err)
    }
}

