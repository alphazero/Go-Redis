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
	Package client provides implementations for Synchronous Redis clients.  
*/
package redis

import (
	"os";
	"strings";
	"bytes";
	"fmt";
	"strconv";
	"log";
)

type synchClient struct {
	conn  SyncConnection;
}

// Create a new Client and connects to the Redis server using the
// default ConnectionSpec.
//
func NewSynchClient () (c Client, err os.Error){
	spec := DefaultSpec();
	c, err = NewSynchClientWithSpec(spec);
	if err != nil { 
		if debug() {log.Stderr("NewSynchClientWithSpec raised error: ", err);}
	}
	if c == nil {
		if debug() { log.Stderr("NewSynchClientWithSpec returned nil Client.");}
		err = os.NewError("NewSynchClientWithSpec returned nil Client.");
	}
	return;
}

// Create a new Client and connects to the Redis server using the
// specified ConnectionSpec.
//
func NewSynchClientWithSpec (spec *ConnectionSpec) (c Client, err os.Error) {
	_c := new(synchClient);
	_c.conn, err = NewSyncConnection (spec);
	if err != nil {
		if debug() {log.Stderr("NewSyncConnection() raised error: ", err);}
		return nil, err;
	}
	return _c, nil;
}

/** ----------------- REDIS INTERFACE ------------- **/

// Redis GET command.
func (c *synchClient) Get (arg0 string) (result []byte, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&GET, [][]byte{arg0bytes});
	if err == nil {result = resp.GetBulkData();}
	return result, err;

}

// Redis TYPE command.
func (c *synchClient) Type (arg0 string) (result KeyType, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&TYPE, [][]byte{arg0bytes});
	if err == nil {result = GetKeyType(resp.GetStringValue());}
	return result, err;
}

// Redis SET command.
func (c *synchClient) Set (arg0 string, arg1 []byte) (err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;

	_, err = c.conn.ServiceRequest(&SET, [][]byte{arg0bytes, arg1bytes});
	return;
}

// Redis SAVE command.
func (c *synchClient) Save () (err Error){
	_, err = c.conn.ServiceRequest(&SAVE, [][]byte{});
	return;

}

// Redis KEYS command.
func (c *synchClient) AllKeys () (result []string, err Error){
	return c.Keys("*");
}

// Redis KEYS command.
func (c *synchClient) Keys (arg0 string) (result []string, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&KEYS, [][]byte{arg0bytes});
	if err == nil {
//		result = strings.Split(bytes.NewBuffer(resp.GetBulkData()).String(), " ", 0);
		result = convAndSplit (resp.GetBulkData());
	}
	return result, err;

}

func convAndSplit (buff []byte) []string {
	return strings.Split(bytes.NewBuffer(buff).String(), " ", 0);
}
/***
// Redis SORT command.
func (c *synchClient) Sort (arg0 string) (result redis.Sort, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SORT, [][]byte{arg0bytes});
	if err == nil {result = resp.GetMultiBulkData();}
	return result, err;

}
***/
// Redis EXISTS command.
func (c *synchClient) Exists (arg0 string) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&EXISTS, [][]byte{arg0bytes});
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

// Redis RENAME command.
func (c *synchClient) Rename (arg0 string, arg1 string) (err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (arg1);

	_, err = c.conn.ServiceRequest(&RENAME, [][]byte{arg0bytes, arg1bytes});
	return;
}

// Redis INFO command.
func (c *synchClient) Info () (result map[string] string, err Error){
	var resp Response;
	resp, err = c.conn.ServiceRequest(&INFO, [][]byte{});
	if err == nil {
		result = parseInfo(resp.GetBulkData());
//		infoStr := bytes.NewBuffer(resp.GetBulkData()).String();
//		infoItems := strings.Split(infoStr, "\r\n", 0);
//		result = make(map[string] string);
//		for _, entry := range infoItems  {
//			etuple := strings.Split(entry, ":", 2);
//			result[etuple[0]] = etuple[1];
//		}
	}
	return result, err;
}

func parseInfo (buff []byte) map[string]string {
	infoStr := bytes.NewBuffer(buff).String();
	infoItems := strings.Split(infoStr, "\r\n", 0);
	result := make(map[string] string);
	for _, entry := range infoItems  {
		etuple := strings.Split(entry, ":", 2);
		result[etuple[0]] = etuple[1];
	}
	return result;
}
// Redis PING command.
func (c *synchClient) Ping () (err Error){
	if c == nil {
		log.Stderr("FAULT in synchclient.Ping(): why is c nil?");
		return NewError(SYSTEM_ERR, "c *synchClient is NIL!");
	}
	else if c.conn == nil {
		log.Stderr("FAULT in synchclient.Ping(): why is c.conn nil?");
		return NewError(SYSTEM_ERR, "c.conn *SynchConnection is NIL!");
	}
	_, err = c.conn.ServiceRequest(&PING, [][]byte{});
	return;
}

// Redis QUIT command.
func (c *synchClient) Quit () (err Error){
	c.conn.Close();
//	_, err = c.conn.ServiceRequest(&QUIT);
	return;
}

// Redis SETNX command.
func (c *synchClient) Setnx (arg0 string, arg1 []byte) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SETNX, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

// Redis GETSET command.
func (c *synchClient) Getset (arg0 string, arg1 []byte) (result []byte, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;

	var resp Response;
	resp, err = c.conn.ServiceRequest(&GETSET, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.GetBulkData();}
	return result, err;

}

// Redis MGET command.
func (c *synchClient) Mget (arg0 string, arg1 []string) (result [][]byte, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := concatAndGetBytes(arg1, " ");

	var resp Response;
	resp, err = c.conn.ServiceRequest(&MGET, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.GetMultiBulkData();}
	return result, err;

}


// Redis INCR command.
func (c *synchClient) Incr (arg0 string) (result int64, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&INCR, [][]byte{arg0bytes});
	if err == nil {result = resp.GetNumberValue();}
	return result, err;

}

// Redis INCRBY command.
func (c *synchClient) Incrby (arg0 string, arg1 int64) (result int64, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&INCRBY, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.GetNumberValue();}
	return result, err;

}

// Redis DECR command.
func (c *synchClient) Decr (arg0 string) (result int64, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&DECR, [][]byte{arg0bytes});
	if err == nil {result = resp.GetNumberValue();}
	return result, err;

}

// Redis DECRBY command.
func (c *synchClient) Decrby (arg0 string, arg1 int64) (result int64, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&DECRBY, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.GetNumberValue();}
	return result, err;

}

// Redis DEL command.
func (c *synchClient) Del (arg0 string) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&DEL, [][]byte{arg0bytes});
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

// Redis RANDOMKEY command.
func (c *synchClient) Randomkey () (result string, err Error){
	var resp Response;
	resp, err = c.conn.ServiceRequest(&RANDOMKEY, [][]byte{});
	if err == nil {result = resp.GetStringValue();}
	return result, err;

}

// Redis RENAMENX command.
func (c *synchClient) Renamenx (arg0 string, arg1 string) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (arg1);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&RENAMENX, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

// Redis DBSIZE command.
func (c *synchClient) Dbsize () (result int64, err Error){
	var resp Response;
	resp, err = c.conn.ServiceRequest(&DBSIZE, [][]byte{});
	if err == nil {result = resp.GetNumberValue();}
	return result, err;

}

// Redis EXPIRE command.
func (c *synchClient) Expire (arg0 string, arg1 int64) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&EXPIRE, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

// Redis TTL command.
func (c *synchClient) Ttl (arg0 string) (result int64, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&TTL, [][]byte{arg0bytes});
	if err == nil {result = resp.GetNumberValue();}
	return result, err;

}

// Redis RPUSH command.
func (c *synchClient) Rpush (arg0 string, arg1 []byte) (err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;

	_, err = c.conn.ServiceRequest(&RPUSH, [][]byte{arg0bytes, arg1bytes});
	return;
}

// Redis LPUSH command.
func (c *synchClient) Lpush (arg0 string, arg1 []byte) (err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;

	_, err = c.conn.ServiceRequest(&LPUSH, [][]byte{arg0bytes, arg1bytes});
	return;
}

// Redis LSET command.
func (c *synchClient) Lset (arg0 string, arg1 int64, arg2 []byte) (err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));
	arg2bytes := arg2;

	_, err = c.conn.ServiceRequest(&LSET, [][]byte{arg0bytes, arg1bytes, arg2bytes});
	return;
}


// Redis LREM command.
func (c *synchClient) Lrem (arg0 string, arg1 []byte, arg2 int64) (result int64, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;
	arg2bytes := strings.Bytes (fmt.Sprintf("%d", arg2));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&LREM, [][]byte{arg0bytes, arg1bytes, arg2bytes});
	if err == nil {result = resp.GetNumberValue();}
	return result, err;

}

// Redis LLEN command.
func (c *synchClient) Llen (arg0 string) (result int64, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&LLEN, [][]byte{arg0bytes});
	if err == nil {result = resp.GetNumberValue();}
	return result, err;

}

// Redis LRANGE command.
func (c *synchClient) Lrange (arg0 string, arg1 int64, arg2 int64) (result [][]byte, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));
	arg2bytes := strings.Bytes (fmt.Sprintf("%d", arg2));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&LRANGE, [][]byte{arg0bytes, arg1bytes, arg2bytes});
	if err == nil {result = resp.GetMultiBulkData();}
	return result, err;

}

// Redis LTRIM command.
func (c *synchClient) Ltrim (arg0 string, arg1 int64, arg2 int64) (err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));
	arg2bytes := strings.Bytes (fmt.Sprintf("%d", arg2));

	_, err = c.conn.ServiceRequest(&LTRIM, [][]byte{arg0bytes, arg1bytes, arg2bytes});
	return;
}

// Redis LINDEX command.
func (c *synchClient) Lindex (arg0 string, arg1 int64) (result []byte, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&LINDEX, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.GetBulkData();}
	return result, err;

}

// Redis LPOP command.
func (c *synchClient) Lpop (arg0 string) (result []byte, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&LPOP, [][]byte{arg0bytes});
	if err == nil {result = resp.GetBulkData();}
	return result, err;

}

// Redis RPOP command.
func (c *synchClient) Rpop (arg0 string) (result []byte, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&RPOP, [][]byte{arg0bytes});
	if err == nil {result = resp.GetBulkData();}
	return result, err;

}

// Redis RPOPLPUSH command.
func (c *synchClient) Rpoplpush (arg0 string, arg1 string) (result []byte, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (arg1);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&RPOPLPUSH, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.GetBulkData();}
	return result, err;

}

// Redis SADD command.
func (c *synchClient) Sadd (arg0 string, arg1 []byte) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SADD, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}


// Redis SREM command.
func (c *synchClient) Srem (arg0 string, arg1 []byte) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SREM, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}


// Redis SISMEMBER command.
func (c *synchClient) Sismember (arg0 string, arg1 []byte) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SISMEMBER, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

// Redis SMOVE command.
func (c *synchClient) Smove (arg0 string, arg1 string, arg2 []byte) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (arg1);
	arg2bytes := arg2;

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SMOVE, [][]byte{arg0bytes, arg1bytes, arg2bytes});
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

// Redis SCARD command.
func (c *synchClient) Scard (arg0 string) (result int64, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SCARD, [][]byte{arg0bytes});
	if err == nil {result = resp.GetNumberValue();}
	return result, err;

}

func concatAndGetBytes(arr []string, delim string) []byte {
	cstr := "";
	for _, s := range arr {
		cstr += s;
		cstr += delim;
	}
	return strings.Bytes(cstr);
}
// Redis SINTER command.
func (c *synchClient) Sinter (arg0 string, arg1 []string) (result [][]byte, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := concatAndGetBytes(arg1, " ");

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SINTER, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.GetMultiBulkData();}
	return result, err;

}

// Redis SINTERSTORE command.
func (c *synchClient) Sinterstore (arg0 string, arg1 []string) (err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := concatAndGetBytes(arg1, " ");

	_, err = c.conn.ServiceRequest(&SINTERSTORE, [][]byte{arg0bytes, arg1bytes});
	return;
}

// Redis SUNION command.
func (c *synchClient) Sunion (arg0 string, arg1 []string) (result [][]byte, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := concatAndGetBytes(arg1, " ");

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SUNION, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.GetMultiBulkData();}
	return result, err;

}

// Redis SUNIONSTORE command.
func (c *synchClient) Sunionstore (arg0 string, arg1 []string) (err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := concatAndGetBytes(arg1, " ");

	_, err = c.conn.ServiceRequest(&SUNIONSTORE, [][]byte{arg0bytes, arg1bytes});
	return;
}

// Redis SDIFF command.
func (c *synchClient) Sdiff (arg0 string, arg1 []string) (result [][]byte, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := concatAndGetBytes(arg1, " ");

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SDIFF, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.GetMultiBulkData();}
	return result, err;

}

// Redis SDIFFSTORE command.
func (c *synchClient) Sdiffstore (arg0 string, arg1 []string) (err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := concatAndGetBytes(arg1, " ");

	_, err = c.conn.ServiceRequest(&SDIFFSTORE, [][]byte{arg0bytes, arg1bytes});
	return;
}

// Redis SMEMBERS command.
func (c *synchClient) Smembers (arg0 string) (result [][]byte, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SMEMBERS, [][]byte{arg0bytes});
	if err == nil {result = resp.GetMultiBulkData();}
	return result, err;

}

// Redis SRANDMEMBER command.
func (c *synchClient) Srandmember (arg0 string) (result []byte, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SRANDMEMBER, [][]byte{arg0bytes});
	if err == nil {result = resp.GetBulkData();}
	return result, err;

}

// Redis ZADD command.
func (c *synchClient) Zadd (arg0 string, arg1 float64, arg2 []byte) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%e", arg1));
	arg2bytes := arg2;

	var resp Response;
	resp, err = c.conn.ServiceRequest(&ZADD, [][]byte{arg0bytes, arg1bytes, arg2bytes});
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

// Redis ZREM command.
func (c *synchClient) Zrem (arg0 string, arg1 []byte) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;

	var resp Response;
	resp, err = c.conn.ServiceRequest(&ZREM, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

// Redis ZCARD command.
func (c *synchClient) Zcard (arg0 string) (result int64, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&ZCARD, [][]byte{arg0bytes});
	if err == nil {result = resp.GetNumberValue();}
	return result, err;

}

// Redis ZSCORE command.
func (c *synchClient) Zscore (arg0 string, arg1 []byte) (result float64, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;

	var resp Response;
	resp, err = c.conn.ServiceRequest(&ZSCORE, [][]byte{arg0bytes, arg1bytes});
	if err == nil {
		buff := resp.GetBulkData();
//		fnum, oserr := strconv.Atof64(bytes.NewBuffer(buff).String());
//		if oserr != nil {
//			err = NewErrorWithCause(SYSTEM_ERR, "Expected a parsable byte representation of a float64 in Zscore!", oserr);
//		}
//		result = fnum;
		result, err = Btof64 (buff);
	}
	return result, err;

}

func Btof64 (buff []byte) (num float64, e Error){
	num, ce := strconv.Atof64(bytes.NewBuffer(buff).String());
	if ce != nil {
		e = NewErrorWithCause(SYSTEM_ERR, "Expected a parsable byte representation of a float64", ce);
	}
	return;
}

// Redis ZRANGE command.
func (c *synchClient) Zrange (arg0 string, arg1 int64, arg2 int64) (result [][]byte, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));
	arg2bytes := strings.Bytes (fmt.Sprintf("%d", arg2));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&ZRANGE, [][]byte{arg0bytes, arg1bytes, arg2bytes});
	if err == nil {result = resp.GetMultiBulkData();}
	return result, err;

}

// Redis ZREVRANGE command.
func (c *synchClient) Zrevrange (arg0 string, arg1 int64, arg2 int64) (result [][]byte, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));
	arg2bytes := strings.Bytes (fmt.Sprintf("%d", arg2));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&ZREVRANGE, [][]byte{arg0bytes, arg1bytes, arg2bytes});
	if err == nil {result = resp.GetMultiBulkData();}
	return result, err;

}

// Redis ZRANGEBYSCORE command.
func (c *synchClient) Zrangebyscore (arg0 string, arg1 float64, arg2 float64) (result [][]byte, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%e", arg1));
	arg2bytes := strings.Bytes (fmt.Sprintf("%e", arg2));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&ZRANGEBYSCORE, [][]byte{arg0bytes, arg1bytes, arg2bytes});
	if err == nil {result = resp.GetMultiBulkData();}
	return result, err;

}

// Redis FLUSHDB command.
func (c *synchClient) Flushdb () (err Error){
	_, err = c.conn.ServiceRequest(&FLUSHDB, [][]byte{});
	return;
}

// Redis FLUSHALL command.
func (c *synchClient) Flushall () (err Error){
	_, err = c.conn.ServiceRequest(&FLUSHALL, [][]byte{});
	return;
}

// Redis MOVE command.
func (c *synchClient) Move (arg0 string, arg1 int64) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&MOVE, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

// Redis BGSAVE command.
func (c *synchClient) Bgsave () (err Error){
	_, err = c.conn.ServiceRequest(&BGSAVE, [][]byte{});
	return;
}

// Redis LASTSAVE command.
func (c *synchClient) Lastsave () (result int64, err Error){
	var resp Response;
	resp, err = c.conn.ServiceRequest(&LASTSAVE, [][]byte{});
	if err == nil {result = resp.GetNumberValue();}
	return result, err;

}
