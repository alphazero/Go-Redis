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
	"os";
	"fmt"
)

// ----------------------------------------------------------------------------
// Interfaces
// ----------------------------------------------------------------------------

// Redis 'KeyType'
type KeyType byte;

// Known Redis key types
const (
	RT_NONE 	KeyType = iota;
	RT_STRING;
	RT_SET;
	RT_LIST;
	RT_ZSET;
)

// Returns KeyType by name
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

// Asynchronous call semantics client interface.
//
// TBD.

type Pipeline interface {
}

// ----------------------------------------------------------------------------
// PROTOCOL SPEC
// These are for internal ops and not of any interest to the end user.  
// ----------------------------------------------------------------------------

// Request type defines the characteristic pattern of a cmd request
//
type RequestType int;
const (
	_ RequestType = iota;
	NO_ARG;
	KEY;
	KEY_KEY;
	KEY_NUM;
	KEY_SPEC;
	KEY_NUM_NUM;
	KEY_VALUE;
	KEY_IDX_VALUE;
	KEY_KEY_VALUE;
	KEY_CNT_VALUE;
	MULTI_KEY;
)

// Response type defines the various flavors of responses from Redis
// per its specification.

type ResponseType int;
const (
	VIRTUAL ResponseType = iota;
	BOOLEAN;
	NUMBER;
	STRING;
	STATUS;
	BULK;
	MULTI_BULK;
)

// Describes a given Redis command
//
type Command struct {
	Code string;
	ReqType RequestType;
	RespType ResponseType;
}

// The supported Command set, with one to one mapping to eponymous Redis command.
//
var (
	AUTH			Command = Command {"AUTH", KEY, 	STATUS};
 	PING			Command = Command {"PING", NO_ARG, 	STATUS};
 	QUIT			Command = Command {"QUIT", NO_ARG, 	VIRTUAL};
 	SET				Command = Command {"SET", KEY_VALUE, 	STATUS};
 	GET				Command = Command {"GET", KEY, 	BULK};
 	GETSET			Command = Command {"GETSET", KEY_VALUE, 	BULK};
 	MGET			Command = Command {"MGET", MULTI_KEY, 	MULTI_BULK};
 	SETNX			Command = Command {"SETNX", KEY_VALUE, 	BOOLEAN};
 	INCR			Command = Command {"INCR", KEY, 	NUMBER};
 	INCRBY			Command = Command {"INCRBY", KEY_NUM, 	NUMBER};
 	DECR			Command = Command {"DECR", KEY, 	NUMBER};
 	DECRBY			Command = Command {"DECRBY", KEY_NUM, 	NUMBER};
 	EXISTS			Command = Command {"EXISTS", KEY, 	BOOLEAN};
 	DEL				Command = Command {"DEL", KEY, 	BOOLEAN};
 	TYPE			Command = Command {"TYPE", KEY, 	STRING};
 	KEYS			Command = Command {"KEYS", KEY, 	BULK};
 	RANDOMKEY		Command = Command {"RANDOMKEY", NO_ARG, 	STRING};
 	RENAME			Command = Command {"RENAME", KEY_KEY, 	STATUS};
 	RENAMENX		Command = Command {"RENAMENX", KEY_KEY, 	BOOLEAN};
 	DBSIZE			Command = Command {"DBSIZE", NO_ARG, 	NUMBER};
 	EXPIRE			Command = Command {"EXPIRE", KEY_NUM, 	BOOLEAN};
 	TTL				Command = Command {"TTL", KEY, 	NUMBER};
 	RPUSH			Command = Command {"RPUSH", KEY_VALUE, 	STATUS};
 	LPUSH			Command = Command {"LPUSH", KEY_VALUE, 	STATUS};
 	LLEN			Command = Command {"LLEN", KEY, 	NUMBER};
 	LRANGE			Command = Command {"LRANGE", KEY_NUM_NUM, 	MULTI_BULK};
 	LTRIM			Command = Command {"LTRIM", KEY_NUM_NUM, 	STATUS};
 	LINDEX			Command = Command {"LINDEX", KEY_NUM, 	BULK};
 	LSET			Command = Command {"LSET", KEY_IDX_VALUE, 	STATUS};
 	LREM			Command = Command {"LREM", KEY_CNT_VALUE, 	NUMBER};
 	LPOP			Command = Command {"LPOP", KEY, 	BULK};
 	RPOP			Command = Command {"RPOP", KEY, 	BULK};
 	RPOPLPUSH		Command = Command {"RPOPLPUSH", KEY_VALUE, 	BULK};
 	SADD			Command = Command {"SADD", KEY_VALUE, 	BOOLEAN};
 	SREM			Command = Command {"SREM", KEY_VALUE, 	BOOLEAN};
 	SCARD			Command = Command {"SCARD", KEY, 	NUMBER};
 	SISMEMBER		Command = Command {"SISMEMBER", KEY_VALUE, 	BOOLEAN};
 	SINTER			Command = Command {"SINTER", MULTI_KEY, 	MULTI_BULK};
 	SINTERSTORE		Command = Command {"SINTERSTORE", MULTI_KEY, 	STATUS};
 	SUNION			Command = Command {"SUNION", MULTI_KEY, 	MULTI_BULK};
 	SUNIONSTORE		Command = Command {"SUNIONSTORE", MULTI_KEY, 	STATUS};
 	SDIFF			Command = Command {"SDIFF", MULTI_KEY, 	MULTI_BULK};
 	SDIFFSTORE		Command = Command {"SDIFFSTORE", MULTI_KEY, 	STATUS};
 	SMEMBERS		Command = Command {"SMEMBERS", KEY, 	MULTI_BULK};
 	SMOVE			Command = Command {"SMOVE", KEY_KEY_VALUE, 	BOOLEAN};
 	SRANDMEMBER		Command = Command {"SRANDMEMBER", KEY, 	BULK};
 	ZADD			Command = Command {"ZADD", KEY_IDX_VALUE, 	BOOLEAN};
 	ZREM			Command = Command {"ZREM", KEY_VALUE, 	BOOLEAN};
 	ZCARD			Command = Command {"ZCARD", KEY, 	NUMBER};
 	ZSCORE			Command = Command {"ZSCORE", KEY_VALUE, 	BULK};
 	ZRANGE			Command = Command {"ZRANGE", KEY_NUM_NUM, 	MULTI_BULK};
 	ZREVRANGE		Command = Command {"ZREVRANGE", KEY_NUM_NUM, 	MULTI_BULK};
 	ZRANGEBYSCORE	Command = Command {"ZRANGEBYSCORE", KEY_NUM_NUM, 	MULTI_BULK};
 	SELECT			Command = Command {"SELECT", KEY, 	STATUS};
 	FLUSHDB			Command = Command {"FLUSHDB", NO_ARG, 	STATUS};
 	FLUSHALL		Command = Command {"FLUSHALL", NO_ARG, 	STATUS};
 	MOVE			Command = Command {"MOVE", KEY_NUM, 	BOOLEAN};
 	SORT			Command = Command {"SORT", KEY_SPEC, 	MULTI_BULK};
 	SAVE			Command = Command {"SAVE", NO_ARG, 	STATUS};
 	BGSAVE			Command = Command {"BGSAVE", NO_ARG, 	STATUS};
 	LASTSAVE		Command = Command {"LASTSAVE", NO_ARG, 	NUMBER};
 	SHUTDOWN		Command = Command {"SHUTDOWN", NO_ARG, 	VIRTUAL};
 	INFO			Command = Command {"INFO", NO_ARG, 	BULK};
 	MONITOR			Command = Command {"MONITOR", NO_ARG, 	VIRTUAL};
 )


// ----------------------------------------------------------------------------
// ERRORS
// ----------------------------------------------------------------------------

// redis.Errors can be either user or system category.
//  
// REDIS_ERR are sent by the Redis server in response to user commands. (For example
// "operation against key holding the wrong value).  If a call to the client method
// has a non nil error with this cateogry, it means you did something wrong.  If you
// want to know the details, see Error.GetMessage() which returns the precise error
// sent from Redis.
//
// If error is of category SYSTEM_ERR, it either means there were errors during exeuction
// (which may be due to your system environment, networking, memory, etc.) OR, there is
// a bug in the package.   Unlike REDIS_ERR cateogry, SYSTEM_ERR category errors may have
// and underlying cause (os.Error) which can be obtained by Error.Cause() method.

type ErrorCategory uint8;
const (
	_ ErrorCategory = iota;
	REDIS_ERR;
	SYSTEM_ERR;
)

// Defines the interfce to get details of an Error.
// 
type Error interface {
	Cause() os.Error;
	Category() ErrorCategory;
	Message() string;
	String() string;
}
type redisError struct {
	msg 	  string;
	category  ErrorCategory;
	cause     os.Error;
}
func (e redisError) Cause() os.Error { return e.cause; }
func (e redisError) Category() ErrorCategory { return e.category; }
func (e redisError) Message() string { return e.msg; }
func (e redisError) String() string {
	var errCat string;
	var causeDetails string = "";
	switch e.category {
	case REDIS_ERR: 
		errCat = "REDIS-ERROR";
	case SYSTEM_ERR: 
		errCat = "SYSTEM-ERROR";
		if e.Cause() != nil { 
			causeDetails = fmt.Sprintf(" [cause: %s]", e.Cause()); 
		}
	}
	return fmt.Sprintf("[go-redis|%d|%s]: %s %s", e.category, errCat, e.msg, causeDetails);
}

// Creates an Error of REDIS_ERR category with the message.
//
func NewRedisError(msg string) Error {
	e := new (redisError);
	e.msg = msg;
	e.category = REDIS_ERR;
	return e;
}

// Creates an Error of specified category with the message.
//
func NewError(t ErrorCategory, msg string) Error {
	e := new (redisError);
	e.msg = msg;
	e.category = t;
	return e;
}

// Creates an Error of specified category with the message and cause.
//
func NewErrorWithCause(cat ErrorCategory, msg string, cause os.Error) Error {
	e := new (redisError);
	e.msg = msg;
	e.category = cat;
	e.cause = cause;
	return e;
}

