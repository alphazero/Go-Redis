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
package client

import (
	"os";
//	"log";
	"strings";
	"bytes";
	"fmt";
	"strconv";
	"connection";
	. "protocol";
	. "redis";
)

type Client struct {
	conn  connection.Endpoint;
}

func Connect () (c *Client, err os.Error){
	spec := connection.DefaultSpec();
	c, err = ConnectWithSpec(spec);
	return;
}

func ConnectWithSpec (spec *connection.Spec) (c *Client, err os.Error) {
	c = new(Client);
	c.conn, err = connection.OpenNew (spec);
	return;
}

/** ----------------- REDIS INTERFACE ------------- **/
/** ----------------- G E N E R A T E ------------- **/




// Redis GET command.
func (c *Client) Get (arg0 string) (result []byte, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&GET, arg0bytes);
	if err == nil {result = resp.GetBulkData();}
	return result, err;

}

// Redis TYPE command.
func (c *Client) Type (arg0 string) (result KeyType, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&TYPE, arg0bytes);
	if err == nil {result = GetKeyType(resp.GetStringValue());}
	return result, err;

}

// Redis SET command.
func (c *Client) Set (arg0 string, arg1 []byte) (err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;

	_, err = c.conn.ServiceRequest(&SET, arg0bytes, arg1bytes);
	return;
}
/*****
// Redis SET command.
func (c *Client) Set (arg0 string, arg1 string) (err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (arg1);

	_, err = c.conn.ServiceRequest(&SET, arg0bytes, arg1bytes);
	return result, err;
}

// Redis SET command.
func (c *Client) Set (arg0 string, arg1 int64) (err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));

	_, err = c.conn.ServiceRequest(&SET, arg0bytes, arg1bytes);
	return result, err;
}

// Redis SET command.
func (c *Client) Set (arg0 string, arg1 io.ReadWriter) (err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := convertSerializableToBytes????????(arg1);

	_, err = c.conn.ServiceRequest(&SET, arg0bytes, arg1bytes);
	return result, err;
}
******/

// Redis SAVE command.
func (c *Client) Save () (err Error){
	_, err = c.conn.ServiceRequest(&SAVE);
	return;

}

// Redis KEYS command.
func (c *Client) AllKeys () (result []string, err Error){
	return c.Keys("*");
}

// Redis KEYS command.
func (c *Client) Keys (arg0 string) (result []string, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&KEYS, arg0bytes);
	if err == nil {
		result = strings.Split(bytes.NewBuffer(resp.GetBulkData()).String(), " ", 0);
	}
	return result, err;

}
/***
// Redis SORT command.
func (c *Client) Sort (arg0 string) (result redis.Sort, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SORT, arg0bytes);
	if err == nil {result = resp.GetMultiBulkData();}
	return result, err;

}
***/
// Redis EXISTS command.
func (c *Client) Exists (arg0 string) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&EXISTS, arg0bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

// Redis RENAME command.
func (c *Client) Rename (arg0 string, arg1 string) (err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (arg1);

	_, err = c.conn.ServiceRequest(&RENAME, arg0bytes, arg1bytes);
	return;
}

// Redis INFO command.
func (c *Client) Info () (result map[string] string, err Error){
	var resp Response;
	resp, err = c.conn.ServiceRequest(&INFO);
	if err == nil {
		infoStr := bytes.NewBuffer(resp.GetBulkData()).String();
		infoItems := strings.Split(infoStr, "\r\n", 0);
		result = make(map[string] string);
		for _, entry := range infoItems  {
			etuple := strings.Split(entry, ":", 2);
			result[etuple[0]] = etuple[1];
		}
	}
	return result, err;
}

// Redis PING command.
func (c *Client) Ping () (err Error){
	_, err = c.conn.ServiceRequest(&PING);
	return;
}

// Redis QUIT command.
func (c *Client) Quit () (err Error){
	_, err = c.conn.ServiceRequest(&QUIT);
	return;
}

// Redis SETNX command.
func (c *Client) Setnx (arg0 string, arg1 []byte) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SETNX, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}
/**
// Redis SETNX command.
func (c *Client) Setnx (arg0 string, arg1 io.ReadWriter) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := convertSerializableToBytes????????(arg1);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SETNX, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

// Redis SETNX command.
func (c *Client) Setnx (arg0 string, arg1 string) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (arg1);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SETNX, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

// Redis SETNX command.
func (c *Client) Setnx (arg0 string, arg1 int64) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SETNX, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}
*/

// Redis GETSET command.
func (c *Client) Getset (arg0 string, arg1 []byte) (result []byte, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;

	var resp Response;
	resp, err = c.conn.ServiceRequest(&GETSET, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBulkData();}
	return result, err;

}

/**
// Redis GETSET command.
func (c *Client) Getset (arg0 string, arg1 string) (result []byte, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (arg1);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&GETSET, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBulkData();}
	return result, err;

}

// Redis GETSET command.
func (c *Client) Getset (arg0 string, arg1 int64) (result []byte, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&GETSET, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBulkData();}
	return result, err;

}

// Redis GETSET command.
func (c *Client) Getset (arg0 string, arg1 io.ReadWriter) (result []byte, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := convertSerializableToBytes????????(arg1);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&GETSET, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBulkData();}
	return result, err;

}
**/

// Redis MGET command.
func (c *Client) Mget (arg0 string, arg1 []string) (result [][]byte, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := concatAndGetBytes(arg1, " ");

	var resp Response;
	resp, err = c.conn.ServiceRequest(&MGET, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetMultiBulkData();}
	return result, err;

}


// Redis INCR command.
func (c *Client) Incr (arg0 string) (result int64, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&INCR, arg0bytes);
	if err == nil {result = resp.GetNumberValue();}
	return result, err;

}

// Redis INCRBY command.
func (c *Client) Incrby (arg0 string, arg1 int64) (result int64, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&INCRBY, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetNumberValue();}
	return result, err;

}

// Redis DECR command.
func (c *Client) Decr (arg0 string) (result int64, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&DECR, arg0bytes);
	if err == nil {result = resp.GetNumberValue();}
	return result, err;

}

// Redis DECRBY command.
func (c *Client) Decrby (arg0 string, arg1 int64) (result int64, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&DECRBY, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetNumberValue();}
	return result, err;

}

// Redis DEL command.
func (c *Client) Del (arg0 string) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&DEL, arg0bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

// Redis RANDOMKEY command.
func (c *Client) Randomkey () (result string, err Error){
	var resp Response;
	resp, err = c.conn.ServiceRequest(&RANDOMKEY);
	if err == nil {result = resp.GetStringValue();}
	return result, err;

}

// Redis RENAMENX command.
func (c *Client) Renamenx (arg0 string, arg1 string) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (arg1);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&RENAMENX, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

// Redis DBSIZE command.
func (c *Client) Dbsize () (result int64, err Error){
	var resp Response;
	resp, err = c.conn.ServiceRequest(&DBSIZE);
	if err == nil {result = resp.GetNumberValue();}
	return result, err;

}

// Redis EXPIRE command.
func (c *Client) Expire (arg0 string, arg1 int64) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&EXPIRE, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

// Redis TTL command.
func (c *Client) Ttl (arg0 string) (result int64, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&TTL, arg0bytes);
	if err == nil {result = resp.GetNumberValue();}
	return result, err;

}

// Redis RPUSH command.
func (c *Client) Rpush (arg0 string, arg1 []byte) (err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;

	_, err = c.conn.ServiceRequest(&RPUSH, arg0bytes, arg1bytes);
	return;
}

/**
// Redis RPUSH command.
func (c *Client) Rpush (arg0 string, arg1 io.ReadWriter) (err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := convertSerializableToBytes????????(arg1);

	_, err = c.conn.ServiceRequest(&RPUSH, arg0bytes, arg1bytes);
	return result, err;

}

// Redis RPUSH command.
func (c *Client) Rpush (arg0 string, arg1 string) (err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (arg1);

	_, err = c.conn.ServiceRequest(&RPUSH, arg0bytes, arg1bytes);
	return result, err;

}

// Redis RPUSH command.
func (c *Client) Rpush (arg0 string, arg1 int64) (err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));

	_, err = c.conn.ServiceRequest(&RPUSH, arg0bytes, arg1bytes);
	return result, err;

}
**/
// Redis LPUSH command.
func (c *Client) Lpush (arg0 string, arg1 []byte) (err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;

	_, err = c.conn.ServiceRequest(&LPUSH, arg0bytes, arg1bytes);
	return;
}


/**
// Redis LPUSH command.
func (c *Client) Lpush (arg0 string, arg1 io.ReadWriter) (err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := convertSerializableToBytes????????(arg1);

	_, err = c.conn.ServiceRequest(&LPUSH, arg0bytes, arg1bytes);
	return;
}

// Redis LPUSH command.
func (c *Client) Lpush (arg0 string, arg1 int64) (err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));

	_, err = c.conn.ServiceRequest(&LPUSH, arg0bytes, arg1bytes);
	return;
}

// Redis LPUSH command.
func (c *Client) Lpush (arg0 string, arg1 string) (err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (arg1);

	_, err = c.conn.ServiceRequest(&LPUSH, arg0bytes, arg1bytes);
	return;
}
**/

// Redis LSET command.
func (c *Client) Lset (arg0 string, arg1 int64, arg2 []byte) (err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));
	arg2bytes := arg2;

	_, err = c.conn.ServiceRequest(&LSET, arg0bytes, arg1bytes, arg2bytes);
	return;
}

/**
// Redis LSET command.
func (c *Client) Lset (arg0 string, arg1 int64, arg2 io.ReadWriter) (err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));
	arg2bytes := convertSerializableToBytes????????(arg2);

	_, err = c.conn.ServiceRequest(&LSET, arg0bytes, arg1bytes, arg2bytes);
	return;
}

// Redis LSET command.
func (c *Client) Lset (arg0 string, arg1 int64, arg2 string) (err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));
	arg2bytes := strings.Bytes (arg2);

	_, err = c.conn.ServiceRequest(&LSET, arg0bytes, arg1bytes, arg2bytes);
	return;
}

// Redis LSET command.
func (c *Client) Lset (arg0 string, arg1 int64, arg2 int64) (err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));
	arg2bytes := strings.Bytes (fmt.Sprintf("%d", arg2));

	_, err = c.conn.ServiceRequest(&LSET, arg0bytes, arg1bytes, arg2bytes);
	return;
}
**/

// Redis LREM command.
func (c *Client) Lrem (arg0 string, arg1 []byte, arg2 int64) (result int64, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;
	arg2bytes := strings.Bytes (fmt.Sprintf("%d", arg2));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&LREM, arg0bytes, arg1bytes, arg2bytes);
	if err == nil {result = resp.GetNumberValue();}
	return result, err;

}

/**
// Redis LREM command.
func (c *Client) Lrem (arg0 string, arg1 io.ReadWriter, arg2 int64) (result int64, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := convertSerializableToBytes????????(arg1);
	arg2bytes := strings.Bytes (fmt.Sprintf("%d", arg2));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&LREM, arg0bytes, arg1bytes, arg2bytes);
	if err == nil {result = resp.GetNumberValue();}
	return result, err;

}

// Redis LREM command.
func (c *Client) Lrem (arg0 string, arg1 string, arg2 int64) (result int64, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (arg1);
	arg2bytes := strings.Bytes (fmt.Sprintf("%d", arg2));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&LREM, arg0bytes, arg1bytes, arg2bytes);
	if err == nil {result = resp.GetNumberValue();}
	return result, err;

}

// Redis LREM command.
func (c *Client) Lrem (arg0 string, arg1 int64, arg2 int64) (result int64, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));
	arg2bytes := strings.Bytes (fmt.Sprintf("%d", arg2));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&LREM, arg0bytes, arg1bytes, arg2bytes);
	if err == nil {result = resp.GetNumberValue();}
	return result, err;

}
*/

// Redis LLEN command.
func (c *Client) Llen (arg0 string) (result int64, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&LLEN, arg0bytes);
	if err == nil {result = resp.GetNumberValue();}
	return result, err;

}

// Redis LRANGE command.
func (c *Client) Lrange (arg0 string, arg1 int64, arg2 int64) (result [][]byte, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));
	arg2bytes := strings.Bytes (fmt.Sprintf("%d", arg2));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&LRANGE, arg0bytes, arg1bytes, arg2bytes);
	if err == nil {result = resp.GetMultiBulkData();}
	return result, err;

}

// Redis LTRIM command.
func (c *Client) Ltrim (arg0 string, arg1 int64, arg2 int64) (err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));
	arg2bytes := strings.Bytes (fmt.Sprintf("%d", arg2));

	_, err = c.conn.ServiceRequest(&LTRIM, arg0bytes, arg1bytes, arg2bytes);
	return;
}

// Redis LINDEX command.
func (c *Client) Lindex (arg0 string, arg1 int64) (result []byte, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&LINDEX, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBulkData();}
	return result, err;

}

// Redis LPOP command.
func (c *Client) Lpop (arg0 string) (result []byte, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&LPOP, arg0bytes);
	if err == nil {result = resp.GetBulkData();}
	return result, err;

}

// Redis RPOP command.
func (c *Client) Rpop (arg0 string) (result []byte, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&RPOP, arg0bytes);
	if err == nil {result = resp.GetBulkData();}
	return result, err;

}

// Redis RPOPLPUSH command.
func (c *Client) Rpoplpush (arg0 string, arg1 string) (result []byte, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (arg1);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&RPOPLPUSH, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBulkData();}
	return result, err;

}

// Redis SADD command.
func (c *Client) Sadd (arg0 string, arg1 []byte) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SADD, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

/**
// Redis SADD command.
func (c *Client) Sadd (arg0 string, arg1 io.ReadWriter) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := convertSerializableToBytes????????(arg1);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SADD, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

// Redis SADD command.
func (c *Client) Sadd (arg0 string, arg1 int64) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SADD, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

// Redis SADD command.
func (c *Client) Sadd (arg0 string, arg1 string) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (arg1);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SADD, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}
*/

// Redis SREM command.
func (c *Client) Srem (arg0 string, arg1 []byte) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SREM, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

/**
// Redis SREM command.
func (c *Client) Srem (arg0 string, arg1 string) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (arg1);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SREM, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

// Redis SREM command.
func (c *Client) Srem (arg0 string, arg1 int64) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SREM, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

// Redis SREM command.
func (c *Client) Srem (arg0 string, arg1 io.ReadWriter) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := convertSerializableToBytes????????(arg1);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SREM, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}
*/

// Redis SISMEMBER command.
func (c *Client) Sismember (arg0 string, arg1 []byte) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SISMEMBER, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

/**
// Redis SISMEMBER command.
func (c *Client) Sismember (arg0 string, arg1 io.ReadWriter) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := convertSerializableToBytes????????(arg1);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SISMEMBER, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

// Redis SISMEMBER command.
func (c *Client) Sismember (arg0 string, arg1 int64) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SISMEMBER, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

// Redis SISMEMBER command.
func (c *Client) Sismember (arg0 string, arg1 string) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (arg1);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SISMEMBER, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}
*/

// Redis SMOVE command.
func (c *Client) Smove (arg0 string, arg1 string, arg2 []byte) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (arg1);
	arg2bytes := arg2;

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SMOVE, arg0bytes, arg1bytes, arg2bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

/**
// Redis SMOVE command.
func (c *Client) Smove (arg0 string, arg1 string, arg2 io.ReadWriter) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (arg1);
	arg2bytes := convertSerializableToBytes????????(arg2);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SMOVE, arg0bytes, arg1bytes, arg2bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

// Redis SMOVE command.
func (c *Client) Smove (arg0 string, arg1 string, arg2 string) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (arg1);
	arg2bytes := strings.Bytes (arg2);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SMOVE, arg0bytes, arg1bytes, arg2bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

// Redis SMOVE command.
func (c *Client) Smove (arg0 string, arg1 string, arg2 int64) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (arg1);
	arg2bytes := strings.Bytes (fmt.Sprintf("%d", arg2));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SMOVE, arg0bytes, arg1bytes, arg2bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}
*/

// Redis SCARD command.
func (c *Client) Scard (arg0 string) (result int64, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SCARD, arg0bytes);
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
func (c *Client) Sinter (arg0 string, arg1 []string) (result [][]byte, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := concatAndGetBytes(arg1, " ");

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SINTER, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetMultiBulkData();}
	return result, err;

}

// Redis SINTERSTORE command.
func (c *Client) Sinterstore (arg0 string, arg1 []string) (err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := concatAndGetBytes(arg1, " ");

	_, err = c.conn.ServiceRequest(&SINTERSTORE, arg0bytes, arg1bytes);
	return;
}

// Redis SUNION command.
func (c *Client) Sunion (arg0 string, arg1 []string) (result [][]byte, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := concatAndGetBytes(arg1, " ");

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SUNION, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetMultiBulkData();}
	return result, err;

}

// Redis SUNIONSTORE command.
func (c *Client) Sunionstore (arg0 string, arg1 []string) (err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := concatAndGetBytes(arg1, " ");

	_, err = c.conn.ServiceRequest(&SUNIONSTORE, arg0bytes, arg1bytes);
	return;
}

// Redis SDIFF command.
func (c *Client) Sdiff (arg0 string, arg1 []string) (result [][]byte, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := concatAndGetBytes(arg1, " ");

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SDIFF, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetMultiBulkData();}
	return result, err;

}

// Redis SDIFFSTORE command.
func (c *Client) Sdiffstore (arg0 string, arg1 []string) (err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := concatAndGetBytes(arg1, " ");

	_, err = c.conn.ServiceRequest(&SDIFFSTORE, arg0bytes, arg1bytes);
	return;
}

// Redis SMEMBERS command.
func (c *Client) Smembers (arg0 string) (result [][]byte, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SMEMBERS, arg0bytes);
	if err == nil {result = resp.GetMultiBulkData();}
	return result, err;

}

// Redis SRANDMEMBER command.
func (c *Client) Srandmember (arg0 string) (result []byte, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&SRANDMEMBER, arg0bytes);
	if err == nil {result = resp.GetBulkData();}
	return result, err;

}

// Redis ZADD command.
func (c *Client) Zadd (arg0 string, arg1 float64, arg2 []byte) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%e", arg1));
	arg2bytes := arg2;

	var resp Response;
	resp, err = c.conn.ServiceRequest(&ZADD, arg0bytes, arg1bytes, arg2bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

/**
// Redis ZADD command.
func (c *Client) Zadd (arg0 string, arg1 float64, arg2 string) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%e", arg1));
	arg2bytes := strings.Bytes (arg2);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&ZADD, arg0bytes, arg1bytes, arg2bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

// Redis ZADD command.
func (c *Client) Zadd (arg0 string, arg1 float64, arg2 int64) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%e", arg1));
	arg2bytes := strings.Bytes (fmt.Sprintf("%d", arg2));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&ZADD, arg0bytes, arg1bytes, arg2bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

// Redis ZADD command.
func (c *Client) Zadd (arg0 string, arg1 float64, arg2 io.ReadWriter) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%e", arg1));
	arg2bytes := convertSerializableToBytes????????(arg2);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&ZADD, arg0bytes, arg1bytes, arg2bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}
*/

// Redis ZREM command.
func (c *Client) Zrem (arg0 string, arg1 []byte) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;

	var resp Response;
	resp, err = c.conn.ServiceRequest(&ZREM, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

/**
// Redis ZREM command.
func (c *Client) Zrem (arg0 string, arg1 io.ReadWriter) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := convertSerializableToBytes????????(arg1);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&ZREM, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

// Redis ZREM command.
func (c *Client) Zrem (arg0 string, arg1 int64) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&ZREM, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

// Redis ZREM command.
func (c *Client) Zrem (arg0 string, arg1 string) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (arg1);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&ZREM, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}
*/

// Redis ZCARD command.
func (c *Client) Zcard (arg0 string) (result int64, err Error){
	arg0bytes := strings.Bytes (arg0);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&ZCARD, arg0bytes);
	if err == nil {result = resp.GetNumberValue();}
	return result, err;

}

// Redis ZSCORE command.
func (c *Client) Zscore (arg0 string, arg1 []byte) (result float64, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;

	var resp Response;
	resp, err = c.conn.ServiceRequest(&ZSCORE, arg0bytes, arg1bytes);
	if err == nil {
		buff := resp.GetBulkData();
		fnum, oserr := strconv.Atof64(bytes.NewBuffer(buff).String());
		if oserr != nil {
			err = NewErrorWithCause(SYSTEM_ERR, "Expected a parsable byte representation of a float64 in Zscore!", oserr);
		}
		result = fnum;
	}
	return result, err;

}

/**
// Redis ZSCORE command.
func (c *Client) Zscore (arg0 string, arg1 io.ReadWriter) (result float64, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := convertSerializableToBytes????????(arg1);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&ZSCORE, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBulkData();}
	return result, err;

}

// Redis ZSCORE command.
func (c *Client) Zscore (arg0 string, arg1 int64) (result float64, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&ZSCORE, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBulkData();}
	return result, err;

}

// Redis ZSCORE command.
func (c *Client) Zscore (arg0 string, arg1 string) (result float64, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (arg1);

	var resp Response;
	resp, err = c.conn.ServiceRequest(&ZSCORE, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBulkData();}
	return result, err;

}
*/

// Redis ZRANGE command.
func (c *Client) Zrange (arg0 string, arg1 int64, arg2 int64) (result [][]byte, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));
	arg2bytes := strings.Bytes (fmt.Sprintf("%d", arg2));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&ZRANGE, arg0bytes, arg1bytes, arg2bytes);
	if err == nil {result = resp.GetMultiBulkData();}
	return result, err;

}

// Redis ZREVRANGE command.
func (c *Client) Zrevrange (arg0 string, arg1 int64, arg2 int64) (result [][]byte, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));
	arg2bytes := strings.Bytes (fmt.Sprintf("%d", arg2));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&ZREVRANGE, arg0bytes, arg1bytes, arg2bytes);
	if err == nil {result = resp.GetMultiBulkData();}
	return result, err;

}

// Redis ZRANGEBYSCORE command.
func (c *Client) Zrangebyscore (arg0 string, arg1 float64, arg2 float64) (result [][]byte, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%e", arg1));
	arg2bytes := strings.Bytes (fmt.Sprintf("%e", arg2));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&ZRANGEBYSCORE, arg0bytes, arg1bytes, arg2bytes);
	if err == nil {result = resp.GetMultiBulkData();}
	return result, err;

}

// Redis FLUSHDB command.
func (c *Client) Flushdb () (err Error){
	_, err = c.conn.ServiceRequest(&FLUSHDB);
	return;
}

// Redis FLUSHALL command.
func (c *Client) Flushall () (err Error){
	_, err = c.conn.ServiceRequest(&FLUSHALL);
	return;
}

// Redis MOVE command.
func (c *Client) Move (arg0 string, arg1 int64) (result bool, err Error){
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := strings.Bytes (fmt.Sprintf("%d", arg1));

	var resp Response;
	resp, err = c.conn.ServiceRequest(&MOVE, arg0bytes, arg1bytes);
	if err == nil {result = resp.GetBooleanValue();}
	return result, err;

}

// Redis BGSAVE command.
func (c *Client) Bgsave () (err Error){
	_, err = c.conn.ServiceRequest(&BGSAVE);
	return;
}

// Redis LASTSAVE command.
func (c *Client) Lastsave () (result int64, err Error){
	var resp Response;
	resp, err = c.conn.ServiceRequest(&LASTSAVE);
	if err == nil {result = resp.GetNumberValue();}
	return result, err;

}











/**
func (c *Client) Set (key string, arg1 []byte) (Error) {
	keybytes := strings.Bytes(key);

	_, err := c.conn.ServiceRequest (&SET, keybytes, arg1);
	return err;
}

func (c *Client) Exists (key string) (result bool, err Error){
	keybytes := strings.Bytes(key);

	var resp Response;
	resp, err = c.conn.ServiceRequest (&EXISTS, keybytes);
	if err == nil { result = resp.GetBooleanValue(); }
	return result, err;
}

func (c *Client) Type (key string) (result string, err Error){
	keybytes := strings.Bytes(key);

	var resp Response;
	resp, err = c.conn.ServiceRequest (&TYPE, keybytes);
	if err == nil { result = resp.GetStringValue(); }
	return result, err;
}

func (c *Client) Incr (key string) (num int64, err Error){
	keybytes := strings.Bytes(key);

	var resp Response;
	resp, err = c.conn.ServiceRequest (&INCR, keybytes);
	if err == nil { num = resp.GetNumberValue(); }
	return num, err;
}

func (c *Client) Ping() (err Error){
	_, err := c.conn.ServiceRequest (&PING);
	return err;
}

func (c *Client) Get(key string) (data []byte, e Error) {
	log.Stdout("=> Get " + key);
	keybytes := strings.Bytes(key);
	
	var resp Response;
	resp, e = c.conn.ServiceRequest (&GET, keybytes);
	if e == nil { data = resp.GetBulkData(); }
	
	return data, e;
}

func (c *Client) Smembers (key string) (data [][]byte, e Error) {
	log.Stdout("=> SMEMBERS " + key);
	keybytes := strings.Bytes(key);
	var resp Response;
	resp, e = c.conn.ServiceRequest (&SMEMBERS, keybytes);
	if e == nil {
		data = resp.GetMultiBulkData();
	}
	return data, e;
}
**/
