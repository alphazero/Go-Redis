package main

import (
	"os";
	"redis";
	"strings";
	"log";
	"fmt";
	"time";
)
func onError (msg string, e os.Error) os.Error {
	log.Stderr (msg, "", e);
	return e;
}
func failedTest (msg string) os.Error {
	log.Stderr (msg);
	return nil;
}
func report (cmd string, delta int64, cnt int) {
	fmt.Printf("---\n");
	fmt.Printf(fmt.Sprintf("cmd: %s\n", cmd));
	fmt.Printf(fmt.Sprintf("%d iterations of %s in %d msecs\n", cnt, cmd, delta/1000000));
}

func main ()  {
	cnt := 20000;
	workers := 50;
	
	doConcurrent (cnt/workers, workers);
}
func doConcurrent (cnt int, workers int) os.Error {

	fmt.Printf("\n\n=== Bench synchclient ================ %d Concurrent Clients -- %d opts each --- \n", workers, cnt);
	fmt.Println ();
	var delta int64;
	
	delta= doConcurrentPing (cnt, workers);
	report ("concurrent PING", delta, cnt*workers);
	
	delta= doConcurrentIncr (cnt, workers);
	report ("concurrent INCR", delta, cnt*workers);
	
	delta= doConcurrentSet (cnt, workers);
	report ("concurrent SET", delta, cnt*workers);
	
	delta= doConcurrentGet (cnt, workers);
	report ("concurrent GET", delta, cnt*workers);
	
	delta= doConcurrentLpush (cnt, workers);
	report ("concurrent LPUSH", delta, cnt*workers);
	
	delta= doConcurrentRpush (cnt, workers);
	report ("concurrent RPUSH", delta, cnt*workers);
	
	return nil;	
}
func makeConcurrentClients(workers int) (clients []redis.Client) {
    clients = make([]redis.Client, workers);
    for i := 0; i < workers; i++ {
		spec := redis.DefaultSpec().Db(13);
		client, _ := redis.NewSynchClientWithSpec (spec);
		clients[i] = client;
    }
    return;
}
func doConcurrentPing (cnt int, workers int) (delta int64)  {

    signal := make(chan int);  // Buffering optional but sensible.
    clients := makeConcurrentClients(workers);
	t0 := time.Nanoseconds();
    for i := 0; i < workers; i++ {
        go doPing2 (signal, clients[i], cnt);
    }
    for i := 0; i < workers; i++ { <-signal }
    for i := 0; i < workers; i++ { clients[i].Quit(); }
    return time.Nanoseconds() - t0;
}
func doConcurrentIncr (cnt int, workers int) (delta int64)  {

    signal := make(chan int);  // Buffering optional but sensible.
    clients := makeConcurrentClients(workers);
	t0 := time.Nanoseconds();
    for i := 0; i < workers; i++ {
        go doIncr2 (signal, clients[i], cnt);
    }
    for i := 0; i < workers; i++ { <-signal }
    for i := 0; i < workers; i++ { clients[i].Quit(); }
    return time.Nanoseconds() - t0;
}
func doConcurrentSet (cnt int, workers int) (delta int64)  {

    signal := make(chan int);  // Buffering optional but sensible.
    clients := makeConcurrentClients(workers);
	t0 := time.Nanoseconds();
    for i := 0; i < workers; i++ {
        go doSet2 (signal, clients[i], cnt);
    }
    for i := 0; i < workers; i++ { <-signal }
    for i := 0; i < workers; i++ { clients[i].Quit(); }
    return time.Nanoseconds() - t0;
}
func doConcurrentGet (cnt int, workers int) (delta int64)  {

    signal := make(chan int);  // Buffering optional but sensible.
    clients := makeConcurrentClients(workers);
	t0 := time.Nanoseconds();
    for i := 0; i < workers; i++ {
        go doGet2 (signal, clients[i], cnt);
    }
    for i := 0; i < workers; i++ { <-signal }
    for i := 0; i < workers; i++ { clients[i].Quit(); }
    return time.Nanoseconds() - t0;
}
func doConcurrentRpush (cnt int, workers int) (delta int64)  {

    signal := make(chan int);  // Buffering optional but sensible.
    clients := makeConcurrentClients(workers);
	t0 := time.Nanoseconds();
    for i := 0; i < workers; i++ {
        go doRpush2 (signal, clients[i], cnt);
    }
    for i := 0; i < workers; i++ { <-signal }
    for i := 0; i < workers; i++ { clients[i].Quit(); }
    return time.Nanoseconds() - t0;
}
func doConcurrentLpush (cnt int, workers int) (delta int64)  {

    signal := make(chan int);  // Buffering optional but sensible.
    clients := makeConcurrentClients(workers);
	t0 := time.Nanoseconds();
    for i := 0; i < workers; i++ {
        go doLpush2 (signal, clients[i], cnt);
    }
    for i := 0; i < workers; i++ { <-signal }
    for i := 0; i < workers; i++ { clients[i].Quit(); }
    return time.Nanoseconds() - t0;
}

func doPing2 (signal chan int, client redis.Client, cnt int)  {
	for i:=0;i<cnt;i++ { 
		client.Ping();
	}
	signal <- 1;
}
func doIncr2 (signal chan int, client redis.Client, cnt int)  {
	key := "ctr";
	for i:=0;i<cnt;i++ { 
		client.Incr(key);
	}
	signal <- 1;
}
func doSet2 (signal chan int, client redis.Client, cnt int)  {
	key := "ctr";
	value := strings.Bytes("foo");
	for i:=0;i<cnt;i++ { 
		client.Set(key, value);
	}
	signal <- 1;
}
func doGet2 (signal chan int, client redis.Client, cnt int)  {
	key := "ctr";
	for i:=0;i<cnt;i++ { 
		client.Get(key);
	}
	signal <- 1;
}
func doLpush2 (signal chan int, client redis.Client, cnt int)  {
	key := "list-L";
	value := strings.Bytes("foo");
	for i:=0;i<cnt;i++ { 
		client.Lpush(key, value);
	}
	signal <- 1;
}
func doRpush2 (signal chan int, client redis.Client, cnt int)  {
	key := "list-R";
	value := strings.Bytes("foo");
	for i:=0;i<cnt;i++ { 
		client.Rpush(key, value);
	}
	signal <- 1;
}
