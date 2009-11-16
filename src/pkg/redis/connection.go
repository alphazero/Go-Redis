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
*/
package redis

import (
	"net";
	"fmt";
	"os";
	"io";
	"bufio";
	"log";
)

const (
	TCP = "tcp";
	LOCALHOST = "127.0.0.1";
)

// ----------------------------------------------------------------------------
// Connection ConnectionSpec
// ----------------------------------------------------------------------------

// Defines the set of parameters that are used by the client connections
//
type ConnectionSpec struct {
	host string;
	port int;
	password string;
	db 	int;
}

// Creates a ConnectionSpec using default settings.
// host is localhost
// port is 6379
// no password is specified (so no AUTH on connect)
// no db is specified (so no SELECT on connect)
//
func DefaultSpec () *ConnectionSpec {
	spec := new (ConnectionSpec);
	spec.host = LOCALHOST;
	spec.port = 6379;
	return spec;
}

// Sets the db for connection spec and returns the reference
// Note that you should not this after you have already connected.
//
func (spec *ConnectionSpec) Db(db int) *ConnectionSpec {
	spec.db = db;
	return spec;
}

// Sets the host for connection spec and returns the reference
// Note that you should not this after you have already connected.
//
func (spec *ConnectionSpec) Host(host string) *ConnectionSpec {
	spec.host = host;
	return spec;
}

// Sets the port for connection spec and returns the reference
// Note that you should not this after you have already connected.
//
func (spec *ConnectionSpec) Port(port int) *ConnectionSpec {
	spec.port = port;
	return spec;
}

// Sets the password for connection spec and returns the reference
// Note that you should not this after you have already connected.
//
func (spec *ConnectionSpec) Password(password string) *ConnectionSpec {
	spec.password = password;
	return spec;
}

func (spec *ConnectionSpec) Addr () string {
	return fmt.Sprintf("%s:%d", spec.host, spec.port);
}

// ----------------------------------------------------------------------------
// Connection SyncConnection
// ----------------------------------------------------------------------------

// Defines the service contract supported by synchronous (Request/Reply)
// connections.

type SyncConnection interface {
	ServiceRequest (cmd *Command, args ...) (Response, Error);
	Close () os.Error;
}

// General control structure used by connections.
//
type _connection struct {
	spec 	*ConnectionSpec;
	conn 	net.Conn;
	reader 	*bufio.Reader;
}
func (hdl _connection) Close() os.Error {
	err := hdl.conn.Close();
//	log.Stdout ("Closed connection");
	return err;
}

// Implementation of SyncConnection.ServiceRequest.
//
func (c _connection) ServiceRequest (cmd *Command, args ...) (resp Response, err Error) {
	
	// TODO: need to consider errors here -- assuming it is always ok to write to Buffer ...
	buff, e1 := CreateRequestBytes(cmd, args);
	if e1 != nil {
		return nil, NewErrorWithCause(SYSTEM_ERR, "ServiceRequest(): failed to create request buffer", e1);
	}
	
	e2 := c.sendRequest(c.conn, buff);
	if e2 != nil {
		return nil, NewErrorWithCause(SYSTEM_ERR, "ServiceRequest(): failed to send request", e2);
	}
	
	resp, e3 := GetResponse(c.reader, cmd);
	if e3 != nil {
		return nil, NewErrorWithCause(SYSTEM_ERR, "ServiceRequest(): failed to get response", e3);
	}
	
	if resp.IsError() {
		log.Stderr("[Go-Redis] REDIS ERROR: ", resp.GetMessage());
		return nil, NewRedisError(resp.GetMessage());
	}
	return;
}

// Creates and opens a new SyncConnection.
//
func NewSyncConnection (spec *ConnectionSpec) (c SyncConnection, err os.Error) {
	hdl := new(_connection);
	if hdl == nil { 
		errmsg := "in NewSyncConnection(): Failed to allocate a new _connection";
		log.Stderr(errmsg);
		return nil, os.NewError (errmsg);
	}
	addr := spec.Addr();
	hdl.conn, err = net.Dial(TCP, "", addr);
	switch {
		case err != nil:
			err = NewErrorWithCause(SYSTEM_ERR, "Could not open connection", err);
		case hdl.conn == nil:
			err = NewError(SYSTEM_ERR, "Could not open connection");
		default:
//			log.Stdout("[Go-Redis] Opened SynchConnection connection to ", addr);
			hdl.reader = bufio.NewReader(hdl.conn);	
			c = hdl;
	}
	return;
}

// ----------------------------------------------------------------------------
// internal ops
// ----------------------------------------------------------------------------

func (hdl _connection) sendRequest (writer io.Writer, data []byte) os.Error {
	if writer == nil {
		log.Stderr ("sendRequest invoked with NIL writer!");
		return os.NewError("<BUG> illegal argument in sendRequest: nil writer.");
	}
	n, e1 := writer.Write(data);
	if e1 != nil {
		log.Stderr ("error on Write: ", e1);
	}
	if n < 6 {
		log.Stderr ("didn't write the whole data: ", n);
		return os.NewError("<BUG> lazy programmer didn't finish sendRequest...");
	}
	return e1;
}

