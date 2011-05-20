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

package redis

// ----------------------------------------------------------------------------
// PROTOCOL SPEC
//
// Various other elements of the package use the artifacts of this file to
// negotiate the Redis protocol. 
//
// Redis version: 1.n
// ----------------------------------------------------------------------------
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
	OK			= true;	 	
	PONG 		= true;		
	ERR 		= false;
)

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
 	KEYS			Command = Command {"KEYS", KEY, 	MULTI_BULK};
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
// TODO	SORT		(RequestType.MULTI_KEY,		ResponseType.MULTI_BULK), 
 	SORT			Command = Command {"SORT", KEY_SPEC, 	MULTI_BULK};
 	SAVE			Command = Command {"SAVE", NO_ARG, 	STATUS};
 	BGSAVE			Command = Command {"BGSAVE", NO_ARG, 	STATUS};
 	LASTSAVE		Command = Command {"LASTSAVE", NO_ARG, 	NUMBER};
 	SHUTDOWN		Command = Command {"SHUTDOWN", NO_ARG, 	VIRTUAL};
 	INFO			Command = Command {"INFO", NO_ARG, 	BULK};
 	MONITOR			Command = Command {"MONITOR", NO_ARG, 	VIRTUAL};
 )

