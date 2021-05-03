package client

import (
	"fmt"
	"github.com/redisTesting/internal/db"
	"github.com/redisTesting/internal/rstring"
	"math/rand"
	"sync"
	"time"
)

var ClientTimeout = 10 * time.Second
var NClientRequests = 1000
var ClientBatchSize = 1
var KeyLen = 8
var ValLen = 8

/*
	A client
*/
type Client struct {
	ClientId uint32
	MasterClient *db.RedisClient  // SET and GET command uses this
	RedisGetChan chan string

	Rand    *rand.Rand
	startSending, endSending, endReceiving time.Time
}

/*
	Initialize a client
*/
func ClientInit(clientId uint32) *Client {
	// TODO: ip is hard coded
	// Connect to master
	masterClient, err := db.NewClient("10.142.0.58:6379")
	if err != nil {
		panic("No connection")
	}
	masterClient.ChangePersistence()

	c := &Client{
		ClientId: clientId,
		MasterClient: masterClient,
		RedisGetChan: make(chan string),
		Rand: rand.New(rand.NewSource(time.Now().UnixNano() * int64(clientId))),
	}
	/*
		SentSoFar, ReceivedSoFar are zeros are initialization
		startSending, endSending, endReceiving retain their default values
	*/
	return c
}

func (c *Client) CloseLoopClient(wg *sync.WaitGroup)  {
	defer wg.Done()

	c.startSending = time.Now()

	ticker := time.NewTicker(ClientTimeout)
MainLoop:
	// TODO: ClientBatchSize not meaningful yet
	for i := 0; i < NClientRequests/ClientBatchSize; i++ {

		key := rstring.RandString(c.Rand, KeyLen)
		value := rstring.RandString(c.Rand, ValLen)

		// Set
		err := c.MasterClient.Set(key, value, 0, 1)
		if err != nil {
			panic(err)
		}

		// Get
		go func() {
			val, _ := c.MasterClient.Get(key)
			c.RedisGetChan <- val
		}()

		select {
		case <-c.RedisGetChan:
			continue
		case <-ticker.C:
			break MainLoop
		}
	}

	c.endSending = time.Now()
	c.endReceiving = time.Now()
	fmt.Println("Time:", c.endSending.Sub(c.startSending))
}

func StartNClients(n int)  {
	var wg sync.WaitGroup

	// Initialize n clients
	var allClients []*Client
	for idx := 0; idx < n; idx ++ {
		cli := ClientInit(uint32(idx))
		allClients = append(allClients, cli)
	}

	// Start all the clients
	for _, cli := range allClients {
		wg.Add(1)
		go cli.CloseLoopClient(&wg)
	}

	wg.Wait()
}