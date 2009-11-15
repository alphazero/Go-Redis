package redis

import (
	"os";
	"fmt"
)

// ----------------------------------------------------------------------------
// Interfaces
// ----------------------------------------------------------------------------

type Client interface {

	// Redis GET command.
	Get (key string) (result []byte, err Error);

	// Redis TYPE command.
	Type (key string) (result string, err Error);

	// Redis SET command.
	Set (key string, arg1 []byte) (Error);

	// Redis SAVE command.
	Save () (Error);

	// Redis KEYS command.
	AllKeys () (result [][]byte, err Error);

	// Redis KEYS command.
	Keys (key string) (result [][]byte, err Error);

	// Redis SORT command.
//	Sort (key string) (result Sort, err Error);

	// Redis EXISTS command.
	Exists (key string) (result bool, err Error);

	// Redis RENAME command.
	Rename (key string, arg1 string) (Error);

	// Redis INFO command.
	Info () (result map[string]string, err Error);

	// Redis PING command.
	Ping () (result Client, err Error);

	// Redis QUIT command.
	Quit () (Error);

	// Redis SETNX command.
	Setnx (key string, arg1 []byte) (result bool, err Error);

	// Redis GETSET command.
	Getset (key string, arg1 int64) (result []byte, err Error);

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
	Flushdb () (result Client, err Error);

	// Redis FLUSHALL command.
	Flushall () (result Client, err Error);

	// Redis MOVE command.
	Move (key string, arg1 int64) (result bool, err Error);

	// Redis BGSAVE command.
	Bgsave () (Error);

	// Redis LASTSAVE command.
	Lastsave () (result int64, err Error);
}

type Pipeline interface {
	Sync () (Client, os.Error);
}

// ----------------------------------------------------------------------------
// PROTOCOL SPEC
// ----------------------------------------------------------------------------

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

type Command struct {
	Code string;
	ReqType RequestType;
	RespType ResponseType;
}

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

type ErrorCategory uint8;
const (
	_ ErrorCategory = iota;
	REDIS_ERR;
	SYSTEM_ERR;
//	RUNTIME;
//	BUG;
)
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

func NewRedisError(msg string) Error {
	e := new (redisError);
	e.msg = msg;
	e.category = REDIS_ERR;
	return e;
}
func NewError(t ErrorCategory, msg string) Error {
	e := new (redisError);
	e.msg = msg;
	e.category = t;
	return e;
}
func NewErrorWithCause(cat ErrorCategory, msg string, cause os.Error) Error {
	e := new (redisError);
	e.msg = msg;
	e.category = cat;
	e.cause = cause;
	return e;
}

