package client

import (
	"fmt"
	cfg "github.com/redisTesting/internal/config"
	"github.com/redisTesting/internal/db"
	"github.com/redisTesting/internal/rstring"
	"github.com/rs/zerolog"
	"math/rand"
	"net"
	"os"
	"sort"
	"sync"
	"time"
)

var KeyPoolSize = 1000 // number of keys generated for selection later

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
	// Connect to master
	addr := net.JoinHostPort(cfg.Conf.MasterIp, cfg.Conf.MasterPort)
	masterClient, err := db.NewClient(addr)
	if err != nil {
		panic("No connection")
	}

	c := &Client{
		ClientId: clientId,
		MasterClient: masterClient,
		CommandDone: make(chan int),
		Rand: rand.New(rand.NewSource(time.Now().UnixNano() * int64(clientId))),

		CommandLog: make([]CommandLog, cfg.Conf.NClientRequests/cfg.Conf.ClientBatchSize),
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

	ClientTimeout := time.Duration(cfg.Conf.ClientTimeout) * time.Second
	ticker := time.NewTicker(ClientTimeout)

	MainLoop:
	for i := 0; i < cfg.Conf.NClientRequests/cfg.Conf.ClientBatchSize; i++ {

		if cfg.Conf.ClientBatchSize > 1 {
			go c.processBatchCommands(i, keyPool)
		} else {
			go c.processOneCommand(i, keyPool)
		}

		select {
		case <-c.CommandDone:
			continue
		case <-ticker.C:
			break MainLoop
		}
	}

	c.endSending = time.Now()
	c.writeToLog()
}

// processOneCommand randomly runs SET or GET for one time
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
		value := rstring.RandString(c.Rand, cfg.Conf.ValLen)
		err := c.MasterClient.Set(key, value, 0)
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
	c.SentSoFar += cfg.Conf.ClientBatchSize
}

// processBatchCommands runs MGET or MSET
func (c *Client) processBatchCommands(idx int, keyPool *[]string) {
	setCmds := make([]string, cfg.Conf.ClientBatchSize*2) // pending MSET requests
	getCmds := make([]string, cfg.Conf.ClientBatchSize)   // pending MGET requests
	setCount, getCount := 0, 0                   // elements counter

	for i := 0; i < cfg.Conf.ClientBatchSize; i ++ {
		// Random 0 or 1
		command := c.Rand.Intn(2)

		key := (*keyPool)[c.Rand.Intn(KeyPoolSize)]

		if command == 0 {
			// Append set command
			value := rstring.RandString(c.Rand, cfg.Conf.ValLen)
			setCmds[setCount] = key
			setCount += 1
			setCmds[setCount] = value
			setCount += 1
		} else {
			// Append get command
			getCmds[getCount] = key
			getCount += 1
		}
	}

	// Run MGET and MSET and track time
	tic := time.Now()
	if setCount > 1 {
		c.MasterClient.Mset(setCmds[:setCount])
	}
	if getCount > 1 {
		c.MasterClient.Mget(getCmds[:getCount])
	}
	toc := time.Now()

	// Save time into CommandLog of client
	c.CommandLog[idx] = CommandLog{
		StartTime: tic,
		EndTime: toc,
		Duration: toc.Sub(tic),
	}

	c.CommandDone <- 1
	c.SentSoFar += cfg.Conf.ClientBatchSize
}

// Write log to file
func (c *Client) writeToLog() {
	RepliedLength := len(c.CommandLog) // assume all replied
	//for i := 0; i < len(c.CommandLog); i++ {
	//	if c.CommandLog[i].Duration == time.Duration(0) {
	//		//c.Logger.Warn("", Int("not replied", i))
	//		RepliedLength = i // update RepliedLength if necessary
	//		break
	//	}
	//}
	//fmt.Println("RepliedLength =", RepliedLength)

	// cmdLogs -- exclude head and tails statistics in BatchedCmdLog:
	cmdLogs := make([]CommandLog, int(float64(RepliedLength)*0.8))
	j := 0
	for i := 0; i < len(c.CommandLog); i++ {
		if i < int(float64(RepliedLength)*0.1) ||
			i >= int(float64(RepliedLength)*0.9) {
			continue
		}
		if j < len(cmdLogs) {
			cmdLogs[j] = c.CommandLog[i]
			j++
		} else {
			break
		}
	}
	//fmt.Println("RepliedLength =", RepliedLength,
	//	"len of cmdLogs = ", len(cmdLogs), ",", j, "items filled")

	maxLatVal := time.Duration(0)
	maxLatIdx := 0
	for i, cmd := range cmdLogs {
		if cmd.Duration > maxLatVal {
			maxLatVal = cmd.Duration
			maxLatIdx = i + int(float64(RepliedLength)*0.1)
		}
	}

	mid80Start := cmdLogs[0].StartTime
	mid80End := cmdLogs[len(cmdLogs)-1].EndTime
	mid80Dur := mid80End.Sub(mid80Start).Seconds()

	mid80FirstRecvTime := cmdLogs[0].EndTime
	mid80RecvTimeDur := mid80End.Sub(mid80FirstRecvTime).Seconds()

	sort.Slice(cmdLogs, func(i, j int) bool {
		return cmdLogs[i].Duration < cmdLogs[j].Duration
	})
	minLat := cmdLogs[0].Duration
	maxLat := cmdLogs[len(cmdLogs)-1].Duration
	p50Lat := cmdLogs[int(float64(len(cmdLogs))*0.5)].Duration
	p95Lat := cmdLogs[int(float64(len(cmdLogs))*0.9)].Duration
	p99Lat := cmdLogs[int(float64(len(cmdLogs))*0.99)].Duration
	var durSum int64
	for _, v := range cmdLogs {
		durSum += v.Duration.Microseconds()
	}
	durAvg := durSum / int64(len(cmdLogs))

	// Make a log folder if necessary
	if err := os.MkdirAll(cfg.Conf.LogDir, os.ModePerm); err != nil {
		panic(err)
	}
	// Create a file
	file, err := os.Create(fmt.Sprintf("logs/client%d.json", c.ClientId))
	if err != nil {
		panic(err)
	}

	logger := zerolog.New(zerolog.MultiLevelWriter(file))

	logger.Warn().
		Uint32("ClientId", c.ClientId).
		Int("TotalSent", c.SentSoFar).
		Int64("minLat", minLat.Microseconds()).
		Int64("maxLat", maxLat.Microseconds()).
		Int("maxLatIdx", maxLatIdx).
		Int64("avgLat", durAvg).
		Int64("p50Lat", p50Lat.Microseconds()).
		Int64("p95Lat", p95Lat.Microseconds()).
		Int64("p99Lat", p99Lat.Microseconds()).
		Int64("sendStart", c.startSending.UnixNano()).
		Int64("sendEnd", c.endSending.UnixNano()).
		Int64("mid80Start", mid80Start.UnixNano()).
		Int64("mid80End", mid80End.UnixNano()).
		Float64("mid80Dur", mid80Dur).
		Float64("mid80RecvTimeDur", mid80RecvTimeDur).
		Int("mid80Requests", len(cmdLogs)).
		Float64("mid80Throughput (cmd/sec)", float64(len(cmdLogs))/mid80Dur).
		Float64("mid80Throughput2 (cmd/sec)", float64(len(cmdLogs))/mid80RecvTimeDur).Msg("")

}

// StartNClients performs some initializations
// and then starts running n clients
func StartNClients(n int)  {
	var wg sync.WaitGroup

	// Initialize settings for Redis master and replicas
	db.SpecialDbInit()

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
	addr := net.JoinHostPort(cfg.Conf.MasterIp, cfg.Conf.MasterPort)
	tempClient, err := db.NewClient(addr)
	if err != nil {
		panic("No connection")
	}

	// Generate many keys and set random values
	// Then append keys to keyPool
	var keyPool []string
	for i := 0; i < keyPoolSize; i ++ {
		key := rstring.RandString(rand.New(rand.NewSource(time.Now().UnixNano())), cfg.Conf.KeyLen)
		value := rstring.RandString(rand.New(rand.NewSource(time.Now().UnixNano())), cfg.Conf.ValLen)
		err := tempClient.Set(key, value, 0)
		if err != nil {
			panic(err)
		}
		keyPool = append(keyPool, key)
	}
	return &keyPool
}