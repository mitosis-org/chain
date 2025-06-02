package main

import (
	"log"
	"os"

	"github.com/mitosis-org/chain/cmd/mito/internal/commands"
)

func main() {
	rootCmd := commands.NewRootCmd()

	if err := rootCmd.Execute(); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
}
