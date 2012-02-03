package main

import (
	"os"
	"flag"
	"log"
	"fmt"
	"time"
	"redis"
)

// ----------------------------------------------------------------------------
// types and props
// ----------------------------------------------------------------------------

// redis task function type def
type redisTask func(id string, signal chan int, client redis.AsyncClient, iterations int)

// task info
type taskSpec struct {
	task redisTask
	name string
}

// array of Tasks to run in sequence
// Add a task to the list to bench to the runner.
// Tasks are run in sequence.
var tasks = []taskSpec{
	taskSpec{doPing, "PING"},
	taskSpec{doSet, "SET"},
	taskSpec{doGet, "GET"},
	taskSpec{doIncr, "INCR"},
	taskSpec{doDecr, "DECR"},
	taskSpec{doLpush, "LPUSH"},
	taskSpec{doLpop, "LPOP"},
	taskSpec{doRpush, "RPUSH"},
	taskSpec{doRpop, "RPOP"},
}

// workers option.  default is equiv to -w=10 on command line
var workers = flag.Int("w", 10, "number of concurrent workers")

// opcnt option.  default is equiv to -n=2000 on command line
var opcnt = flag.Int("n", 10000, "number of task iterations per worker")

// ----------------------------------------------------------------------------
// benchmarker
// ----------------------------------------------------------------------------

func main() {
	// DEBUG
	log.SetPrefix("[go-redis|bench] ")
	flag.Parse()

	fmt.Printf("\n\n=== Bench synchclient ================ %d Concurrent Clients -- %d opts each --- \n\n", *workers, *opcnt)
	for _, task := range tasks {
		benchTask(task, *opcnt, *workers, true)
	}

}

// Use a single redis.AsyncClient with specified number
// of workers to bench concurrent load on the async client
func benchTask(taskspec taskSpec, iterations int, workers int, printReport bool) (delta int64, err os.Error) {
	signal := make(chan int, workers) // Buffering optional but sensible.
	spec := redis.DefaultSpec().Db(13).Password("go-redis")
	client, e := redis.NewAsynchClientWithSpec(spec)
	if e != nil {
		log.Println("Error creating client for worker: ", e)
		return -1, e
	}
	//    defer client.Quit()        // will be deprecated soon
	defer client.RedisClient().Quit()

	t0 := time.Nanoseconds()
	for i := 0; i < workers; i++ {
		id := fmt.Sprintf("%d", i)
		go taskspec.task(id, signal, client, iterations)
	}
	for i := 0; i < workers; i++ {
		<-signal
	}
	delta = time.Nanoseconds() - t0
	//	for i := 0; i < workers; i++ {
	//		clients[i].Quit()
	//	}
	//
	if printReport {
		report("concurrent "+taskspec.name, delta, iterations*workers)
	}

	return
}
func report(cmd string, delta int64, cnt int) {
	fmt.Printf("---\n")
	fmt.Printf("cmd: %s\n", cmd)
	fmt.Printf("%d iterations of %s in %d msecs\n", cnt, cmd, delta/1000000)
	fmt.Printf("---\n\n")
}

// ----------------------------------------------------------------------------
// redis tasks
// ----------------------------------------------------------------------------

func doPing(id string, signal chan int, client redis.AsyncClient, cnt int) {
	for i := 0; i < cnt; i++ {
		client.Ping()
	}
	signal <- 1
}
func doIncr(id string, signal chan int, client redis.AsyncClient, cnt int) {
	key := "ctr-" + id
	for i := 0; i < cnt; i++ {
		client.Incr(key)
	}
	signal <- 1
}
func doDecr(id string, signal chan int, client redis.AsyncClient, cnt int) {
	key := "ctr-" + id
	for i := 0; i < cnt; i++ {
		client.Decr(key)
	}
	signal <- 1
}
func doSet(id string, signal chan int, client redis.AsyncClient, cnt int) {
	key := "set-" + id
	value := []byte("foo")
	for i := 0; i < cnt; i++ {
		client.Set(key, value)
	}
	signal <- 1
}
func doGet(id string, signal chan int, client redis.AsyncClient, cnt int) {
	key := "set-" + id
	for i := 0; i < cnt; i++ {
		client.Get(key)
	}
	signal <- 1
}
func doLpush(id string, signal chan int, client redis.AsyncClient, cnt int) {
	key := "list-L-" + id
	value := []byte("foo")
	for i := 0; i < cnt; i++ {
		client.Lpush(key, value)
	}
	signal <- 1
}
func doLpop(id string, signal chan int, client redis.AsyncClient, cnt int) {
	key := "list-L-" + id
	for i := 0; i < cnt; i++ {
		client.Lpop(key)
	}
	signal <- 1
}

func doRpush(id string, signal chan int, client redis.AsyncClient, cnt int) {
	key := "list-R-" + id
	value := []byte("foo")
	for i := 0; i < cnt; i++ {
		client.Rpush(key, value)
	}
	signal <- 1
}
func doRpop(id string, signal chan int, client redis.AsyncClient, cnt int) {
	key := "list-R" + id
	for i := 0; i < cnt; i++ {
		client.Rpop(key)
	}
	signal <- 1
}
