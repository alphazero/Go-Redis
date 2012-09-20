//   Copyright 2009-2012 Joubin Houshyar
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

package redis

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

// connection socket modes
const (
	TCP  = "tcp"  // tcp/ip socket
	UNIX = "unix" // unix domain socket
)

// various defaults for the connections
// exported for user convenience.
const (
	DefaultReqChanSize          = 1000000
	DefaultRespChanSize         = 1000000
	DefaultTCPReadBuffSize      = 1024 * 256
	DefaultTCPWriteBuffSize     = 1024 * 256
	DefaultTCPReadTimeoutNSecs  = 10 * time.Second
	DefaultTCPWriteTimeoutNSecs = 10 * time.Second
	DefaultTCPLinger            = 0 // -n: finish io; 0: discard, +n: wait for n secs to finish
	DefaultTCPKeepalive         = true
	DefaultHeartbeatSecs        = 1 * time.Second
	DefaultProtocol             = REDIS_DB
)

// Redis specific default settings
// exported for user convenience.
const (
	DefaultRedisPassword = ""
	DefaultRedisDB       = 0
	DefaultRedisPort     = 6379
	DefaultRedisHost     = "127.0.0.1"
)

// ----------------------------------------------------------------------------
// Connection ConnectionSpec
// ----------------------------------------------------------------------------

type Protocol int

const (
	REDIS_DB Protocol = iota
	REDIS_PUBSUB
)

func (p Protocol) String() string {
	switch p {
	case REDIS_DB:
		return "Protocol:REDIS_DB"
	case REDIS_PUBSUB:
		return "Protocol:PubSub"
	}
	return "BUG - unknown protocol value"
}

// Defines the set of parameters that are used by the client connections
//
type ConnectionSpec struct {
	host       string        // redis connection host
	port       int           // redis connection port
	password   string        // redis connection password
	db         int           // Redis connection db #
	rBufSize   int           // tcp read buffer size
	wBufSize   int           // tcp write buffer size
	rTimeout   time.Duration // tcp read timeout
	wTimeout   time.Duration // tcp write timeout
	keepalive  bool          // keepalive flag
	lingerspec int           // -n: finish io; 0: discard, +n: wait for n secs to finish
	reqChanCap int           // async request channel capacity - see DefaultReqChanSize
	rspChanCap int           // async response channel capacity - see DefaultRespChanSize
	heartbeat  time.Duration // 0 means no heartbeat
	protocol   Protocol
}

// Creates a ConnectionSpec using default settings.
// using the DefaultXXX consts of redis package.
func DefaultSpec() *ConnectionSpec {
	return &ConnectionSpec{
		DefaultRedisHost,
		DefaultRedisPort,
		DefaultRedisPassword,
		DefaultRedisDB,
		DefaultTCPReadBuffSize,
		DefaultTCPWriteBuffSize,
		DefaultTCPReadTimeoutNSecs,
		DefaultTCPWriteTimeoutNSecs,
		DefaultTCPKeepalive,
		DefaultTCPLinger,
		DefaultReqChanSize,
		DefaultRespChanSize,
		DefaultHeartbeatSecs,
		DefaultProtocol,
	}
}

// Sets the db for connection spec and returns the reference
// Note that you should not this after you have already connected.
func (spec *ConnectionSpec) Db(db int) *ConnectionSpec {
	spec.db = db
	return spec
}

// Sets the host for connection spec and returns the reference
// Note that you should not this after you have already connected.
func (spec *ConnectionSpec) Host(host string) *ConnectionSpec {
	spec.host = host
	return spec
}

// Sets the port for connection spec and returns the reference
// Note that you should not this after you have already connected.
func (spec *ConnectionSpec) Port(port int) *ConnectionSpec {
	spec.port = port
	return spec
}

// Sets the password for connection spec and returns the reference
// Note that you should not this after you have already connected.
func (spec *ConnectionSpec) Password(password string) *ConnectionSpec {
	spec.password = password
	return spec
}

// return the address as string.
func (spec *ConnectionSpec) Heartbeat(period time.Duration) *ConnectionSpec {
	spec.heartbeat = period
	return spec
}

// return the address as string.
func (spec *ConnectionSpec) Protocol(protocol Protocol) *ConnectionSpec {
	spec.protocol = protocol
	return spec
}

// ----------------------------------------------------------------------------
// SyncConnection API
// ----------------------------------------------------------------------------

// Defines the service contract supported by synchronous (Request/Reply)
// connections.

type SyncConnection interface {
	ServiceRequest(cmd *Command, args [][]byte) (Response, Error)
}

// ----------------------------------------------------------------------------
// AsyncConnection API
// ----------------------------------------------------------------------------

// Defines the service contract supported by asynchronous (Request/FutureReply)
// connections.

type AsyncConnection interface {
	QueueRequest(cmd *Command, args [][]byte) (*PendingResponse, Error)
}

// Handle to a future response
type PendingResponse struct {
	future interface{} // TODO: stop using runtime casts
}

// ----------------------------------------------------------------------------
// PubSubConnection API
// ----------------------------------------------------------------------------

// Defines the service contract supported by asynchronous (Request/FutureReply)
// connections.

// REVU - does this need its own special connection structure or
// can we get away with just using the asyncConnHdl structure?
// REVU - BEST is to reuse asyncConnHdl and simple have that support this
//	      interface.
type PubSubConnection interface {
	// TODO - this connector needs to be given a channel to feed responses
	// into as it does not use futures.
	//	SetOutputChannel(<-chan message) bool
	// TODO - then either a generic ServiceRequest or explicit Sub/UnSub/Quit
	// REVU - best to keep it simple and use a single method so using Command
	//	    -
}

// ----------------------------------------------------------------------------
// Generic Conn handle and methods - supports SyncConnection interface
// ----------------------------------------------------------------------------

// General control structure used by connections.
//
type connHdl struct {
	spec      *ConnectionSpec
	conn      net.Conn // may want to change this to TCPConn - TODO REVU
	reader    *bufio.Reader
	connected bool // TODO
}

// Returns minimal info string for logging, etc
func (c *connHdl) String() string {
	return fmt.Sprintf("conn<redis-server@%s:%d [db %d]>", c.spec.host, c.spec.port, c.spec.db)
}

// Creates and opens a new connection to server per ConnectionSpec.
// The new connection is wrapped by a new connHdl with its bufio.Reader
// delegating to the net.Conn's reader.
//
func newConnHdl(spec *ConnectionSpec) (hdl *connHdl, err Error) {
	loginfo := "newConnHdl"

	if hdl = new(connHdl); hdl == nil {
		return nil, NewError(SYSTEM_ERR, fmt.Sprintf("%s(): failed to allocate connHdl", loginfo))
	}

	var mode, addr string
	if spec.port == 0 { // REVU - no special values (it was a contrib) TODO add flag to connspec.
		mode = UNIX
		addr = spec.host
	} else {
		mode = TCP
		addr = fmt.Sprintf("%s:%d", spec.host, spec.port)
		_, e := net.ResolveTCPAddr(TCP, addr)
		if e != nil {
			msg := fmt.Sprintf("%s(): failed to resolve remote address %s", loginfo, addr)
			return nil, NewErrorWithCause(SYSTEM_ERR, msg, e)
		}
	}

	conn, e := net.Dial(mode, addr)
	switch {
	case e != nil:
		err = NewErrorWithCause(SYSTEM_ERR, fmt.Sprintf("%s(): could not open connection", loginfo), e)
		return nil, withError(err)
	case conn == nil:
		err = NewError(SYSTEM_ERR, fmt.Sprintf("%s(): net.Dial returned nil, nil (?)", loginfo))
		return nil, withError(err)
	default:
		configureConn(conn, spec)
		hdl.spec = spec
		hdl.conn = conn
		hdl.connected = true
		bufsize := 4096
		hdl.reader = bufio.NewReaderSize(conn, bufsize)
	}
	return hdl, nil
}

func configureConn(conn net.Conn, spec *ConnectionSpec) {
	// REVU [jh] - TODO look into this 09-13-2012
	// these two -- the most important -- are causing problems on my osx/64
	// where a "service unavailable" pops up in the async reads
	// but we absolutely need to be able to use timeouts.
	//			conn.SetReadTimeout(spec.rTimeout);
	//			conn.SetWriteTimeout(spec.wTimeout);
	if tcp, ok := conn.(*net.TCPConn); ok {
		tcp.SetLinger(spec.lingerspec)
		tcp.SetKeepAlive(spec.keepalive)
		tcp.SetReadBuffer(spec.rBufSize)
		tcp.SetWriteBuffer(spec.wBufSize)
	}
}

// connect event handler will issue AUTH/SELECT on new connection
// if required.
func (c *connHdl) connect() (e Error) {
	if c.spec.password != DefaultRedisPassword {
		_, e = c.ServiceRequest(&AUTH, [][]byte{[]byte(c.spec.password)})
		if e != nil {
			log.Printf("<ERROR> Authentication failed - %s", e.Message())
			return
		}
	}
	if c.spec.db != DefaultRedisDB {
		_, e = c.ServiceRequest(&SELECT, [][]byte{[]byte(fmt.Sprintf("%d", c.spec.db))})
		if e != nil {
			log.Printf("<ERROR> REDIS_DB Select failed - %s", e.Message())
			return
		}
	}
	// REVU - pretty please TODO do the customized log
	//	log.Printf("<INFO> %s - CONNECTED", c)
	return
}

func (hdl *connHdl) disconnect() Error {
	// silently ignore repeated calls to closed connections
	if hdl.connected {
		if e := hdl.conn.Close(); e != nil {
			return NewErrorWithCause(SYSTEM_ERR, "on connHdl.Close()", e)
		}
		hdl.connected = false
		// REVU - pretty please TODO do the customized log
		//		log.Printf("<INFO> %s - DISCONNECTED", hdl)
	}
	return nil
}

// Creates a new SyncConnection using the provided ConnectionSpec.
// Note that this function will also connect to the specified redis server.
func NewSyncConnection(spec *ConnectionSpec) (c SyncConnection, err Error) {
	connHdl, e := newConnHdl(spec)
	if e != nil {
		return nil, e
	}

	e = connHdl.connect()
	return connHdl, e
}

// Implementation of SyncConnection.ServiceRequest.
func (hdl *connHdl) ServiceRequest(cmd *Command, args [][]byte) (resp Response, err Error) {
	loginfo := "connHdl.ServiceRequest"
	errmsg := ""
	if !hdl.connected {
		return nil, NewError(REDIS_ERR, "Connection"+hdl.String()+" is closed")
	}
	if cmd == &QUIT {
		return nil, hdl.disconnect()
	}
	ok := false
	buff, e := CreateRequestBytes(cmd, args) // 2<<<
	if e == nil {
		e = sendRequest(hdl.conn, buff)
		if e == nil {
			// REVU - this demands resp to be non-nil even in case of io errors
			// TODO - refactor this
			resp, e = GetResponse(hdl.reader, cmd)
			//			fmt.Printf("DEBUG-TEMP - resp: %s\n", resp)   // REVU REMOVE TODO
			if e == nil {
				if resp.IsError() {
					redismsg := fmt.Sprintf(" [%s]: %s", cmd.Code, resp.GetMessage())
					err = NewRedisError(redismsg)
				}
				ok = true
			} else {
				errmsg = fmt.Sprintf("%s(%s): failed to get response", loginfo, cmd.Code)
			}
		} else {
			errmsg = fmt.Sprintf("%s(%s): failed to send request", loginfo, cmd.Code)
		}
	} else {
		errmsg = fmt.Sprintf("%s(%s): failed to create request buffer", loginfo, cmd.Code)
	}

	if !ok {
		return resp, withError(NewErrorWithCause(SYSTEM_ERR, errmsg, e)) // log it on debug
	}

	return
}

// ----------------------------------------------------------------------------
// Asynchronous connection handle and friends
// ----------------------------------------------------------------------------

type status_code byte

const (
	ok status_code = iota
	info
	warning
	error_
	reqerr
	inierr
	snderr
	rcverr
)

type interrupt_code byte

// connection process control interrupt codes
const (
	_ interrupt_code = iota
	start
	pause
	stop
)

type workerCtl chan interrupt_code

type event_code byte

const (
	_ event_code = iota
	ready
	working
	faulted
	quit_processed
)

type workerStatus struct {
	id       int
	event    event_code
	taskinfo *taskStatus
	ctlchan  *workerCtl
}
type taskStatus struct {
	code  status_code
	error error
}

var ok_status = taskStatus{ok, nil}

type workerTask func(*asyncConnHdl, workerCtl) (*interrupt_code, *taskStatus)

// Defines the data corresponding to a requested service call through the
// QueueRequest method of AsyncConnection
// not used yet.
type asyncRequestInfo struct {
	id      int64
	stat    status_code
	cmd     *Command
	outbuff *[]byte
	future  interface{}
	error   Error
}
type asyncReqPtr *asyncRequestInfo

// control structure used by asynch connections.
const (
	heartbeatworker int = iota
	manager
	requesthandler
	responsehandler
)

type asyncConnHdl struct {
	super  *connHdl
	writer *bufio.Writer

	nextid       int64
	pendingReqs  chan asyncReqPtr
	pendingResps chan asyncReqPtr
	faults       chan asyncReqPtr

	reqProcCtl   workerCtl
	rspProcCtl   workerCtl
	heartbeatCtl workerCtl
	managerCtl   workerCtl

	feedback chan workerStatus

	shutdown   chan bool
	isShutdown bool
}

func (c *asyncConnHdl) String() string {
	return fmt.Sprintf("async [%s] %s", c.spec().protocol, c.super.String())
}

func (c *asyncConnHdl) spec() *ConnectionSpec {
	return c.super.spec
}

// Creates a new asyncConnHdl with a new connHdl as its delegated 'super'.
// Note it does not start the processing goroutines for the channels.
// REVU - PubSub could be checked here
func newAsyncConnHdl(spec *ConnectionSpec) (async *asyncConnHdl, err Error) {
	connHdl, err := newConnHdl(spec)
	if err == nil && connHdl != nil {
		async = new(asyncConnHdl)
		if async != nil {
			async.super = connHdl
			//			var e error
			async.writer = bufio.NewWriterSize(connHdl.conn, spec.wBufSize)

			async.pendingReqs = make(chan asyncReqPtr, spec.reqChanCap)
			async.pendingResps = make(chan asyncReqPtr, spec.rspChanCap)
			async.faults = make(chan asyncReqPtr, spec.reqChanCap) // REVU - not sure about sizing

			async.reqProcCtl = make(workerCtl)
			async.rspProcCtl = make(workerCtl)
			async.heartbeatCtl = make(workerCtl)
			async.managerCtl = make(workerCtl)

			async.feedback = make(chan workerStatus)
			async.shutdown = make(chan bool, 1)

			async.isShutdown = false

			return
		}
	}
	// reached on errors only
	if debug() {
		log.Println("Error creating asyncConnHdl: ", err)
	}
	return nil, err
}

// Creates and opens a new AsyncConnection and starts the goroutines for
// request and response processing
// interaction with redis (AUTH &| SELECT)
func NewAsynchConnection(spec *ConnectionSpec) (conn AsyncConnection, err Error) {
	var async *asyncConnHdl
	if async, err = newAsyncConnHdl(spec); err == nil {
		if err = async.connect(); err == nil {
			async.startup()
		} else {
			async = nil // do not return a ref if we could not connect
		}
	}
	return async, err
}

// ----------------------------------------------------------------------------
// asyncConnHdl support for AsyncConnection interface
// ----------------------------------------------------------------------------

// TODO Quit - see REVU notes added for adding Quit to async in body
func (c *asyncConnHdl) QueueRequest(cmd *Command, args [][]byte) (*PendingResponse, Error) {

	if c.isShutdown {
		return nil, NewError(SYSTEM_ERR, "Connection is shutdown.")
	}

	select {
	case <-c.shutdown:
		c.isShutdown = true
		c.shutdown <- true // put it back REVU likely to be a bug under heavy load ..
		//		log.Println("<DEBUG> we're shutdown and not accepting any more requests ...")
		return nil, NewError(SYSTEM_ERR, "Connection is shutdown.")
	default:
	}

	future := CreateFuture(cmd)
	request := &asyncRequestInfo{0, 0, cmd, nil, future, nil}

	buff, e1 := CreateRequestBytes(cmd, args)
	if e1 == nil {
		request.outbuff = &buff
		c.pendingReqs <- request // REVU - opt 1 TODO is handling QUIT and sending stop to workers
	} else {
		errmsg := fmt.Sprintf("Failed to create asynchrequest - %s aborted", cmd.Code)
		request.stat = inierr
		request.error = NewErrorWithCause(SYSTEM_ERR, errmsg, e1) // only makes sense if using go ...
		request.future.(FutureResult).onError(request.error)

		return nil, request.error // remove if restoring go
	}
	//}();
	// done.
	return &PendingResponse{future}, nil
}

// ----------------------------------------------------------------------------
// asyncConnHdl support for PubSubConnection interface
// ----------------------------------------------------------------------------

//// TODO Quit - see REVU notes added for adding Quit to async in body
//func (c *asyncConnHdl) ServicePubSubRequest(cmd *Command, args [][]byte) (, Error) {
//
//	if c.isShutdown {
//		return nil, NewError(SYSTEM_ERR, "Connection is shutdown.")
//	}
//
//	select {
//	case <-c.shutdown:
//		c.isShutdown = true
//		c.shutdown <- true // put it back REVU likely to be a bug under heavy load ..
//		//		log.Println("<DEBUG> we're shutdown and not accepting any more requests ...")
//		return nil, NewError(SYSTEM_ERR, "Connection is shutdown.")
//	default:
//	}
//
//	future := CreateFuture(cmd)
//	request := &asyncRequestInfo{0, 0, cmd, nil, future, nil}
//
//	buff, e1 := CreateRequestBytes(cmd, args)
//	if e1 == nil {
//		request.outbuff = &buff
//		c.pendingReqs <- request // REVU - opt 1 TODO is handling QUIT and sending stop to workers
//	} else {
//		errmsg := fmt.Sprintf("Failed to create asynchrequest - %s aborted", cmd.Code)
//		request.stat = inierr
//		request.error = NewErrorWithCause(SYSTEM_ERR, errmsg, e1) // only makes sense if using go ...
//		request.future.(FutureResult).onError(request.error)
//
//		return nil, request.error // remove if restoring go
//	}
//	//}();
//	// done.
//	return &PendingResponse{future}, nil
//}

// ----------------------------------------------------------------------------
// asyncConnHdl internal ops
// ----------------------------------------------------------------------------

// connect event handler.
// See connHdl#connect
func (c *asyncConnHdl) connect() (e Error) {
	return c.super.connect()
}

// REVU - TODO opt 2 for Quit here
func (c *asyncConnHdl) disconnect() (e Error) {

	panic("asyncConnHdl.disconnect NOT IMLEMENTED!")
	//	return
}

// responsible for managing the various moving parts of the asyncConnHdl
func (c *asyncConnHdl) startup() {

	protocol := c.spec().protocol

	go c.worker(manager, "manager", managementTask, c.managerCtl, nil)
	c.managerCtl <- start

	// heartbeat only on REDIS_DB protocol
	if protocol == REDIS_DB {
		go c.worker(heartbeatworker, "heartbeat", heartbeatTask, c.heartbeatCtl, c.feedback)
		c.heartbeatCtl <- start
	}

	go c.worker(heartbeatworker, "request-processor", reqProcessingTask, c.reqProcCtl, c.feedback)
	c.reqProcCtl <- start

	var rspProcTask workerTask
	switch protocol {
	case REDIS_DB:
		rspProcTask = dbRspProcessingTask
	case REDIS_PUBSUB:
		rspProcTask = msgProcessingTask
	}
	go c.worker(heartbeatworker, "response-processor", rspProcTask, c.rspProcCtl, c.feedback)
	c.rspProcCtl <- start

	// REVU - pretty please TODO do the customized log
	//	log.Printf("<INFO> %s - READY", c)
}

// This could find a happy home in a generalized worker package ...
// TODO
func (c *asyncConnHdl) worker(id int, name string, task workerTask, ctl workerCtl, fb chan workerStatus) {
	// REVU - pretty please TODO do the customized log
	//	log.Printf("<INFO> %s - %s STARTED.", c, name)
	var signal interrupt_code
	var tstat *taskStatus

	// todo: add startup hook for worker

await_signal:
	//	log.Println(name, "_worker: await_signal.")
	signal = <-ctl

on_interrupt:
	//		log.Println(name, "_worker: on_interrupt: ", signal);
	switch signal {
	case stop:
		goto before_stop
	case pause:
		goto await_signal
	case start:
		goto work
	}

work:
	//		fmt.Println(name, "_worker: work!");
	select {
	case signal = <-ctl:
		goto on_interrupt
	default:
		is, stat := task(c, ctl) // todo is a task context type
		if stat == nil {
			log.Println("<BUG> nil stat from worker ", name)
		}
		if stat.code != ok {
			//			fmt.Println(name, "_worker: task error!")
			tstat = stat
			goto on_error
		} else if is != nil {
			signal = *is
			goto on_interrupt
		}
		goto work
	}

on_error:
	//log.Println(name, "_worker: on_error!")
	// TODO: log it, send it, and go back to wait_start:
	log.Println(name, "_worker task raised error: ", tstat)
	fb <- workerStatus{id, faulted, tstat, &ctl}
	goto await_signal

before_stop:
	//	fmt.Println(name, "_worker: before_stop!")
	// TODO: add shutdown hook for worker

	// REVU - pretty please TODO do the customized log
	//	log.Printf("<INFO> %s - %s STOPPED.", c, name)
}

// ----------------------------------------------------------------------------
// asynchronous tasks (go routine/workers)
// ----------------------------------------------------------------------------

func managementTask(c *asyncConnHdl, ctl workerCtl) (sig *interrupt_code, te *taskStatus) {
	/* TODO:
	connected:
		clearChannels();
		startWokers();
	disconnected:
		?
	on_fault:
		disconnect();
		goto disconnected;
	on_exit:
		?
	*/
	//	log.Println("MGR: do task ...")
	select {
	case stat := <-c.feedback:
		// do the shutdown for now -- TODO: try reconnect
		if stat.event == faulted || stat.event == quit_processed {
			if stat.event == faulted {
				// REVU - pretty please TODO do the customized log
				log.Printf("<INFO> - %s (manager task) FAULT EVENT ", c)
			}
			// REVU - pretty please TODO do the customized log
			//			log.Printf("<INFO> %s - (manager task) SHUTTING DOWN ...", c)
			c.shutdown <- true

			// REVU - pretty please TODO do the customized log
			//			log.Printf("<INFO> %s - (manager task) RAISING SIGNAL STOP ...", c)
			go func() { c.reqProcCtl <- stop }()
			go func() { c.rspProcCtl <- stop }()
			go func() { c.heartbeatCtl <- stop }()
			go func() { c.managerCtl <- stop }()
		}
	case s := <-ctl:
		return &s, &ok_status
	}

	return nil, &ok_status
}

// Task:
// "One PING only" after receiving a timed tick per ConnectionsSpec period
// - can be interrupted while waiting on ticker
// - uses TryGet on PING response (1 sec as of now)
//
// KNOWN ISSUE:
// possible connection is OK but TryGet timesout and now we just log a warning
// but it is a good measure of latencies in the pipeline and perhaps could be used
// to autoconfigure the params to decrease latencies ... TODO

func heartbeatTask(c *asyncConnHdl, ctl workerCtl) (sig *interrupt_code, te *taskStatus) {
	var async AsyncConnection = c
	select {
	//	case <-NewTimer(ns1Sec * c.spec().heartbeat):
	case <-time.NewTimer(c.spec().heartbeat).C:
		response, e := async.QueueRequest(&PING, [][]byte{})
		if e != nil {
			return nil, &taskStatus{reqerr, e}
		}
		stat, re, timedout := response.future.(FutureBool).TryGet(1 * time.Second)
		if re != nil {
			log.Printf("ERROR: Heartbeat recieved error response on PING: %d\n", re)
			return nil, &taskStatus{error_, re}
		} else if timedout {
			log.Println("Warning: Heartbeat timeout on get PING response.")
		} else {
			// flytrap
			if stat != true {
				log.Println("<BUG> Heartbeat recieved false stat on PING while response error was nil")
				return nil, &taskStatus{error_, NewError(SYSTEM_ERR, "BUG false stat on PING w/out error")}
			}
		}
	case sig := <-ctl:
		return &sig, &ok_status
	}
	return nil, &ok_status
}

// Task:
// process one pending response at a time
// - can be interrupted while waiting on the pending responses queue
// - buffered reader takes care of minimizing network io
//
// KNOWN BUG:
// until we figure out what's the problem with read timeout, can not
// be interrupted if hanging on a read
func dbRspProcessingTask(c *asyncConnHdl, ctl workerCtl) (sig *interrupt_code, te *taskStatus) {

	var req asyncReqPtr
	select {
	case sig := <-ctl:
		// interrupted
		return &sig, &ok_status
	case req = <-c.pendingResps:
		// continue to process
	}

	// process response to asyncRequest
	reader := c.super.reader
	cmd := req.cmd

	resp, e3 := GetResponse(reader, cmd) // REVU - protocol modified to handle VIRTUALS
	if e3 != nil {
		// system error
		log.Println("<TEMP DEBUG> Request sent to faults chan on error in GetResponse: ", e3)
		req.stat = rcverr
		req.error = NewErrorWithCause(SYSTEM_ERR, "GetResponse os.Error", e3)
		c.faults <- req
		return nil, &taskStatus{rcverr, e3}
	}

	// if responsed processed was for cmd QUIT then signal the rest of the crew
	// REVU - ok, a bit hacky but it works.
	if cmd == &QUIT {
		c.feedback <- workerStatus{0, quit_processed, nil, nil}
		fakesig := pause
		c.isShutdown = true
		SetFutureResult(req.future, cmd, resp)
		return &fakesig, &ok_status
	}

	SetFutureResult(req.future, cmd, resp)
	return nil, &ok_status
}

func msgProcessingTask(c *asyncConnHdl, ctl workerCtl) (sig *interrupt_code, te *taskStatus) {
	var req asyncReqPtr
	select {
	case sig := <-ctl:
		// interrupted
		return &sig, &ok_status
	case req = <-c.pendingResps:
		c.pendingResps <- req
		// continue to process
	}
	reader := c.super.reader
	//	cmd := &SUBSCRIBE

	message, e3 := getPubSubResponse(reader, nil)
	if e3 != nil {
		// system error
		log.Println("<TEMP DEBUG> on error in msgProcessingTask: ", e3)
		req.stat = rcverr
		req.error = NewErrorWithCause(SYSTEM_ERR, "GetResponse os.Error", e3)
		c.faults <- req
		return nil, &taskStatus{rcverr, e3}
	}
	log.Printf("MSG IN: %s\n", message)
	//
	panic("msgProcessingTask not implemented")
}

// Pending further tests, this addresses bug in earlier version
// and can be interrupted

func reqProcessingTask(c *asyncConnHdl, ctl workerCtl) (ic *interrupt_code, ts *taskStatus) {

	var err error
	var errmsg string

	bytecnt := 0
	blen := 0
	bufsize := c.spec().wBufSize
	var sig interrupt_code

	select {
	case req := <-c.pendingReqs:
		blen, err := c.processAsyncRequest(req)
		if err != nil {
			errmsg = fmt.Sprintf("processAsyncRequest error in initial phase")
			goto proc_error
		}
		bytecnt += blen
	case sig := <-ctl:
		return &sig, &ok_status
	}

	for bytecnt < bufsize {
		select {
		case req := <-c.pendingReqs:
			blen, err = c.processAsyncRequest(req)
			if err != nil {
				errmsg = fmt.Sprintf("processAsyncRequest error in batch phase")
				goto proc_error
			}
			if blen > 0 {
				bytecnt += blen
			} else {
				bytecnt = 0
			}

		case sig = <-ctl:
			ic = &sig
			goto done

		default:
			goto done
		}
	}

done:
	c.writer.Flush()
	return ic, &ok_status

proc_error:
	log.Println(errmsg, err)
	return nil, &taskStatus{snderr, err}
}

// ----------------------------------------------------------------------------
// asyncConnHdl internal ops
// ----------------------------------------------------------------------------

func (c *asyncConnHdl) processAsyncRequest(req asyncReqPtr) (blen int, e error) {
	//	req := <-c.pendingReqs;
	req.id = c.nextId()
	blen = len(*req.outbuff)
	e = sendRequest(c.writer, *req.outbuff)
	if e == nil {
		req.outbuff = nil
		select {
		case c.pendingResps <- req:
		default:
			c.writer.Flush()
			c.pendingResps <- req
			blen = 0
		}
	} else {
		log.Println("<BUG> lazy programmer >> ERROR in processRequest goroutine -req requeued for now")
		// TODO: set stat on future & inform conn control and put it in faulted list
		c.pendingReqs <- req
	}
	return
}

// request id needs to be unique in context of associated connection
// only one goroutine calls this so no need to provide concurrency guards
func (c *asyncConnHdl) nextId() (id int64) {
	id = c.nextid
	c.nextid++
	return
}
