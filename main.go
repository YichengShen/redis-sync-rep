package main

import (
	"fmt"
	"github.com/redisTesting/client"
)

func main()  {
	client, err := client.NewClient("localhost:6379")
	if err != nil {
		panic("No connection")
	}

	err = client.Set("key", "This key has been set", 0, 0)
	if err != nil {
		panic(err)
	}
	fmt.Println(client.Get("key"))

}
