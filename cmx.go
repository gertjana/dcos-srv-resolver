package main

import (
	"log"
	"os"
)

func main() {
	if err := os.Chmod("dumb-init", 0755); err != nil {
		log.Fatal(err)
	}
}
