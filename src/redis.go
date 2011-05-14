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
	"flag"
//	"runtime";
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
/* 
	The synchronous client interface provides blockng call semantics supported by
	a distinct request/reply sequence at the connector level.
	
	Method names map one to one to the Redis command set.

	All methods may return an redis.Error, which is either a Redis error (from
	the server), or a system error indicating a runtime issue (or bug).
	See Error in this package for details of its interface.

	Usage example:

	func usingRedisSync () Error {
		spec := DefaultConnectionSpec();
		pipeline := NewAsynchClient(spec);

		// futureBytes is a FutureBytes
		value, reqErr := pipline.Get("my-key");
		if reqErr != nil { return withError (reqErr); }
	}

*/

type Client interface {

	// Redis GET command.
	Get (key string) (result []byte, err Error);

	// Redis TYPE command.
	Type (key string) (result KeyType, err Error);

	// Redis SET command.
	Set (key string, arg1 []byte) (Error);

	// Redis SAVE command.
	Save () (Error);

	// Redis KEYS command using "*" wildcard 
	AllKeys () (result []string, err Error);

	// Redis KEYS command.
	Keys (key string) (result []string, err Error);

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
	Setnx (key string, arg1 []byte) (result bool, err Error);

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
	Rpush (key string, arg1 []byte) (Error);

	// Redis LPUSH command.
	Lpush (key string, arg1 []byte) (Error);

	// Redis LSET command.
	Lset (key string, arg1 int64, arg2 []byte) (Error);

	// Redis LREM command.
	Lrem (key string, arg1 []byte, arg2 int64) (result int64, err Error);

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

	// Redis SREM command.
	Srem (key string, arg1 []byte) (result bool, err Error);

	// Redis SISMEMBER command.
	Sismember (key string, arg1 []byte) (result bool, err Error);

	// Redis SMOVE command.
	Smove (key string, arg1 string, arg2 []byte) (result bool, err Error);

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

	// Redis ZREM command.
	Zrem (key string, arg1 []byte) (result bool, err Error);

	// Redis ZCARD command.
	Zcard (key string) (result int64, err Error);

	// Redis ZSCORE command.
	Zscore (key string, arg1 []byte) (result float64, err Error);

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


/* 
	The asynchronous client interface provides asynchronous call semantics with
	future results supporting both blocking and try-and-timeout result accessors.
	
	Each method provides a type-safe future result return value, in addition to
	any (system) errors encountered in queuing the request.  
	
	The returned value may be ignored by clients that are not interested in the
	future response (for example on SET("foo", data)).  ALternatively, the caller
	may retain the future result referenced and perform blocking and/or timed wait
	gets on the expected response.
	
	[Try]Gets on the future result will return any Redis errors that were sent by
	the server, or, Go-Redis (system) errors encountered in processing the response.

	Usage example:

	func usingRedisAsync () Error {
		spec := DefaultConnectionSpec();
		pipeline := NewRedisPipeline(spec);

		// futureBytes is a FutureBytes
		futureBytes, reqErr := pipline.Get("my-key");
		if reqErr != nil { return withError (reqErr); }

		// ....

		[]byte, execErr := futureBytes.Get();
		if execErr != nil { return withError (execErr); }

		// or using timeouts

		timeout := 1000000; // 1 msec
		[]byte, execErr, ok := futureBytes.TryGet (timeout);
		if !ok { 
			// we timedout 
		}
		else {
			if execErr != nil { return withError (execErr); }
		}
	}

*/

type AsyncClient interface {

	// Redis GET command.
	Get (key string) (result FutureBytes, err Error);

	// Redis TYPE command.
	Type (key string) (result FutureKeyType, err Error);

	// Redis SET command.
	Set (key string, arg1 []byte) (status FutureBool, err Error);

	// Redis SAVE command.
	Save () (status FutureBool, err Error);

	// Redis KEYS command using "*" wildcard 
	AllKeys () (result FutureKeys, err Error);

	// Redis KEYS command.
	Keys (key string) (result FutureKeys, err Error);

	// Redis EXISTS command.
	Exists (key string) (result FutureBool, err Error);

	// Redis RENAME command.
	Rename (key string, arg1 string) (status FutureBool, err Error);

	// Redis INFO command.
	Info () (result FutureInfo, err Error);

	// Redis PING command.
	Ping () (status FutureBool, err Error);

	// Redis QUIT command.
	Quit () (status FutureBool, err Error);

	// Redis SETNX command.
	Setnx (key string, arg1 []byte) (result FutureBool, err Error);

	// Redis GETSET command.
	Getset (key string, arg1 []byte) (result FutureBytes, err Error);

	// Redis MGET command.
	Mget (key string, arg1 []string) (result FutureBytesArray, err Error);

	// Redis INCR command.
	Incr (key string) (result FutureInt64, err Error);

	// Redis INCRBY command.
	Incrby (key string, arg1 int64) (result FutureInt64, err Error);

	// Redis DECR command.
	Decr (key string) (result FutureInt64, err Error);

	// Redis DECRBY command.
	Decrby (key string, arg1 int64) (result FutureInt64, err Error);

	// Redis DEL command.
	Del (key string) (result FutureBool, err Error);

	// Redis RANDOMKEY command.
	Randomkey () (result FutureString, err Error);

	// Redis RENAMENX command.
	Renamenx (key string, arg1 string) (result FutureBool, err Error);

	// Redis DBSIZE command.
	Dbsize () (result FutureInt64, err Error);

	// Redis EXPIRE command.
	Expire (key string, arg1 int64) (result FutureBool, err Error);

	// Redis TTL command.
	Ttl (key string) (result FutureInt64, err Error);

	// Redis RPUSH command.
	Rpush (key string, arg1 []byte) (status FutureBool, err Error);

	// Redis LPUSH command.
	Lpush (key string, arg1 []byte) (status FutureBool, err Error);

	// Redis LSET command.
	Lset (key string, arg1 int64, arg2 []byte) (status FutureBool, err Error);

	// Redis LREM command.
	Lrem (key string, arg1 []byte, arg2 int64) (result FutureInt64, err Error);

	// Redis LLEN command.
	Llen (key string) (result FutureInt64, err Error);

	// Redis LRANGE command.
	Lrange (key string, arg1 int64, arg2 int64) (result FutureBytesArray, err Error);

	// Redis LTRIM command.
	Ltrim (key string, arg1 int64, arg2 int64) (status FutureBool, err Error);

	// Redis LINDEX command.
	Lindex (key string, arg1 int64) (result FutureBytes, err Error);

	// Redis LPOP command.
	Lpop (key string) (result FutureBytes, err Error);

	// Redis RPOP command.
	Rpop (key string) (result FutureBytes, err Error);

	// Redis RPOPLPUSH command.
	Rpoplpush (key string, arg1 string) (result FutureBytes, err Error);

	// Redis SADD command.
	Sadd (key string, arg1 []byte) (result FutureBool, err Error);

	// Redis SREM command.
	Srem (key string, arg1 []byte) (result FutureBool, err Error);

	// Redis SISMEMBER command.
	Sismember (key string, arg1 []byte) (result FutureBool, err Error);

	// Redis SMOVE command.
	Smove (key string, arg1 string, arg2 []byte) (result FutureBool, err Error);

	// Redis SCARD command.
	Scard (key string) (result FutureInt64, err Error);

	// Redis SINTER command.
	Sinter (key string, arg1 []string) (result FutureBytesArray, err Error);

	// Redis SINTERSTORE command.
	Sinterstore (key string, arg1 []string) (status FutureBool, err Error);

	// Redis SUNION command.
	Sunion (key string, arg1 []string) (result FutureBytesArray, err Error);

	// Redis SUNIONSTORE command.
	Sunionstore (key string, arg1 []string) (status FutureBool, err Error);

	// Redis SDIFF command.
	Sdiff (key string, arg1 []string) (result FutureBytesArray, err Error);

	// Redis SDIFFSTORE command.
	Sdiffstore (key string, arg1 []string) (status FutureBool, err Error);

	// Redis SMEMBERS command.
	Smembers (key string) (result FutureBytesArray, err Error);

	// Redis SRANDMEMBER command.
	Srandmember (key string) (result FutureBytes, err Error);

	// Redis ZADD command.
	Zadd (key string, arg1 float64, arg2 []byte) (result FutureBool, err Error);

	// Redis ZREM command.
	Zrem (key string, arg1 []byte) (result FutureBool, err Error);

	// Redis ZCARD command.
	Zcard (key string) (result FutureInt64, err Error);

	// Redis ZSCORE command.
	Zscore (key string, arg1 []byte) (result FutureFloat64, err Error);

	// Redis ZRANGE command.
	Zrange (key string, arg1 int64, arg2 int64) (result FutureBytesArray, err Error);

	// Redis ZREVRANGE command.
	Zrevrange (key string, arg1 int64, arg2 int64) (result FutureBytesArray, err Error);

	// Redis ZRANGEBYSCORE command.
	Zrangebyscore (key string, arg1 float64, arg2 float64) (result FutureBytesArray, err Error);

	// Redis FLUSHDB command.
	Flushdb () (status FutureBool, err Error);

	// Redis FLUSHALL command.
	Flushall () (status FutureBool, err Error);

	// Redis MOVE command.
	Move (key string, arg1 int64) (result FutureBool, err Error);

	// Redis BGSAVE command.
	Bgsave () (status FutureBool, err Error);

	// Redis LASTSAVE command.
	Lastsave () (result FutureInt64, err Error);
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
func init () { 
//	runtime.GOMAXPROCS(2);
//	flag.Parse(); 
}

// redis:d
//
// global debug flag for redis package components.
// 
var _debug *bool = flag.Bool ("redis:d", false, "debug flag for go-redis"); // TEMP: should default to false
func debug() bool { return *_debug; }

