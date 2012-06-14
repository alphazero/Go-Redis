//   Copyright 2009-2012 Joubin Houshyar
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

// Package redis provides both clients and connectors for the Redis
// server.  Both synchronous and asynchronous interaction modes are
// supported.  Asynchronous clients (using the asynchronous connection)
// use pipelining.
//
// Synchronous semantics are defined by redis.Client interface type
//
//
// Usage example:
//
//  func usingRedisSync () Error {
//      spec := DefaultConnectionSpec();
//      pipeline := NewAsynchClient(spec);
//
//      value, reqErr := pipeline.Get("my-key");
//      if reqErr != nil { return withError (reqErr); }
//  }
//
// Asynchronous semantics are defined by redis.AsyncClient interface type.
// Note that these clients use a pipelining connector and a single instance
// can be used by multiple go routines.  Use of pipelining increases throughput
// at the cost of latency.  If low latency is more important to you than
// throughput, and you require async call semantics, then you should use only
// 1 go routine per AsyncClient connection.
//
// Usage example without timeouts (caller blocks on Get until the response
// from Redis has been processed.)
//
//  func usingRedisAsync () Error {
//      spec := DefaultConnectionSpec();
//      pipeline := NewRedisPipeline(spec);
//
//      // async invoke of GET my-key
//      // futureBytes is a FutureBytes that will have the result of the
//      // Redis GET operation.
//      futureBytes, reqErr := pipline.Get("my-key");
//      if reqErr != nil {
//          return withError (reqErr);
//      }
//
//      // ... note that you could issue additional redis commands here ...
//
//      []byte, execErr := futureBytes.Get();
//      if execErr != nil {
//          return withError (execErr);
//      }
//  }
//
// Usage example with timeouts - same Redis op as above but here we use
// TryGet on the Future result with a timeout of 1 msecs:
//
//  func usingRedisAsync () Error {
//      spec := DefaultConnectionSpec();
//      pipeline := NewRedisPipeline(spec);
//
//      // futureBytes is a FutureBytes
//      futureBytes, reqErr := pipline.Get("my-key");
//      if reqErr != nil { return withError (reqErr); }
//
//      // ... note that you could issue additional redis commands here ...
//
//      []byte, execErr := futureBytes.Get();
//      if execErr != nil {
//          return withError (execErr);
//      }
//
//      timeout := 1000000; // wait 1 msec for result
//      []byte, execErr, ok := futureBytes.TryGet (timeout);
//      if !ok {
//          .. handle timeout here
//      }
//      else {
//          if execErr != nil {
//              return withError (execErr);
//          }
//      }
//  }
//
package redis

import (
	"flag"
	"runtime";
)

// Common interface supported by all clients
// to consolidate common ops
type RedisClient interface {

	// Redis QUIT command.
	Quit() (err Error)
}

// The synchronous call semantics Client interface.
//
// Method names map one to one to the Redis command set.
// All methods may return an redis.Error, which is either
// a system error (runtime issue or bug) or Redis error (i.e. user error)
// See Error in this package for details of its interface.
//
// The synchronous client interface provides blocking call semantics supported by
// a distinct request/reply sequence at the connector level.
//
// Method names map one to one to the Redis command set.
//
// All methods may return an redis.Error, which is either a Redis error (from
// the server), or a system error indicating a runtime issue (or bug).
// See Error in this package for details of its interface.
type Client interface {

	// psuedo inheritance to coerce to RedisClient type
	// for the common API (not "algorithm", not "data", but interface ...)
	RedisClient() RedisClient

	// Redis GET command.
	Get(key string) (result []byte, err Error)

	// Redis TYPE command.
	Type(key string) (result KeyType, err Error)

	// Redis SET command.
	Set(key string, arg1 []byte) Error

	// Redis SAVE command.
	Save() Error

	// Redis KEYS command using "*" wildcard 
	AllKeys() (result []string, err Error)

	// Redis KEYS command.
	Keys(key string) (result []string, err Error)

	// Redis EXISTS command.
	Exists(key string) (result bool, err Error)

	// Redis RENAME command.
	Rename(key string, arg1 string) Error

	// Redis INFO command.
	Info() (result map[string]string, err Error)

	// Redis PING command.
	Ping() Error

	// Redis QUIT command.
	Quit() Error

	// Redis SETNX command.
	Setnx(key string, arg1 []byte) (result bool, err Error)

	// Redis GETSET command.
	Getset(key string, arg1 []byte) (result []byte, err Error)

	// Redis MGET command.
	Mget(key string, arg1 []string) (result [][]byte, err Error)

	// Redis INCR command.
	Incr(key string) (result int64, err Error)

	// Redis INCRBY command.
	Incrby(key string, arg1 int64) (result int64, err Error)

	// Redis DECR command.
	Decr(key string) (result int64, err Error)

	// Redis DECRBY command.
	Decrby(key string, arg1 int64) (result int64, err Error)

	// Redis DEL command.
	Del(key string) (result bool, err Error)

	// Redis RANDOMKEY command.
	Randomkey() (result string, err Error)

	// Redis RENAMENX command.
	Renamenx(key string, arg1 string) (result bool, err Error)

	// Redis DBSIZE command.
	Dbsize() (result int64, err Error)

	// Redis EXPIRE command.
	Expire(key string, arg1 int64) (result bool, err Error)

	// Redis TTL command.
	Ttl(key string) (result int64, err Error)

	// Redis RPUSH command.
	Rpush(key string, arg1 []byte) Error

	// Redis LPUSH command.
	Lpush(key string, arg1 []byte) Error

	// Redis LSET command.
	Lset(key string, arg1 int64, arg2 []byte) Error

	// Redis LREM command.
	Lrem(key string, arg1 []byte, arg2 int64) (result int64, err Error)

	// Redis LLEN command.
	Llen(key string) (result int64, err Error)

	// Redis LRANGE command.
	Lrange(key string, arg1 int64, arg2 int64) (result [][]byte, err Error)

	// Redis LTRIM command.
	Ltrim(key string, arg1 int64, arg2 int64) Error

	// Redis LINDEX command.
	Lindex(key string, arg1 int64) (result []byte, err Error)

	// Redis LPOP command.
	Lpop(key string) (result []byte, err Error)

	// Redis BLPOP command.
	Blpop(key string, timeout int) (result [][]byte, err Error)

	// Redis RPOP command.
	Rpop(key string) (result []byte, err Error)

	// Redis BRPOP command.
	Brpop(key string, timeout int) (result [][]byte, err Error)

	// Redis RPOPLPUSH command.
	Rpoplpush(key string, arg1 string) (result []byte, err Error)

	// Redis BRPOPLPUSH command.
	Brpoplpush(key string, arg1 string, timeout int) (result [][]byte, err Error)

	// Redis SADD command.
	Sadd(key string, arg1 []byte) (result bool, err Error)

	// Redis SREM command.
	Srem(key string, arg1 []byte) (result bool, err Error)

	// Redis SISMEMBER command.
	Sismember(key string, arg1 []byte) (result bool, err Error)

	// Redis SMOVE command.
	Smove(key string, arg1 string, arg2 []byte) (result bool, err Error)

	// Redis SCARD command.
	Scard(key string) (result int64, err Error)

	// Redis SINTER command.
	Sinter(key string, arg1 []string) (result [][]byte, err Error)

	// Redis SINTERSTORE command.
	Sinterstore(key string, arg1 []string) Error

	// Redis SUNION command.
	Sunion(key string, arg1 []string) (result [][]byte, err Error)

	// Redis SUNIONSTORE command.
	Sunionstore(key string, arg1 []string) Error

	// Redis SDIFF command.
	Sdiff(key string, arg1 []string) (result [][]byte, err Error)

	// Redis SDIFFSTORE command.
	Sdiffstore(key string, arg1 []string) Error

	// Redis SMEMBERS command.
	Smembers(key string) (result [][]byte, err Error)

	// Redis SRANDMEMBER command.
	Srandmember(key string) (result []byte, err Error)

	// Redis ZADD command.
	Zadd(key string, arg1 float64, arg2 []byte) (result bool, err Error)

	// Redis ZREM command.
	Zrem(key string, arg1 []byte) (result bool, err Error)

	// Redis ZCARD command.
	Zcard(key string) (result int64, err Error)

	// Redis ZSCORE command.
	Zscore(key string, arg1 []byte) (result float64, err Error)

	// Redis ZRANGE command.
	Zrange(key string, arg1 int64, arg2 int64) (result [][]byte, err Error)

	// Redis ZREVRANGE command.
	Zrevrange(key string, arg1 int64, arg2 int64) (result [][]byte, err Error)

	// Redis ZRANGEBYSCORE command.
	Zrangebyscore(key string, arg1 float64, arg2 float64) (result [][]byte, err Error)

	// Redis HGET command.
	Hget(key string, hashkey string) (result []byte, err Error)

	// Redis HSET command.
	Hset(key string, hashkey string, arg1 []byte) Error

	// Redis HGETALL command.
	Hgetall(key string) (result [][]byte, err Error)

	// Redis FLUSHDB command.
	Flushdb() Error

	// Redis FLUSHALL command.
	Flushall() Error

	// Redis MOVE command.
	Move(key string, arg1 int64) (result bool, err Error)

	// Redis BGSAVE command.
	Bgsave() Error

	// Redis LASTSAVE command.
	Lastsave() (result int64, err Error)
}

// The asynchronous client interface provides asynchronous call semantics with
// future results supporting both blocking and try-and-timeout result accessors.
//
// Each method provides a type-safe future result return value, in addition to
// any (system) errors encountered in queuing the request.
//
// The returned value may be ignored by clients that are not interested in the
// future response (for example on SET("foo", data)).  Alternatively, the caller
// may retain the future result referenced and perform blocking and/or timed wait
// gets on the expected response.
//
// Get() or TryGet() on the future result will return any Redis errors that were sent by
// the server, or, Go-Redis (system) errors encountered in processing the response.
type AsyncClient interface {

	// psuedo inheritance to coerce to RedisClient type
	// for the common API (not "algorithm", not "data", but interface ...)
	RedisClient() RedisClient

	// Redis GET command.
	Get(key string) (result FutureBytes, err Error)

	// Redis TYPE command.
	Type(key string) (result FutureKeyType, err Error)

	// Redis SET command.
	Set(key string, arg1 []byte) (status FutureBool, err Error)

	// Redis SAVE command.
	Save() (status FutureBool, err Error)

	// Redis KEYS command using "*" wildcard 
	AllKeys() (result FutureKeys, err Error)

	// Redis KEYS command.
	Keys(key string) (result FutureKeys, err Error)

	// Redis EXISTS command.
	Exists(key string) (result FutureBool, err Error)

	// Redis RENAME command.
	Rename(key string, arg1 string) (status FutureBool, err Error)

	// Redis INFO command.
	Info() (result FutureInfo, err Error)

	// Redis PING command.
	Ping() (status FutureBool, err Error)

	// Redis QUIT command.
	Quit() (err Error)

	// Redis SETNX command.
	Setnx(key string, arg1 []byte) (result FutureBool, err Error)

	// Redis GETSET command.
	Getset(key string, arg1 []byte) (result FutureBytes, err Error)

	// Redis MGET command.
	Mget(key string, arg1 []string) (result FutureBytesArray, err Error)

	// Redis INCR command.
	Incr(key string) (result FutureInt64, err Error)

	// Redis INCRBY command.
	Incrby(key string, arg1 int64) (result FutureInt64, err Error)

	// Redis DECR command.
	Decr(key string) (result FutureInt64, err Error)

	// Redis DECRBY command.
	Decrby(key string, arg1 int64) (result FutureInt64, err Error)

	// Redis DEL command.
	Del(key string) (result FutureBool, err Error)

	// Redis RANDOMKEY command.
	Randomkey() (result FutureString, err Error)

	// Redis RENAMENX command.
	Renamenx(key string, arg1 string) (result FutureBool, err Error)

	// Redis DBSIZE command.
	Dbsize() (result FutureInt64, err Error)

	// Redis EXPIRE command.
	Expire(key string, arg1 int64) (result FutureBool, err Error)

	// Redis TTL command.
	Ttl(key string) (result FutureInt64, err Error)

	// Redis RPUSH command.
	Rpush(key string, arg1 []byte) (status FutureBool, err Error)

	// Redis LPUSH command.
	Lpush(key string, arg1 []byte) (status FutureBool, err Error)

	// Redis LSET command.
	Lset(key string, arg1 int64, arg2 []byte) (status FutureBool, err Error)

	// Redis LREM command.
	Lrem(key string, arg1 []byte, arg2 int64) (result FutureInt64, err Error)

	// Redis LLEN command.
	Llen(key string) (result FutureInt64, err Error)

	// Redis LRANGE command.
	Lrange(key string, arg1 int64, arg2 int64) (result FutureBytesArray, err Error)

	// Redis LTRIM command.
	Ltrim(key string, arg1 int64, arg2 int64) (status FutureBool, err Error)

	// Redis LINDEX command.
	Lindex(key string, arg1 int64) (result FutureBytes, err Error)

	// Redis LPOP command.
	Lpop(key string) (result FutureBytes, err Error)

	// Redis RPOP command.
	Rpop(key string) (result FutureBytes, err Error)

	// Redis RPOPLPUSH command.
	Rpoplpush(key string, arg1 string) (result FutureBytes, err Error)

	// Redis SADD command.
	Sadd(key string, arg1 []byte) (result FutureBool, err Error)

	// Redis SREM command.
	Srem(key string, arg1 []byte) (result FutureBool, err Error)

	// Redis SISMEMBER command.
	Sismember(key string, arg1 []byte) (result FutureBool, err Error)

	// Redis SMOVE command.
	Smove(key string, arg1 string, arg2 []byte) (result FutureBool, err Error)

	// Redis SCARD command.
	Scard(key string) (result FutureInt64, err Error)

	// Redis SINTER command.
	Sinter(key string, arg1 []string) (result FutureBytesArray, err Error)

	// Redis SINTERSTORE command.
	Sinterstore(key string, arg1 []string) (status FutureBool, err Error)

	// Redis SUNION command.
	Sunion(key string, arg1 []string) (result FutureBytesArray, err Error)

	// Redis SUNIONSTORE command.
	Sunionstore(key string, arg1 []string) (status FutureBool, err Error)

	// Redis SDIFF command.
	Sdiff(key string, arg1 []string) (result FutureBytesArray, err Error)

	// Redis SDIFFSTORE command.
	Sdiffstore(key string, arg1 []string) (status FutureBool, err Error)

	// Redis SMEMBERS command.
	Smembers(key string) (result FutureBytesArray, err Error)

	// Redis SRANDMEMBER command.
	Srandmember(key string) (result FutureBytes, err Error)

	// Redis ZADD command.
	Zadd(key string, arg1 float64, arg2 []byte) (result FutureBool, err Error)

	// Redis ZREM command.
	Zrem(key string, arg1 []byte) (result FutureBool, err Error)

	// Redis ZCARD command.
	Zcard(key string) (result FutureInt64, err Error)

	// Redis ZSCORE command.
	Zscore(key string, arg1 []byte) (result FutureFloat64, err Error)

	// Redis ZRANGE command.
	Zrange(key string, arg1 int64, arg2 int64) (result FutureBytesArray, err Error)

	// Redis ZREVRANGE command.
	Zrevrange(key string, arg1 int64, arg2 int64) (result FutureBytesArray, err Error)

	// Redis ZRANGEBYSCORE command.
	Zrangebyscore(key string, arg1 float64, arg2 float64) (result FutureBytesArray, err Error)

	// Redis FLUSHDB command.
	Flushdb() (status FutureBool, err Error)

	// Redis FLUSHALL command.
	Flushall() (status FutureBool, err Error)

	// Redis MOVE command.
	Move(key string, arg1 int64) (result FutureBool, err Error)

	// Redis BGSAVE command.
	Bgsave() (status FutureBool, err Error)

	// Redis LASTSAVE command.
	Lastsave() (result FutureInt64, err Error)
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
func init() {
		runtime.GOMAXPROCS(2);
		flag.Parse();
}

// redis:d
//
// global debug flag for redis package components.
// 
var _debug *bool = flag.Bool("redis:d", false, "debug flag for go-redis") // TEMP: should default to false
func debug() bool {
	return *_debug
}
