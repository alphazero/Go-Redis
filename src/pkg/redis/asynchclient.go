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
/*
	"bytes";
	"fmt";
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

// ----------------------- aync interface


func (c *async) Incr (arg0 string) (result FutureInt64, err Error) {
	arg0bytes := strings.Bytes (arg0);

	resp, err := c.conn.QueueRequest(&INCR, [][]byte{arg0bytes});
	if err == nil {result = resp.future.(FutureInt64);}
	return result, err;
}
func (c *async) Get (arg0 string) (result FutureBytes, err Error) {
	arg0bytes := strings.Bytes (arg0);

	resp, err := c.conn.QueueRequest(&GET, [][]byte{arg0bytes});
	if err == nil {result = resp.future.(FutureBytes);}
	return result, err;
}
func (c *async) Set (arg0 string, arg1 []byte) (result FutureString, err Error) {
//	log.Stdout("asynchClient.Set: about to issue request SET for: ", arg0);
	arg0bytes := strings.Bytes (arg0);
	arg1bytes := arg1;

	resp, err:= c.conn.QueueRequest(&SET, [][]byte{arg0bytes, arg1bytes});
	if err == nil {result = resp.future.(FutureString);}
	return;
}
func (c *async) Exists (arg0 string) (result FutureBool, err Error) {
	arg0bytes := strings.Bytes (arg0);

	resp, err := c.conn.QueueRequest(&EXISTS, [][]byte{arg0bytes});
	if err == nil {result = resp.future.(FutureBool);}
	return result, err;

}
/*
func (c *async) Sismembers (key string) (FutureBytesArray, Error) {}
func (c *async) Incr (key string) (FutureInt64, Error) {}
func (c *async) Randomkey () (FutureString, Error) {}
*/