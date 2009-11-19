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
//	"reflect";
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
	if debug() {log.Stdout ("[Go-Redis] Closed connection: ", hdl);}
	return err;
}

// Implementation of SyncConnection.ServiceRequest.
//
func (c _connection) ServiceRequest (cmd *Command, args ...) (resp Response, err Error) {
	
	// TODO: need to consider errors here -- assuming it is always ok to write to Buffer ...
	buff, e1 := CreateRequestBytes(cmd, args);
	if e1 != nil {
		if debug() {log.Stderr("[Go-Redis] ERROR [", cmd.Code, "] :","CreateRequestBytes(): failed to get response", e1);}
		return nil, NewErrorWithCause(SYSTEM_ERR, "ServiceRequest(): failed to create request buffer", e1);
	}
	
	e2 := sendRequest(c.conn, buff);
	if e2 != nil {
		if debug() {log.Stderr("[Go-Redis] ERROR [", cmd.Code, "] :","sendRequest(): failed to send request", e2);}
		return nil, NewErrorWithCause(SYSTEM_ERR, "ServiceRequest(): failed to send request", e2);
	}
	
	resp, e3 := GetResponse(c.reader, cmd);
	if e3 != nil {
		if debug() {log.Stderr("[Go-Redis] ERROR [", cmd.Code, "] :","GetResponse(): failed to get response", e3);}
		return nil, NewErrorWithCause(SYSTEM_ERR, "ServiceRequest(): failed to get response", e3);
	}
	
	if resp.IsError() {
		if debug() {log.Stderr("[Go-Redis] REDIS ERROR [", cmd.Code, "] :",resp.GetMessage());}
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
		if debug() {log.Stderr(errmsg);}
		return nil, os.NewError (errmsg);
	}
	addr := spec.Addr();
	hdl.conn, err = net.Dial(TCP, "", addr);
	switch {
		case err != nil:
			err = NewErrorWithCause(SYSTEM_ERR, "Could not open connection", err);
			if debug() {log.Stderr(err);}
		case hdl.conn == nil:
			err = NewError(SYSTEM_ERR, "Could not open connection");
			if debug() {log.Stderr(err);}
		default:
			if debug() {log.Stdout("[Go-Redis] Opened SynchConnection connection to ", addr);}
			hdl.reader = bufio.NewReader(hdl.conn);	
			c = hdl;
	}
	return;
}
func newConnection (spec *ConnectionSpec) (c *_connection, err os.Error) {
	hdl := new(_connection);
	if hdl == nil { 
		errmsg := "in NewSyncConnection(): Failed to allocate a new _connection";
		if debug() {log.Stderr(errmsg);}
		return nil, os.NewError (errmsg);
	}
	addr := spec.Addr();
	hdl.conn, err = net.Dial(TCP, "", addr);
	switch {
		case err != nil:
			err = NewErrorWithCause(SYSTEM_ERR, "Could not open connection", err);
			if debug() {log.Stderr(err);}
		case hdl.conn == nil:
			err = NewError(SYSTEM_ERR, "Could not open connection");
			if debug() {log.Stderr(err);}
		default:
			if debug() {log.Stdout("[Go-Redis] Opened SynchConnection connection to ", addr);}
			hdl.reader = bufio.NewReader(hdl.conn);	
			c = hdl;
	}
	return;
}

// ----------------------------------------------------------------------------
// Asynchronous connections
// ----------------------------------------------------------------------------


// Defines the service contract supported by asynchronous (Request/FutureReply)
// connections.

type AsyncConnection interface {
	QueueRequest (cmd *Command, args ...) (*PendingResponse, Error);
}

type pendingRequest struct {
	cmd *Command;
//	args *reflect.StructValue;
//	args interface{};
	outbuff []byte;
	future interface{};
}

// control structure used by asynch connections.
//
type _async_connection struct {
	super 		*_connection;
	pending_reqs 	chan *pendingRequest;
	pending_resps 	chan *pendingRequest;	
	
//	req_handler  chan interface{};
//	resp_handler chan interface{};
}

// Creates and opens a new SyncConnection.
//
func NewAsynchConnection (spec *ConnectionSpec) (c AsyncConnection, err os.Error) {
	hdl := new(_async_connection);
	if hdl == nil { 
		errmsg := "in NewAsynchConnection(): Failed to allocate a new _async_connection";
		if debug() {log.Stderr(errmsg);}
		return nil, os.NewError (errmsg);
	}
	super, serr := newConnection (spec);
	if serr != nil {
		if debug () { log.Stderr(serr); }
		return nil, serr;
	}
	hdl.super = super;
	hdl.pending_reqs = make (chan *pendingRequest, 1000000);
	hdl.pending_resps = make (chan *pendingRequest, 1000000);

	go hdl.processRequests();
	go hdl.processResponses();
			
	return hdl, nil;
}

func (c _async_connection) processRequests ()  {
	if debug () {log.Stdout("processing requests for connection: ", c);}
	for {
		if debug() {log.Stdout("==>> process requests: ==>> about to TAKE FROM pending REQ ...");}
		req := <-c.pending_reqs;
		if debug() {log.Stdout("==>> process requests: ==>> ... TOOK cmd", req.cmd.Code, "args", req.outbuff);}
		
//		buff, e := CreateRequestBytes(req.cmd, req.args);
//		if e!=nil {
//			log.Stderr("PROBLEM: ", e);
//		}
		sendRequest(c.super.conn, req.outbuff);
		if debug() {log.Stdout("==>> process requests: ==>> sent request");}
		req.outbuff = nil;
		
		if debug() {log.Stdout("==>> process requests: ==>> about to QUEUE TO PENDING ...");}
		c.pending_resps<- req;
		if debug() {log.Stdout("==>> process requests: ==>> queued to pending responses");}
	}
}

func (c _async_connection) processResponses () {
	if debug () {log.Stdout("processing responses for connection: ", c);}
	for {
		if debug() {log.Stdout("<<== process response: about to TAKE FROM CHANNEL ...");}
		req := <-c.pending_resps;
		if debug() {log.Stdout("<<== process response: ... TOOK pending request for:  ", req);}
		
		reader := c.super.reader;
		cmd := req.cmd;
		
		r, e3 := GetResponse (reader, cmd);
		if e3!= nil {
			log.Stderr("<<== <PROBLEM!!>", e3);
			// this is a SYSTEM_ERR, such as broken connection, etc.
			// TDB: connection must handle this, but right now doesn't
			
		}
		if r.IsError() {
			errorResponse := NewRedisError(r.GetMessage());
			if debug() {log.Stdout("<<== process responses", errorResponse);}
			req.future.(FutureResult).onError(errorResponse);
		}
		else {
			switch cmd.RespType { // this is actually the raw response type ...
			case BOOLEAN:
				if debug() {log.Stdout("<<== process responses: about to SET on future ...");}
				req.future.(FutureBool).set(r.GetBooleanValue());
				if debug() {log.Stdout("<<== process responses: ... DID SET on future");}

			case BULK: 			
				req.future.(FutureBytes).set(r.GetBulkData());

			case MULTI_BULK:	
				req.future.(FutureBytesArray).set(r.GetMultiBulkData());

			case NUMBER:			
				req.future.(FutureInt64).set(r.GetNumberValue());

			case STATUS:		
				req.future.(FutureString).set(r.GetStringValue());

			case STRING:		
				req.future.(FutureString).set(r.GetStringValue());

		//	case VIRTUAL:		// FutureString?
		//	    resp, err = getVirtualResponse ();
			}
		}
	}
}

type PendingResponse struct {
	future interface{}
}
type _asyncResponse struct {
	future interface {}
}
func (asr _asyncResponse) getFutureResult() (FutureResult) {
	return asr.future.(FutureResult);
}
func (c _async_connection) QueueRequest (cmd *Command, v ...) (*PendingResponse, Error) {

	// create the pending request
	//
	buff, e1 := CreateRequestBytes(cmd, v);
	if e1 != nil {
		error := NewErrorWithCause(SYSTEM_ERR, "Failed to create request bytes", e1);
		if debug() { log.Stderr(error); }
		return nil, error;
	}
	request := new (pendingRequest);
	request.cmd = cmd;
	request.outbuff = buff;

	
	// create its specific future type
	switch cmd.RespType {
		case BOOLEAN:
			request.future = newFutureBool();
			if debug() {log.Stderr("QueueRequest: request.future", request.future);}
		case BULK: 			
			request.future = newFutureBytes();
		case MULTI_BULK:	
			request.future = newFutureBytesArray();
		case NUMBER:			
			request.future = newFutureInt64();
		case STATUS:		
			request.future = newFutureString();
		case STRING:		
			request.future = newFutureString();
	}
	
	// create pending response to be returned to caller
	// both point to the same future
	//	
	response := new(PendingResponse);
	response.future = request.future;
	
	// send it to the pending requests channel to queue the request
	//
	c.pending_reqs<- request;
	
	// done.
	return response, nil ;
}


// ----------------------------------------------------------------------------
// internal ops
// ----------------------------------------------------------------------------

//func (hdl _connection) sendRequest (writer io.Writer, data []byte) os.Error {
func sendRequest (writer io.Writer, data []byte) os.Error {
	if debug() { log.Stdout("sendRequest: data size", len(data), "data:", data);}
	if writer == nil {
		log.Stderr ("sendRequest invoked with NIL writer!");
		return os.NewError("<BUG> illegal argument in sendRequest: nil writer.");
	}
	n, e1 := writer.Write(data);
	if e1 != nil {
		log.Stderr ("error on Write: ", e1);
	}
	if n < len(data) {
		log.Stderr ("didn't write the whole data: ", n);
		return os.NewError("<BUG> lazy programmer didn't finish sendRequest...");
	}
	return e1;
}

