//   Copyright 2009 Joubin Houshyar
// 
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//    
//   http://www.apache.org/licenses/LICENSE-2.0
//    
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//

/*
	Package redis defines the interfaces for Redis clients. 
*/
package redis

import (
	"flag";
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// Redis 'KeyType'
//
type KeyType byte;

// Known Redis key types
//
const (
	RT_NONE 	KeyType = iota;
	RT_STRING;
	RT_SET;
	RT_LIST;
	RT_ZSET;
)

// Returns KeyType by name
//
func GetKeyType (typename string) (keytype KeyType) {
	switch {
	case typename == "none": keytype =  RT_NONE;
	case typename == "string": keytype =  RT_STRING;
	case typename == "list": keytype =  RT_LIST;
	case typename == "set": keytype =  RT_SET;
	case typename == "zset": keytype =  RT_ZSET;
	}
	return;
}

// Not yet used -- TODO: decide if returning status (say for Set) for non error cases 
// really buys us anything beyond (useless) consistency.
//
type Status bool;
const (
	OK 	 bool	= true;	 	
	PONG 		= true;		
	ERR 		= false;
)

// ----------------------------------------------------------------------------
// Interfaces
// ----------------------------------------------------------------------------


// The synchronous call semantics Client interface.
//
// Method names map one to one to the Redis command set.
// All methods may return an redis.Error, which is either
// a system error (runtime issue or bug) or Redis error (i.e. user error)
// See Error in this package for details of its interface.

type Client interface {

	// Redis GET command.
	Get (key string) (result []byte, err Error);

	// Redis TYPE command.
	Type (key string) (result KeyType, err Error);

	// Redis SET command.
	Set (key string, arg1 []byte) (Error);

	// Redis SET command.
//	Set (key string, arg1 string) (Error);

	// Redis SET command.
//	Set (key string, arg1 int64) (Error);

	// Redis SET command.
//	Set (key string, arg1 io.ReadWriter) (Error);

	// Redis SAVE command.
	Save () (Error);

	// Redis KEYS command using "*" wildcard 
	AllKeys () (result []string, err Error);

	// Redis KEYS command.
	Keys (key string) (result []string, err Error);

	// Redis SORT command.
//	Sort (key string) (result redis.Sort, err Error);

	// Redis EXISTS command.
	Exists (key string) (result bool, err Error);

	// Redis RENAME command.
	Rename (key string, arg1 string) (Error);

	// Redis INFO command.
	Info () (result map[string]string, err Error);

	// Redis PING command.
	Ping () (Error);

	// Redis QUIT command.
	Quit () (Error);

	// Redis SETNX command.
//	Setnx (key string, arg1 io.ReadWriter) (result bool, err Error);

	// Redis SETNX command.
	Setnx (key string, arg1 []byte) (result bool, err Error);

	// Redis SETNX command.
//	Setnx (key string, arg1 string) (result bool, err Error);

	// Redis SETNX command.
//	Setnx (key string, arg1 int64) (result bool, err Error);

	// Redis GETSET command.
//	Getset (key string, arg1 string) (result []byte, err Error);

	// Redis GETSET command.
//	Getset (key string, arg1 int64) (result []byte, err Error);

	// Redis GETSET command.
//	Getset (key string, arg1 io.ReadWriter) (result []byte, err Error);

	// Redis GETSET command.
	Getset (key string, arg1 []byte) (result []byte, err Error);

	// Redis MGET command.
	Mget (key string, arg1 []string) (result [][]byte, err Error);

	// Redis INCR command.
	Incr (key string) (result int64, err Error);

	// Redis INCRBY command.
	Incrby (key string, arg1 int64) (result int64, err Error);

	// Redis DECR command.
	Decr (key string) (result int64, err Error);

	// Redis DECRBY command.
	Decrby (key string, arg1 int64) (result int64, err Error);

	// Redis DEL command.
	Del (key string) (result bool, err Error);

	// Redis RANDOMKEY command.
	Randomkey () (result string, err Error);

	// Redis RENAMENX command.
	Renamenx (key string, arg1 string) (result bool, err Error);

	// Redis DBSIZE command.
	Dbsize () (result int64, err Error);

	// Redis EXPIRE command.
	Expire (key string, arg1 int64) (result bool, err Error);

	// Redis TTL command.
	Ttl (key string) (result int64, err Error);

	// Redis RPUSH command.
//	Rpush (key string, arg1 io.ReadWriter) (Error);

	// Redis RPUSH command.
	Rpush (key string, arg1 []byte) (Error);

	// Redis RPUSH command.
//	Rpush (key string, arg1 string) (Error);

	// Redis RPUSH command.
//	Rpush (key string, arg1 int64) (Error);

	// Redis LPUSH command.
//	Lpush (key string, arg1 io.ReadWriter) (Error);

	// Redis LPUSH command.
//	Lpush (key string, arg1 int64) (Error);

	// Redis LPUSH command.
//	Lpush (key string, arg1 string) (Error);

	// Redis LPUSH command.
	Lpush (key string, arg1 []byte) (Error);

	// Redis LSET command.
//	Lset (key string, arg1 int64, arg2 io.ReadWriter) (Error);

	// Redis LSET command.
	Lset (key string, arg1 int64, arg2 []byte) (Error);

	// Redis LSET command.
//	Lset (key string, arg1 int64, arg2 string) (Error);

	// Redis LSET command.
//	Lset (key string, arg1 int64, arg2 int64) (Error);

	// Redis LREM command.
//	Lrem (key string, arg1 io.ReadWriter, arg2 int64) (result int64, err Error);

	// Redis LREM command.
	Lrem (key string, arg1 []byte, arg2 int64) (result int64, err Error);

	// Redis LREM command.
//	Lrem (key string, arg1 string, arg2 int64) (result int64, err Error);

	// Redis LREM command.
//	Lrem (key string, arg1 int64, arg2 int64) (result int64, err Error);

	// Redis LLEN command.
	Llen (key string) (result int64, err Error);

	// Redis LRANGE command.
	Lrange (key string, arg1 int64, arg2 int64) (result [][]byte, err Error);

	// Redis LTRIM command.
	Ltrim (key string, arg1 int64, arg2 int64) (Error);

	// Redis LINDEX command.
	Lindex (key string, arg1 int64) (result []byte, err Error);

	// Redis LPOP command.
	Lpop (key string) (result []byte, err Error);

	// Redis RPOP command.
	Rpop (key string) (result []byte, err Error);

	// Redis RPOPLPUSH command.
	Rpoplpush (key string, arg1 string) (result []byte, err Error);

	// Redis SADD command.
	Sadd (key string, arg1 []byte) (result bool, err Error);

	// Redis SADD command.
//	Sadd (key string, arg1 io.ReadWriter) (result bool, err Error);

	// Redis SADD command.
//	Sadd (key string, arg1 int64) (result bool, err Error);

	// Redis SADD command.
//	Sadd (key string, arg1 string) (result bool, err Error);

	// Redis SREM command.
	Srem (key string, arg1 []byte) (result bool, err Error);

	// Redis SREM command.
//	Srem (key string, arg1 string) (result bool, err Error);

	// Redis SREM command.
//	Srem (key string, arg1 int64) (result bool, err Error);

	// Redis SREM command.
//	Srem (key string, arg1 io.ReadWriter) (result bool, err Error);

	// Redis SISMEMBER command.
//	Sismember (key string, arg1 io.ReadWriter) (result bool, err Error);

	// Redis SISMEMBER command.
//	Sismember (key string, arg1 int64) (result bool, err Error);

	// Redis SISMEMBER command.
//	Sismember (key string, arg1 string) (result bool, err Error);

	// Redis SISMEMBER command.
	Sismember (key string, arg1 []byte) (result bool, err Error);

	// Redis SMOVE command.
//	Smove (key string, arg1 string, arg2 io.ReadWriter) (result bool, err Error);

	// Redis SMOVE command.
	Smove (key string, arg1 string, arg2 []byte) (result bool, err Error);

	// Redis SMOVE command.
//	Smove (key string, arg1 string, arg2 string) (result bool, err Error);

	// Redis SMOVE command.
//	Smove (key string, arg1 string, arg2 int64) (result bool, err Error);

	// Redis SCARD command.
	Scard (key string) (result int64, err Error);

	// Redis SINTER command.
	Sinter (key string, arg1 []string) (result [][]byte, err Error);

	// Redis SINTERSTORE command.
	Sinterstore (key string, arg1 []string) (Error);

	// Redis SUNION command.
	Sunion (key string, arg1 []string) (result [][]byte, err Error);

	// Redis SUNIONSTORE command.
	Sunionstore (key string, arg1 []string) (Error);

	// Redis SDIFF command.
	Sdiff (key string, arg1 []string) (result [][]byte, err Error);

	// Redis SDIFFSTORE command.
	Sdiffstore (key string, arg1 []string) (Error);

	// Redis SMEMBERS command.
	Smembers (key string) (result [][]byte, err Error);

	// Redis SRANDMEMBER command.
	Srandmember (key string) (result []byte, err Error);

	// Redis ZADD command.
	Zadd (key string, arg1 float64, arg2 []byte) (result bool, err Error);

	// Redis ZADD command.
//	Zadd (key string, arg1 float64, arg2 string) (result bool, err Error);

	// Redis ZADD command.
//	Zadd (key string, arg1 float64, arg2 int64) (result bool, err Error);

	// Redis ZADD command.
//	Zadd (key string, arg1 float64, arg2 io.ReadWriter) (result bool, err Error);

	// Redis ZREM command.
//	Zrem (key string, arg1 io.ReadWriter) (result bool, err Error);

	// Redis ZREM command.
//	Zrem (key string, arg1 int64) (result bool, err Error);

	// Redis ZREM command.
//	Zrem (key string, arg1 string) (result bool, err Error);

	// Redis ZREM command.
	Zrem (key string, arg1 []byte) (result bool, err Error);

	// Redis ZCARD command.
	Zcard (key string) (result int64, err Error);

	// Redis ZSCORE command.
//	Zscore (key string, arg1 io.ReadWriter) (result float64, err Error);

	// Redis ZSCORE command.
//	Zscore (key string, arg1 int64) (result float64, err Error);

	// Redis ZSCORE command.
	Zscore (key string, arg1 []byte) (result float64, err Error);

	// Redis ZSCORE command.
//	Zscore (key string, arg1 string) (result float64, err Error);

	// Redis ZRANGE command.
	Zrange (key string, arg1 int64, arg2 int64) (result [][]byte, err Error);

	// Redis ZREVRANGE command.
	Zrevrange (key string, arg1 int64, arg2 int64) (result [][]byte, err Error);

	// Redis ZRANGEBYSCORE command.
	Zrangebyscore (key string, arg1 float64, arg2 float64) (result [][]byte, err Error);

	// Redis FLUSHDB command.
	Flushdb () (Error);

	// Redis FLUSHALL command.
	Flushall () (Error);

	// Redis MOVE command.
	Move (key string, arg1 int64) (result bool, err Error);

	// Redis BGSAVE command.
	Bgsave () (Error);

	// Redis LASTSAVE command.
	Lastsave () (result int64, err Error);
}


// The asynchronous call semantics client interface.
//
//  WIP:  what the usage will look like
//
//	func usingRedisAsync () Error {
//		spec := DefaultConnectionSpec();
//		pipeline := NewRedisPipeline(spec);
//		
//		// futureBytes is a FutureBytes
//		futureBytes, reqErr := pipline.Get("my-key");
//		if reqErr != nil { return withError (reqErr); }
//		
//		// ....
//		
//		[]byte, execErr := futureBytes.Get();
//		if execErr != nil { return withError (execErr); }
//		
//		// or using timeouts
//		
//		timeout := 1000000; // 1 msec
//		[]byte, execErr, ok := futureBytes.TryGet (timeout);
//		if !ok { 
//			// we timedout 
//		}
//		else {
//			if execErr != nil { return withError (execErr); }
//		}
//	}

type AsyncClient interface {

// THESE ARE WIP/TEMP ... SUBJECT TO CHANGE.
//
//	Info () (FutureInfo, Error);
	Exists (key string) (FutureBool, Error);
//	Sismembers (key string) (FutureBytesArray, Error);
//	Incr (key string) (FutureInt64, Error);
//	Randomkey () (FutureString, Error);
	Get (key string) (FutureBytes, Error);
}

// ----------------------------------------------------------------------------
// package initiatization and internal ops and flags
// ----------------------------------------------------------------------------

// ----------------
// flags
//
// go-redis will make use of command line flags where available.  flag names
// for this package are all prefixed by "redis:" to prevent possible name collisions.
//
func init () { flag.Parse(); }

// redis:d
//
// global debug flag for redis package components.
// 
var _debug *bool = flag.Bool ("redis:d", false, "debug flag for go-redis"); // TEMP: should default to false
func debug() bool { return *_debug; }

