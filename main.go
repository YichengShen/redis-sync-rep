package main

import (
	"github.com/redisTesting/deployment/analysis"
	cfg "github.com/redisTesting/internal/config"
	"github.com/redisTesting/roles/client"
	"os"
)

func main()  {
	// Remove logs
	err := os.RemoveAll(cfg.Conf.LogDir)
	if err != nil {
		panic(err)
	}

	// Run n clients
	client.StartNClients(cfg.Conf.NClients)

	// Analysis
	analysis.RunAnalysis(cfg.Conf.LogDir)
}

