package main

import (
	"hetzner-dns/internal/app"
	"os"
)

func main() {
	os.Exit(app.Run())
}
