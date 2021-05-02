package main

import (
	"fmt"
	"github.com/redisTesting/client"
)

func main()  {
	client, err := client.NewClient("10.142.0.58:6379")
	if err != nil {
		panic("No connection")
	}

	client.ChangePersistence()

	err = client.Set("key", "This key has been set", 0, 1)
	if err != nil {
		panic(err)
	}
	fmt.Println(client.Get("key"))

	client.PrintInfo()
}
