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
func main ()  {
	cnt := 20000;
	
	doOne(cnt);
}

func doOne (cnt int) os.Error {

	var delta int64;
	spec := redis.DefaultSpec().Db(13);
	
	fmt.Printf("\n\n=== Bench synchclient ================ 1 Client -- %d opts --- \n", cnt );
	fmt.Println ();
	
	client, e := redis.NewSynchClientWithSpec (spec);
	if e != nil {  return onError ("on NewSynchClient call: ", e); }
	if client == nil {  return failedTest ("NewSynchClient returned nil!"); }

	client.Flushdb();
		
	delta = doPing(client, cnt);
	report ("PING", delta, cnt);
	
	delta = doIncr(client, cnt);
	report ("INCR", delta, cnt);
	
	delta = doSet(client, cnt);
	report ("SET", delta, cnt);
	
	delta = doGet(client, cnt);
	report ("GET", delta, cnt);
	
	delta = doSadd(client, cnt);
	report ("SADD", delta, cnt);
	
	delta = doLpush(client, cnt);
	report ("LPUSH", delta, cnt);
	
	delta = doRpush(client, cnt);
	report ("RPUSH", delta, cnt);
	
	delta = doLpop(client, cnt);
	report ("LPOP", delta, cnt);
	
	delta = doRpop(client, cnt);
	report ("RPOP", delta, cnt);
	
	client.Quit();
	return nil;
}

func report (cmd string, delta int64, cnt int) {
	fmt.Printf("---\n");
	fmt.Printf(fmt.Sprintf("cmd: %s\n", cmd));
	fmt.Printf(fmt.Sprintf("%d iterations of %s in %d msecs\n", cnt, cmd, delta/1000000));
}

func doPing (client redis.Client, cnt int) (delta int64) {
	t0 := time.Nanoseconds();
	for i:=0;i<cnt;i++ {	
		client.Ping();
	}
	delta = time.Nanoseconds() - t0;
	client.Flushdb();
	return;
}
func doIncr (client redis.Client, cnt int) (delta int64) {
	key := "ctr";
	t0 := time.Nanoseconds();
	for i:=0;i<cnt;i++ {	
		client.Incr(key);
	}
	delta = time.Nanoseconds() - t0;
	client.Flushdb();
	return;
}
func doSet (client redis.Client, cnt int) (delta int64) {
	key := "ctr";
	value := strings.Bytes("foo");
	t0 := time.Nanoseconds();
	for i:=0;i<cnt;i++ {	
		client.Set(key, value);
	}
	delta = time.Nanoseconds() - t0;
	client.Flushdb();
	return;
}
func doGet (client redis.Client, cnt int) (delta int64) {
	key := "ctr";
	t0 := time.Nanoseconds();
	for i:=0;i<cnt;i++ {	
		client.Get(key);
	}
	delta = time.Nanoseconds() - t0;
	client.Flushdb();
	return;
}
func doSadd (client redis.Client, cnt int) (delta int64) {
	key := "set";
	value:= strings.Bytes("one");
	t0 := time.Nanoseconds();
	for i:=0;i<cnt;i++ {	
		client.Sadd(key, value);
	}
	delta = time.Nanoseconds() - t0;
	client.Flushdb();
	return;
}
func doLpush (client redis.Client, cnt int) (delta int64) {
	key := "list-L";
	value:= strings.Bytes("foo");
	t0 := time.Nanoseconds();
	for i:=0;i<cnt;i++ {	
		client.Lpush(key, value);
	}
	delta = time.Nanoseconds() - t0;
	return;
}
func doRpush (client redis.Client, cnt int) (delta int64) {
	key := "list-R";
	value:= strings.Bytes("foo");
	t0 := time.Nanoseconds();
	for i:=0;i<cnt;i++ {	
		client.Lpush(key, value);
	}
	delta = time.Nanoseconds() - t0;
	return;
}
func doLpop (client redis.Client, cnt int) (delta int64) {
	key := "list-L";
	t0 := time.Nanoseconds();
	for i:=0;i<cnt;i++ {	
		client.Lpop(key);
	}
	delta = time.Nanoseconds() - t0;
	return;
}
func doRpop (client redis.Client, cnt int) (delta int64) {
	key := "list-R";
	t0 := time.Nanoseconds();
	for i:=0;i<cnt;i++ {	
		client.Lpop(key);
	}
	delta = time.Nanoseconds() - t0;
	return;
}
