package main

import (
	"strings";
	"bytes";
	"log";
	"redis";
	"connection";
	"client";
)


func check (c redis.Client) {
	eping := c.Ping();
	if eping != nil { log.Stderr("error on Ping: ", eping); }
	
	es := c.Set("foobar", strings.Bytes("as usual"));
	if es != nil { log.Stderr("error on Set: ", es); }

	val, e := c.Get("foobar");
	if e != nil { log.Stderr("error on Get: ", e); }
	else { log.Stdout("(Get) foobar => ", bytes.NewBuffer(val).String()); }
	
	members, e2 := c.Smembers ("my-set");	
	if e2 != nil { log.Stderr("error on Smembers: ", e2); }
	else if members != nil {
		if len(members) > 0 {
			for i, member := range members {
				log.Stdout("[",i,"] => ", member);
			}
		}
		else {
			log.Stdout ("my-set IS EMPTY");
		}	
	}
	else {
		log.Stdout ("my-set smembers RETURNED NIL");
	}	
	
	members, e2 = c.Smembers ("set");	
	if e2 != nil { log.Stderr("error on Smembers: ", e2); }
	else if members != nil {
		if len(members) > 0 {
			for i, member := range members {
				log.Stdout("set:[",i,"] => ", member);
			}
		}
		else {
			log.Stdout ("set IS EMPTY");
		}	
	}
	else {
		log.Stdout ("set smembers RETURNED NIL");
	}	
	
	members, e2 = c.Smembers ("s");	
	if e2 != nil { log.Stderr("error on Smembers: ", e2); }
	else if members != nil {
		if len(members) > 0 {
			for i, member := range members {
				log.Stdout("s:[",i,"] => ", member);
			}
		}
		else {
			log.Stdout ("s IS EMPTY");
		}	
	}
	else {
		log.Stdout ("s smembers RETURNED NIL");
	}	
}

func main () {
/*
	list := make([][]byte, 10);
	for i, item := range list {
		log.Stdout ("i: ", i);
		log.Stdout ("item: ", item);
		list[i] = make([]byte, 22);
	} 
	for i, item := range list {
		log.Stdout ("i: ", i);
		log.Stdout ("item: ", len(item));
	} 
*/	
	c, err := client.Connect ();
	if(err != nil) {
		log.Stderr("Error on client.Connect: ", err);
	}
	else {
	/** BENCHMARK 
		val := strings.Bytes("333");
		for i :=0; i<100000; i++ {
			c.Set ("foo", val);
		}
		if true {
			return; 
		}
	***/
		c.Ping();
		c.Set("A-STRING", strings.Bytes("as usual"));
		c.Incr("A-COUNTER");
		t, _ := c.Type("A-COUNTER");
		log.Stderr("type of key is ", t);
		
		b1, _ := c.Exists("no-such");
		log.Stderr("Exists? ", b1);
		
		b2, _ := c.Exists("A-COUNTER");
		log.Stderr("Exists? ", b2);
//		var cntr int64;
		for i :=0; i<2; i++ {
			cntr, e := c.Incr("cntr");
			if e != nil {
				log.Stderr(e);
			} 
			else {
				log.Stdout(i, " >> ", cntr);
			}
		}
		check(c);
	
	}

	spec := connection.DefaultSpec();
	spec.Db(10).Password("foo").Port(6379);
	
	c2, err := client.ConnectWithSpec (spec);
	if(err != nil) {
		log.Stderr("Error on client.Connect: ", err);
		return;
	}
	c2.Ping();
	check(c2);
}
