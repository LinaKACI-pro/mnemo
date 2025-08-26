package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/LinaKACI-pro/mnemo/internal/config"
	"github.com/LinaKACI-pro/mnemo/store/sqlite"
)

var version = "0.0.1"

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %s", err)
	}

	db, err := sqlite.New(cfg.DbDriver, cfg.DbPath)
	if err != nil {
		log.Fatalf("failed to init sqlite store: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Search
	results, err := db.Search(ctx, "hello")
	if err != nil {
		log.Fatalf("search failed: %v", err)
	}

	for _, e := range results {
		fmt.Printf("Found entry: %d - %s\n", e.ID, e.Value)
	}

	showVersion := flag.Bool("v", false, "show version")
	flag.Parse()

	if *showVersion {
		fmt.Println("mnemo version", version)
		os.Exit(0)
	}

	fmt.Println("Hello from mnemo!")
}
