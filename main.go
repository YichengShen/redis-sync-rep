package main

import (
	"github.com/redisTesting/roles/client"
)

func main()  {
	client.StartNClients(100)
}

