package client

import (
	"fmt"
	"github.com/redisTesting/internal/db"
	"github.com/redisTesting/internal/rstring"
	"math/rand"
	"sync"
	"time"
)

var KeyPoolSize = 1000

var ClientTimeout = 10 * time.Second
var NClientRequests = 1000
var ClientBatchSize = 1
var KeyLen = 8
var ValLen = 8

// a request is one or more kv-store operating commands (i.e, read or write command), depends on ClientBatchSize
type CommandLog struct {
	StartTime   time.Time     	// the start time of this command
	EndTime 	time.Time     	// the finish time of this command
	Duration    time.Duration 	// the latency of this command
}

/*
	A client
*/
type Client struct {
	ClientId 					uint32

	MasterClient 				*db.RedisClient  // SET and GET command uses this
	CommandDone 				chan int

	Rand    					*rand.Rand
	CommandLog 					[]CommandLog
	SentSoFar 					int
	startSending, endSending 	time.Time
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
		CommandDone: make(chan int),
		Rand: rand.New(rand.NewSource(time.Now().UnixNano() * int64(clientId))),

		CommandLog: make([]CommandLog, NClientRequests/ClientBatchSize),
	}
	/*
		SentSoFar is zero at initialization
		startSending, endSending retain their default values
	*/
	return c
}

// CloseLoopClient starts running a client
// It keeps issuing SET or GET commands
func (c *Client) CloseLoopClient(wg *sync.WaitGroup, keyPool *[]string)  {
	defer wg.Done()

	c.startSending = time.Now()

	ticker := time.NewTicker(ClientTimeout)

	MainLoop:
	// TODO: ClientBatchSize not meaningful yet
	for i := 0; i < NClientRequests/ClientBatchSize; i++ {

		go c.processOneCommand(i, keyPool)

		select {
		case <-c.CommandDone:
			continue
		case <-ticker.C:
			break MainLoop
		}
	}

	c.endSending = time.Now()
	fmt.Println("Time:", c.endSending.Sub(c.startSending))
}

// processOneCommand randomly runs SET or GET
// and keeps track of time
func (c *Client) processOneCommand(i int, keyPool *[]string)  {
	// Randomly get a key from keyPool
	// The key is used for SET or GET
	key := (*keyPool)[c.Rand.Intn(KeyPoolSize)]

	// Random 0 or 1
	command := c.Rand.Intn(2)

	// Run SET or GET and keep track of time
	tic := time.Now()
	if command == 0 {
		// Set
		value := rstring.RandString(c.Rand, ValLen)
		err := c.MasterClient.Set(key, value, 0, 1)
		if err != nil {
			panic(err)
		}
	} else {
		// Get
		_, err := c.MasterClient.Get(key)
		if err != nil {
			panic(err)
		}
	}
	toc := time.Now()

	// Save time into CommandLog of client
	c.CommandLog[i] = CommandLog{
		StartTime: tic,
		EndTime: toc,
		Duration: toc.Sub(tic),
	}

	c.CommandDone <- 1
	c.SentSoFar += ClientBatchSize
}

// StartNClients performs some initializations
// and then starts running n clients
func StartNClients(n int)  {
	var wg sync.WaitGroup

	// Initialize a pool of keys
	keyPool := InitKeyPool(KeyPoolSize)

	// Initialize n clients
	var allClients []*Client
	for idx := 0; idx < n; idx ++ {
		cli := ClientInit(uint32(idx))
		allClients = append(allClients, cli)
	}

	// Start all the clients
	for _, cli := range allClients {
		wg.Add(1)
		go cli.CloseLoopClient(&wg, keyPool)
	}

	wg.Wait()
}

// Initialize a key pool to pick from later
func InitKeyPool(keyPoolSize int) *[]string {
	// Temporary client to set some keys
	tempClient, err := db.NewClient("10.142.0.58:6379")
	if err != nil {
		panic("No connection")
	}
	tempClient.ChangePersistence()

	// Generate many keys and set random values
	// Then append keys to keyPool
	var keyPool []string
	for i := 0; i < keyPoolSize; i ++ {
		key := rstring.RandString(rand.New(rand.NewSource(time.Now().UnixNano())), KeyLen)
		value := rstring.RandString(rand.New(rand.NewSource(time.Now().UnixNano())), ValLen)
		err := tempClient.Set(key, value, 0, 1)
		if err != nil {
			panic(err)
		}
		keyPool = append(keyPool, key)
	}
	return &keyPool
}