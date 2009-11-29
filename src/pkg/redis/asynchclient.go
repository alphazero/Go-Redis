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

import (
	"strings";
	"os";
	"log";
	"fmt";
/*
	"bytes";
	"strconv";
*/
)

type async struct {
	conn  AsyncConnection;
}




// Create a new Client and connects to the Redis server using the
// default ConnectionSpec.
//
func NewAsynchClient () (c AsyncClient, err os.Error){
	spec := DefaultSpec();
	c, err = NewAsynchClientWithSpec(spec);
	return;
}

// Create a new Client and connects to the Redis server using the
// specified ConnectionSpec.
//
func NewAsynchClientWithSpec (spec *ConnectionSpec) (c AsyncClient, err os.Error) {
	_c := new(async);
	_c.conn, err = NewAsynchConnection (spec);
	if err != nil {
		if debug() {log.Stderr("NewAsyncConnection() raised error: ", err);}
		return nil, err;
	}
	return _c, nil;
}
/** ----------------- REDIS INTERFACE ------------- **/

// Redis GET command.
func (c *async) Get (arg0 string) (result FutureBytes, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&GET, [][]byte{arg0bytes});
	if err == nil {result = resp.future.(FutureBytes);}
	return result, err;

}

// Redis TYPE command.
func (c *async) Type (arg0 string) (result FutureKeyType, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&TYPE, [][]byte{arg0bytes});
	if err == nil {result = newFutureKeyType(resp.future.(FutureString));}
	return result, err;
}

// Redis SET command.
func (c *async) Set (arg0 string, arg1 []byte) (stat FutureBool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;

	resp, err := c.conn.QueueRequest(&SET, [][]byte{arg0bytes, arg1bytes});
	if err == nil {stat = resp.future.(FutureBool);}

	return;
}

// Redis SAVE command.
func (c *async) Save () (stat FutureBool, err Error){
	resp, err := c.conn.QueueRequest(&SAVE, [][]byte{});
	if err == nil {stat = resp.future.(FutureBool);}

	return;

}

// Redis KEYS command.
func (c *async) AllKeys () (result FutureKeys, err Error){
	return c.Keys("*");
}

// Redis KEYS command.
func (c *async) Keys (arg0 string) (result FutureKeys, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&KEYS, [][]byte{arg0bytes});
	if err == nil {
		result = newFutureKeys(resp.future.(FutureBytes));
	}
	return result, err;
}

/***
// Redis SORT command.
func (c *async) Sort (arg0 string) (result redis.Sort, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&SORT, [][]byte{arg0bytes});
	if err == nil {result = resp.GetMultiBulkData();}
	return result, err;

}
***/
// Redis EXISTS command.
func (c *async) Exists (arg0 string) (result FutureBool, err Error){
	arg0bytes := strings.Bytes (arg0);

	resp, err := c.conn.QueueRequest(&EXISTS, [][]byte{arg0bytes});
	if err == nil {result = resp.future.(FutureBool);}
	return result, err;

}

// Redis RENAME command.
func (c *async) Rename (arg0 string, arg1 string) (stat FutureBool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (arg1);

	resp, err := c.conn.QueueRequest(&RENAME, [][]byte{arg0bytes, arg1bytes});
	if err == nil {stat = resp.future.(FutureBool);}

	return;
}

// Redis INFO command.
func (c *async) Info () (result FutureInfo, err Error){
	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&INFO, [][]byte{});
	if err == nil {
		result = newFutureInfo(resp.future.(FutureBytes));
	}
	return result, err;
}

// Redis PING command.
func (c *async) Ping () (stat FutureBool, err Error){
	resp, err := c.conn.QueueRequest(&PING, [][]byte{});
	if err == nil {stat = resp.future.(FutureBool);}

	return;
}

// Redis QUIT command.
func (c *async) Quit () (stat FutureBool, err Error){
//	c.conn.Close();
////	resp, err = c.conn.QueueRequest(&QUIT);
//	if err == nil {stat = resp.future.(FutureBool);}

	return nil, NewError(SYSTEM_ERR, "<BUG> Lazy programmer hasn't implemented Quit!");;
}

// Redis SETNX command.
func (c *async) Setnx (arg0 string, arg1 []byte) (result FutureBool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;

	resp, err := c.conn.QueueRequest(&SETNX, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.future.(FutureBool);}
	return result, err;

}

// Redis GETSET command.
func (c *async) Getset (arg0 string, arg1 []byte) (result FutureBytes, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&GETSET, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.future.(FutureBytes);}
	return result, err;

}

// Redis MGET command.
func (c *async) Mget (arg0 string, arg1 []string) (result FutureBytesArray, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := concatAndGetBytes(arg1, " ");

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&MGET, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.future.(FutureBytesArray);}
	return result, err;

}


// Redis INCR command.
func (c *async) Incr (arg0 string) (result FutureInt64, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&INCR, [][]byte{arg0bytes});
	if err == nil {result = resp.future.(FutureInt64);}
	return result, err;

}

// Redis INCRBY command.
func (c *async) Incrby (arg0 string, arg1 int64) (result FutureInt64, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&INCRBY, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.future.(FutureInt64);}
	return result, err;

}

// Redis DECR command.
func (c *async) Decr (arg0 string) (result FutureInt64, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&DECR, [][]byte{arg0bytes});
	if err == nil {result = resp.future.(FutureInt64);}
	return result, err;

}

// Redis DECRBY command.
func (c *async) Decrby (arg0 string, arg1 int64) (result FutureInt64, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&DECRBY, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.future.(FutureInt64);}
	return result, err;

}

// Redis DEL command.
func (c *async) Del (arg0 string) (result FutureBool, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&DEL, [][]byte{arg0bytes});
	if err == nil {result = resp.future.(FutureBool);}
	return result, err;

}

// Redis RANDOMKEY command.
func (c *async) Randomkey () (result FutureString, err Error){
	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&RANDOMKEY, [][]byte{});
	if err == nil {result = resp.future.(FutureString);}
	return result, err;

}

// Redis RENAMENX command.
func (c *async) Renamenx (arg0 string, arg1 string) (result FutureBool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (arg1);

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&RENAMENX, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.future.(FutureBool);}
	return result, err;

}

// Redis DBSIZE command.
func (c *async) Dbsize () (result FutureInt64, err Error){
	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&DBSIZE, [][]byte{});
	if err == nil {result = resp.future.(FutureInt64);}
	return result, err;

}

// Redis EXPIRE command.
func (c *async) Expire (arg0 string, arg1 int64) (result FutureBool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&EXPIRE, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.future.(FutureBool);}
	return result, err;

}

// Redis TTL command.
func (c *async) Ttl (arg0 string) (result FutureInt64, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&TTL, [][]byte{arg0bytes});
	if err == nil {result = resp.future.(FutureInt64);}
	return result, err;

}

// Redis RPUSH command.
func (c *async) Rpush (arg0 string, arg1 []byte) (stat FutureBool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;

	resp, err := c.conn.QueueRequest(&RPUSH, [][]byte{arg0bytes, arg1bytes});
	if err == nil {stat = resp.future.(FutureBool);}

	return;
}

// Redis LPUSH command.
func (c *async) Lpush (arg0 string, arg1 []byte) (result FutureBool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;

	resp, err := c.conn.QueueRequest(&LPUSH, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.future.(FutureBool);}

	return;
}

// Redis LSET command.
func (c *async) Lset (arg0 string, arg1 int64, arg2 []byte) (stat FutureBool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));
	arg2bytes := arg2;

	resp, err := c.conn.QueueRequest(&LSET, [][]byte{arg0bytes, arg1bytes, arg2bytes});
	if err == nil {stat = resp.future.(FutureBool);}

	return;
}


// Redis LREM command.
func (c *async) Lrem (arg0 string, arg1 []byte, arg2 int64) (result FutureInt64, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;
	arg2bytes := strings.Bytes (fmt.Sprintf("%d", arg2));

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&LREM, [][]byte{arg0bytes, arg1bytes, arg2bytes});
	if err == nil {result = resp.future.(FutureInt64);}
	return result, err;

}

// Redis LLEN command.
func (c *async) Llen (arg0 string) (result FutureInt64, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&LLEN, [][]byte{arg0bytes});
	if err == nil {result = resp.future.(FutureInt64);}
	return result, err;

}

// Redis LRANGE command.
func (c *async) Lrange (arg0 string, arg1 int64, arg2 int64) (result FutureBytesArray, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));
	arg2bytes := strings.Bytes (fmt.Sprintf("%d", arg2));

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&LRANGE, [][]byte{arg0bytes, arg1bytes, arg2bytes});
	if err == nil {result = resp.future.(FutureBytesArray);}
	return result, err;

}

// Redis LTRIM command.
func (c *async) Ltrim (arg0 string, arg1 int64, arg2 int64) (stat FutureBool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));
	arg2bytes := strings.Bytes (fmt.Sprintf("%d", arg2));

	resp, err := c.conn.QueueRequest(&LTRIM, [][]byte{arg0bytes, arg1bytes, arg2bytes});
	if err == nil {stat = resp.future.(FutureBool);}

	return;
}

// Redis LINDEX command.
func (c *async) Lindex (arg0 string, arg1 int64) (result FutureBytes, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&LINDEX, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.future.(FutureBytes);}
	return result, err;

}

// Redis LPOP command.
func (c *async) Lpop (arg0 string) (result FutureBytes, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&LPOP, [][]byte{arg0bytes});
	if err == nil {result = resp.future.(FutureBytes);}
	return result, err;

}

// Redis RPOP command.
func (c *async) Rpop (arg0 string) (result FutureBytes, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&RPOP, [][]byte{arg0bytes});
	if err == nil {result = resp.future.(FutureBytes);}
	return result, err;

}

// Redis RPOPLPUSH command.
func (c *async) Rpoplpush (arg0 string, arg1 string) (result FutureBytes, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (arg1);

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&RPOPLPUSH, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.future.(FutureBytes);}
	return result, err;

}

// Redis SADD command.
func (c *async) Sadd (arg0 string, arg1 []byte) (result FutureBool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&SADD, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.future.(FutureBool);}
	return result, err;

}


// Redis SREM command.
func (c *async) Srem (arg0 string, arg1 []byte) (result FutureBool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&SREM, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.future.(FutureBool);}
	return result, err;

}


// Redis SISMEMBER command.
func (c *async) Sismember (arg0 string, arg1 []byte) (result FutureBool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&SISMEMBER, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.future.(FutureBool);}
	return result, err;

}

// Redis SMOVE command.
func (c *async) Smove (arg0 string, arg1 string, arg2 []byte) (result FutureBool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (arg1);
	arg2bytes := arg2;

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&SMOVE, [][]byte{arg0bytes, arg1bytes, arg2bytes});
	if err == nil {result = resp.future.(FutureBool);}
	return result, err;

}

// Redis SCARD command.
func (c *async) Scard (arg0 string) (result FutureInt64, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&SCARD, [][]byte{arg0bytes});
	if err == nil {result = resp.future.(FutureInt64);}
	return result, err;

}

// Redis SINTER command.
func (c *async) Sinter (arg0 string, arg1 []string) (result FutureBytesArray, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := concatAndGetBytes(arg1, " ");

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&SINTER, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.future.(FutureBytesArray);}
	return result, err;

}

// Redis SINTERSTORE command.
func (c *async) Sinterstore (arg0 string, arg1 []string) (stat FutureBool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := concatAndGetBytes(arg1, " ");

	resp, err := c.conn.QueueRequest(&SINTERSTORE, [][]byte{arg0bytes, arg1bytes});
	if err == nil {stat = resp.future.(FutureBool);}

	return;
}

// Redis SUNION command.
func (c *async) Sunion (arg0 string, arg1 []string) (result FutureBytesArray, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := concatAndGetBytes(arg1, " ");

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&SUNION, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.future.(FutureBytesArray);}
	return result, err;

}

// Redis SUNIONSTORE command.
func (c *async) Sunionstore (arg0 string, arg1 []string) (stat FutureBool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := concatAndGetBytes(arg1, " ");

	resp, err := c.conn.QueueRequest(&SUNIONSTORE, [][]byte{arg0bytes, arg1bytes});
	if err == nil {stat = resp.future.(FutureBool);}

	return;
}

// Redis SDIFF command.
func (c *async) Sdiff (arg0 string, arg1 []string) (result FutureBytesArray, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := concatAndGetBytes(arg1, " ");

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&SDIFF, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.future.(FutureBytesArray);}
	return result, err;

}

// Redis SDIFFSTORE command.
func (c *async) Sdiffstore (arg0 string, arg1 []string) (stat FutureBool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := concatAndGetBytes(arg1, " ");

	resp, err := c.conn.QueueRequest(&SDIFFSTORE, [][]byte{arg0bytes, arg1bytes});
	if err == nil {stat = resp.future.(FutureBool);}

	return;
}

// Redis SMEMBERS command.
func (c *async) Smembers (arg0 string) (result FutureBytesArray, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&SMEMBERS, [][]byte{arg0bytes});
	if err == nil {result = resp.future.(FutureBytesArray);}
	return result, err;

}

// Redis SRANDMEMBER command.
func (c *async) Srandmember (arg0 string) (result FutureBytes, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&SRANDMEMBER, [][]byte{arg0bytes});
	if err == nil {result = resp.future.(FutureBytes);}
	return result, err;

}

// Redis ZADD command.
func (c *async) Zadd (arg0 string, arg1 float64, arg2 []byte) (result FutureBool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%e", arg1));
	arg2bytes := arg2;

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&ZADD, [][]byte{arg0bytes, arg1bytes, arg2bytes});
	if err == nil {result = resp.future.(FutureBool);}
	return result, err;

}

// Redis ZREM command.
func (c *async) Zrem (arg0 string, arg1 []byte) (result FutureBool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&ZREM, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.future.(FutureBool);}
	return result, err;

}

// Redis ZCARD command.
func (c *async) Zcard (arg0 string) (result FutureInt64, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&ZCARD, [][]byte{arg0bytes});
	if err == nil {result = resp.future.(FutureInt64);}
	return result, err;

}

// Redis ZSCORE command.
func (c *async) Zscore (arg0 string, arg1 []byte) (result FutureFloat64, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&ZSCORE, [][]byte{arg0bytes, arg1bytes});
	if err == nil {
		result = newFutureFloat64(resp.future.(FutureBytes));
	}
	return result, err;

}

// Redis ZRANGE command.
func (c *async) Zrange (arg0 string, arg1 int64, arg2 int64) (result FutureBytesArray, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));
	arg2bytes := strings.Bytes (fmt.Sprintf("%d", arg2));

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&ZRANGE, [][]byte{arg0bytes, arg1bytes, arg2bytes});
	if err == nil {result = resp.future.(FutureBytesArray);}
	return result, err;

}

// Redis ZREVRANGE command.
func (c *async) Zrevrange (arg0 string, arg1 int64, arg2 int64) (result FutureBytesArray, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));
	arg2bytes := strings.Bytes (fmt.Sprintf("%d", arg2));

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&ZREVRANGE, [][]byte{arg0bytes, arg1bytes, arg2bytes});
	if err == nil {result = resp.future.(FutureBytesArray);}
	return result, err;

}

// Redis ZRANGEBYSCORE command.
func (c *async) Zrangebyscore (arg0 string, arg1 float64, arg2 float64) (result FutureBytesArray, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%e", arg1));
	arg2bytes := strings.Bytes (fmt.Sprintf("%e", arg2));

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&ZRANGEBYSCORE, [][]byte{arg0bytes, arg1bytes, arg2bytes});
	if err == nil {result = resp.future.(FutureBytesArray);}
	return result, err;

}

// Redis FLUSHDB command.
func (c *async) Flushdb () (stat FutureBool, err Error){
	resp, err := c.conn.QueueRequest(&FLUSHDB, [][]byte{});
	if err == nil {stat = resp.future.(FutureBool);}

	return;
}

// Redis FLUSHALL command.
func (c *async) Flushall () (stat FutureBool, err Error){
	resp, err := c.conn.QueueRequest(&FLUSHALL, [][]byte{});
	if err == nil {stat = resp.future.(FutureBool);}

	return;
}

// Redis MOVE command.
func (c *async) Move (arg0 string, arg1 int64) (result FutureBool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));

	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&MOVE, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.future.(FutureBool);}
	return result, err;

}

// Redis BGSAVE command.
func (c *async) Bgsave () (stat FutureBool, err Error){
	resp, err := c.conn.QueueRequest(&BGSAVE, [][]byte{});
	if err == nil {stat = resp.future.(FutureBool);}

	return;
}

// Redis LASTSAVE command.
func (c *async) Lastsave () (result FutureInt64, err Error){
	var resp *PendingResponse;
	resp, err = c.conn.QueueRequest(&LASTSAVE, [][]byte{});
	if err == nil {result = resp.future.(FutureInt64);}
	return result, err;

}
