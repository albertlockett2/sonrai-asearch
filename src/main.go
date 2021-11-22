package main

import (
	"fmt"
	"github.com/sonraisecurity/sonrai-asearch/src/collector"
	"github.com/sonraisecurity/sonrai-asearch/src/manager"
	"github.com/sonraisecurity/sonrai-asearch/src/results"
	"github.com/sonraisecurity/sonrai-asearch/src/worker"
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatalf("usage: ./asearch <manager/worker/collector/results>")
	}

	switch args[1] {
	case "collector":
		collector.Main()
	case "manager":
		manager.Main()
	case "results":
		results.Main()
	case "worker":
		worker.Main()
	default:
		log.Fatalf(fmt.Sprintf("unknown arg %s", args[0]))
	}
}
