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
	"net";
	"fmt";
	"strconv";
	"os";
	"io";
	"bufio";
	"log";
)

const (
	TCP = "tcp";
	LOCALHOST = "127.0.0.1";
	ns1Sec = 1000000;
	ns1MSec = 1000;
)

// various default sizes for the connections
// exported for user convenience if nedded
const(
	DefaultReqChanSize			= 100000;
	DefaultRespChanSize			= 100000;
	
	DefaultTCPReadBuffSize		= 1024 * 256;
	DefaultTCPWriteBuffSize		= 1024 * 256;
	DefaultTCPReadTimeoutNSecs	= ns1Sec * 10;
	DefaultTCPWriteTimeoutNSecs	= ns1Sec * 10;
	DefaultTCPLinger			= -1;
	DefaultTCPKeepalive			= true;
)

// Redis specific default settings
// exported for user convenience if nedded
const (
	DefaultRedisPassword = "";
	DefaultRedisDB = 0;
	DefaultRedisPort = 6379;
	DefaultRedisHost = LOCALHOST;
)
// ----------------------------------------------------------------------------
// Connection ConnectionSpec
// ----------------------------------------------------------------------------

// Defines the set of parameters that are used by the client connections
//
type ConnectionSpec struct {
	host		string;
	port		int;
	password	string;
	db			int;
	// tcp specific sspecs
	rBufSize	int;
	wBufSize	int;
	rTimeout	int64;
	wTimeout	int64;
	keepalive	bool;
	lingerspec	int; // -n: finish io; 0: discard, +n: wait for n secs to finish
}

// Creates a ConnectionSpec using default settings.
// using the DefaultXXX consts of redis package.
func DefaultSpec () *ConnectionSpec {
	return &ConnectionSpec {
		DefaultRedisHost,
		DefaultRedisPort,
		DefaultRedisPassword,
		DefaultRedisDB,
		DefaultTCPReadBuffSize,
		DefaultTCPWriteBuffSize,
		DefaultTCPReadTimeoutNSecs,
		DefaultTCPWriteTimeoutNSecs,
		DefaultTCPKeepalive,
		DefaultTCPLinger
	};
}

// Sets the db for connection spec and returns the reference
// Note that you should not this after you have already connected.
func (spec *ConnectionSpec) Db(db int) *ConnectionSpec {
	spec.db = db;
	return spec;
}

// Sets the host for connection spec and returns the reference
// Note that you should not this after you have already connected.
func (spec *ConnectionSpec) Host(host string) *ConnectionSpec {
	spec.host = host;
	return spec;
}

// Sets the port for connection spec and returns the reference
// Note that you should not this after you have already connected.
func (spec *ConnectionSpec) Port(port int) *ConnectionSpec {
	spec.port = port;
	return spec;
}

// Sets the password for connection spec and returns the reference
// Note that you should not this after you have already connected.
func (spec *ConnectionSpec) Password(password string) *ConnectionSpec {
	spec.password = password;
	return spec;
}

// return the address as string.
func (spec *ConnectionSpec) Addr () string {
	return spec.host + ":" + strconv.Itoa(int(spec.port));
//	return fmt.Sprintf("%s:%d", spec.host, spec.port);
}

// ----------------------------------------------------------------------------
// Generic Conn handle and methods
// ----------------------------------------------------------------------------

// General control structure used by connections.
//
type syncConnHDL struct {
	spec 	*ConnectionSpec;
	conn 	net.Conn;
	reader 	*bufio.Reader;
}

// Creates and opens a new connection to server per ConnectionSpec.
// The new connection is wrapped by a new syncConnHDL with its bufio.Reader
// delegating to the net.Conn's reader. 
//
func newConnHDL (spec *ConnectionSpec) (hdl *syncConnHDL, err os.Error) {
	here := "newConnHDL";

	if hdl = new(syncConnHDL); hdl == nil { 
		return nil, withNewError (fmt.Sprintf("%s(): failed to allocate syncConnHDL", here));
	}
	addr := spec.Addr(); 
	raddr, e:= net.ResolveTCPAddr(addr); 
	if e != nil {
		return nil, withNewError (fmt.Sprintf("%s(): failed to resolve remote address %s", here, addr));
	}	
	conn, e:= net.DialTCP(TCP, nil, raddr);
	switch {
		case e != nil:
			err = withOsError (fmt.Sprintf("%s(): could not open connection", here), e);
		case conn == nil:
			err = withNewError (fmt.Sprintf("%s(): net.Dial returned nil, nil (?)", here));
		default:
			configureConn(conn, spec);
			hdl.conn = conn;
			hdl.reader, e = bufio.NewReaderSize(conn, 4096);
			if debug() {log.Stdout("[Go-Redis] Opened SynchConnection connection to ", addr);}
	}
	return hdl, err;
}

func configureConn (conn *net.TCPConn, spec *ConnectionSpec) {
	// these two -- the most important -- are causing problems on my osx/64
	// where a "service unavailable" pops up in the async reads 
	// but we absolutely need to be able to use timeouts.
//			conn.SetReadTimeout(spec.rTimeout);	
//			conn.SetWriteTimeout(spec.wTimeout);	
	conn.SetLinger(spec.lingerspec);
	conn.SetKeepAlive(spec.keepalive);
	conn.SetReadBuffer(spec.rBufSize);
	conn.SetWriteBuffer(spec.wBufSize);
}
// closes the syncConnHDL's net.Conn connection.
// Is public so that syncConnHDL struct can be used as SyncConnection (TODO: review that.)
//
func (hdl syncConnHDL) Close() os.Error {
	err := hdl.conn.Close();
	if debug() {log.Stdout ("[Go-Redis] Closed connection: ", hdl);}
	return err;
}


// ----------------------------------------------------------------------------
// Connection SyncConnection
// ----------------------------------------------------------------------------

// Defines the service contract supported by synchronous (Request/Reply)
// connections.

type SyncConnection interface {
//	ServiceRequest (cmd *Command, args ...) (Response, Error);
	ServiceRequest (cmd *Command, args [][]byte) (Response, Error);
	Close () os.Error;
}

// Creates a new SyncConnection using the provided ConnectionSpec
func NewSyncConnection (spec *ConnectionSpec) (c SyncConnection, err os.Error) {
	return newConnHDL (spec);
}

// Implementation of SyncConnection.ServiceRequest.
//
func (chdl *syncConnHDL) ServiceRequest (cmd *Command, args [][]byte) (resp Response, err Error) {
	here := "syncConnHDL.ServiceRequest";
	errmsg := "";
	ok := false;
	buff, e := CreateRequestBytes2 (cmd, args);  // 2<<<
	if e == nil {
		e = sendRequest(chdl.conn, buff);
		if e == nil {
			resp, e = GetResponse(chdl.reader, cmd);
			if e == nil {
				if resp.IsError() {
					redismsg := fmt.Sprintf(" [%s]: %s", cmd.Code, resp.GetMessage());
					err = NewRedisError(redismsg);
				}
				ok = true;
			}
			else { errmsg = fmt.Sprintf("%s(%s): failed to get response", here, cmd.Code); }
		}
		else { errmsg = fmt.Sprintf("%s(%s): failed to send request", here, cmd.Code); }
	}
	else { errmsg = fmt.Sprintf("%s(%s): failed to create request buffer", here, cmd.Code); } 
	
	if !ok {
		return resp, withError(NewErrorWithCause(SYSTEM_ERR, errmsg, e)); // log it on debug
	}
	
	return;
}

// ----------------------------------------------------------------------------
// Asynchronous connections
// ----------------------------------------------------------------------------

// Defines the service contract supported by asynchronous (Request/FutureReply)
// connections.

type AsyncConnection interface {
//	QueueRequest (cmd *Command, args ...) (*PendingResponse, Error);
	QueueRequest (cmd *Command, args [][]byte) (*PendingResponse, Error);
}

// Handle to a future response
type PendingResponse struct {
	future interface{}
}

// Creates and opens a new AsyncConnection and starts the goroutines for 
// request and response processing

func NewAsynchConnection (spec *ConnectionSpec) (AsyncConnection, os.Error) {
	async, err:= newAsyncHDL(spec);
//	go async.processRequests();
	go async.batchProcessRequestsI (spec);
	go async.processResponses();
	return async, err;
}

// Defines the data corresponding to a requested service call through the
// QueueRequest method of AsyncConnection

type pendingRequest struct {
	cmd			*Command;
	outbuff		*[]byte;
	future		interface{};
}

// control structure used by asynch connections.

type asyncConnHDL struct {
	super			*syncConnHDL;
	pending_reqs 	chan *pendingRequest;
	pending_resps 	chan *pendingRequest;	
	
	/*
	// TODO: we'll need these so the go routines can coordinate (on lifecycle and error
	// events) with the AsyncConnection	
	
	req_handler  chan interface{};
	resp_handler chan interface{};
	*/
}

// Creates a new asyncConnHDL with a new syncConnHDL as its delegated 'super'.
// Note it does not start the processing goroutines for the channels.

func newAsyncHDL (spec *ConnectionSpec) (async *asyncConnHDL, err os.Error) {
	here := "newAsynConnHDL";
	super, err := newConnHDL (spec);
	if err == nil {
		async = new(asyncConnHDL);
		if async != nil { 
			async.super = super;
			async.pending_reqs = make (chan *pendingRequest, DefaultReqChanSize);
			async.pending_resps = make (chan *pendingRequest, DefaultRespChanSize);
		}
		else {
			return nil, withNewError (fmt.Sprintf("%s(): failed to allocate asyncConnHDL", here));
		}
	}
	else {
		return nil, withOsError (fmt.Sprintf("%s(): Error creating syncConnHDL", here), err);
	}
			
	return;
}

func (c *asyncConnHDL) batchProcessRequestsI (spec *ConnectionSpec)  {
	if debug () {log.Stdout("begin processing requests for connection [using glued-writes]: ", c);}

	maxbytes := spec.wBufSize/2;
	buff, e := bufio.NewWriterSize(c.super.conn, maxbytes);
	if e!= nil {
		log.Stderr("<BUG> lazy programmer >> ERROR in processRequests goroutine");
	}
	for {
		bytecnt := 0;
		itemcnt := len(c.pending_reqs);
		flush := false;
		for itemcnt > 0 {
			flush = true;
			req := <-c.pending_reqs;
			e := sendRequest(buff, *req.outbuff);
			bytecnt += len(*req.outbuff);
			if e == nil {
				req.outbuff = nil;
				c.pending_resps<- req;
			}
			else {
				// put it back for now ...
				log.Stderr("<BUG> lazy programmer >> ERROR in processRequests goroutine -req requeued");
				c.pending_reqs<- req;
				break;
			}
			itemcnt = len(c.pending_reqs);
			if bytecnt > maxbytes { // i know ..
				break;
			}
		}
		if flush { buff.Flush(); fmt.Printf(""); /* magic */}
		else { fmt.Printf(""); /* magic */ }
	}
	if debug () {log.Stdout("stopped processing requests for connection: ", c);}
}

//// (as of now) used by a goroutine to process pending requests.
//func (c *asyncConnHDL) processRequests ()  {
//	if debug () {log.Stdout("begin processing requests for connection: ", c);}
//	for {
//		req := <-c.pending_reqs;
//		itemcnt := len(c.pending_reqs);
//		if itemcnt > 0 {
//			log.Stdout (">>> still in channel: ", itemcnt);
//		}
//		e := sendRequest(c.super.conn, *req.outbuff);
//		if e == nil {
//			req.outbuff = nil;
//			c.pending_resps<- req;
//		}
//		else {
//			// TODO: need a way for this goroutines to gracefully shutdown
//			// and let the owning connection know there are network issues
//			// & TBD
//			log.Stderr("<BUG> lazy programmer hasn't addressed failures in processRequests goroutine");
//			break;
//		}
//	}
//	if debug () {log.Stdout("stopped processing requests for connection: ", c);}
//}

// (as of now) used by a goroutine to process pending responses.

func (c *asyncConnHDL) processResponses () {
	if debug () {log.Stdout("begin processing responses for connection: ", c);}
	for {
		req:= <-c.pending_resps;
		reader:= c.super.reader;
		cmd:= req.cmd;
		
		r, e3:= GetResponse (reader, cmd);
		if e3!= nil {
			log.Stderr("<BUG> lazy programmer hasn't addressed failures in processResponses goroutine");
			log.Stderr(e3);
			break;
		}
		
		if r.IsError() {
			errorResponse := NewRedisError(r.GetMessage());
			req.future.(FutureResult).onError(errorResponse);
		}
		else {
			switch cmd.RespType {
			case BOOLEAN:
				req.future.(FutureBool).set(r.GetBooleanValue());

			case BULK: 			
				req.future.(FutureBytes).set(r.GetBulkData());

			case MULTI_BULK:	
				req.future.(FutureBytesArray).set(r.GetMultiBulkData());

			case NUMBER:			
				req.future.(FutureInt64).set(r.GetNumberValue());

			case STATUS:		
				req.future.(FutureString).set(r.GetMessage());

			case STRING:		
				req.future.(FutureString).set(r.GetStringValue());

		//	case VIRTUAL:		// FutureString?
		//	    resp, err = getVirtualResponse ();
			}
		}
	}
	if debug () {log.Stdout("stopped processing responses for connection: ", c);}
}

// Implementation of AsyncConnection.QueueRequest;

func (c *asyncConnHDL) QueueRequest (cmd *Command, args [][]byte) (*PendingResponse, Error) {
	here := "syncConnHDL.ServiceRequest";

	buff, e1 := CreateRequestBytes2(cmd, args);
	if e1 != nil {
		errmsg:= fmt.Sprintf("%s(%s): failed to create request bytes", here, cmd.Code);	
		return nil, withError (NewErrorWithCause(SYSTEM_ERR, errmsg, e1));
	}
	
	request := new (pendingRequest);
	request.cmd = cmd;
	request.outbuff = &buff;
	
	// create its specific future type
	switch cmd.RespType {
		case BOOLEAN:
			request.future = newFutureBool();
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


// Either writes all the bytes or it fails and returns an error
//
func sendRequest (w io.Writer, data []byte) (e os.Error) {
	here := "syncConnHDL.sendRequest";
	if w == nil {
		return withNewError (fmt.Sprintf("<BUG> in %s(): nil Writer", here));
	}
	
	n, e := w.Write(data);
	if e != nil {
		var msg string;
		switch {
		case e == os.EAGAIN:		
			// socket timeout -- don't handle that yet but may in future ..
			msg = fmt.Sprintf("%s(): timeout (os.EAGAIN) error on Write", here);
		default:
			// anything else
			msg = fmt.Sprintf("%s(): error on Write", here);
		}
		return withOsError(msg, e);
	}
	
	// doc isn't too clear but the underlying netFD may return n<len(data) AND
	// e == nil, but that's precisely what we're testing for here.  
	// presumably we can try sending the remaining bytes but that is precisely
	// what netFD.Write is doing (and it couldn't) so ...
	if n < len(data) {
		msg := fmt.Sprintf("%s(): connection Write wrote %d bytes only.", here, n);
		return withNewError(msg);
	}
	return;
}

