package main

import "ecommerce-go/internal/bootstrap"

// Thin entrypoint: delegates wiring + server start to the bootstrap package.
func main() {
	bootstrap.Run()
}
