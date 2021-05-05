package main

import (
	"github.com/redisTesting/deployment/analysis"
	"github.com/redisTesting/roles/client"
	"os"
)

func main()  {
	// Remove logs
	err := os.RemoveAll("logs")
	if err != nil {
		panic(err)
	}

	// Run n clients
	client.StartNClients(100)

	// Analysis
	analysis.RunAnalysis("logs")
}

