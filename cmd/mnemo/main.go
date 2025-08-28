package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/LinaKACI-pro/mnemo/internal/config"
	"github.com/LinaKACI-pro/mnemo/store/sqlite"
)

var version = "0.0.1"

func main() {
	showVersion := flag.Bool("v", false, "show version")
	flag.Parse()

	if *showVersion {
		fmt.Println("mnemo version", version)
		os.Exit(0)
	}

	// Load Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %s", err)
	}

	// Initialize SQLite store
	db, err := sqlite.New(cfg.DbDriver, cfg.DbPath)
	if err != nil {
		log.Fatalf("failed to init sqlite store: %v", err)
	}
	defer db.Close()

	fmt.Println("Hello from mnemo!")
}
