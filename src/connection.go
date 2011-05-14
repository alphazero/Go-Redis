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
	"net"
	"fmt"
	"os"
	"io"
	"bufio"
	"log"
	//	"time";
)

const (
	TCP       = "tcp"
	LOCALHOST = "127.0.0.1"
	ns1MSec   = 1000000
	ns1Sec    = ns1MSec * 1000
)

// various default sizes for the connections
// exported for user convenience if nedded
const (
	DefaultReqChanSize  = 1000000
	DefaultRespChanSize = 1000000

	DefaultTCPReadBuffSize      = 1024 * 256
	DefaultTCPWriteBuffSize     = 1024 * 256
	DefaultTCPReadTimeoutNSecs  = ns1Sec * 10
	DefaultTCPWriteTimeoutNSecs = ns1Sec * 10
	DefaultTCPLinger            = 0
	DefaultTCPKeepalive         = true
	DefaultHeartbeatSecs        = 1
)

// Redis specific default settings
// exported for user convenience if nedded
const (
	DefaultRedisPassword = ""
	DefaultRedisDB       = 0
	DefaultRedisPort     = 6379
	DefaultRedisHost     = LOCALHOST
)
// ----------------------------------------------------------------------------
// Connection ConnectionSpec
// ----------------------------------------------------------------------------

// Defines the set of parameters that are used by the client connections
//
type ConnectionSpec struct {
	host     string
	port     int
	password string
	db       int
	// tcp specific specs
	rBufSize   int
	wBufSize   int
	rTimeout   int64
	wTimeout   int64
	keepalive  bool
	lingerspec int // -n: finish io; 0: discard, +n: wait for n secs to finish
	// async specs
	reqChanCap int
	rspChanCap int
	//
	heartbeat int64 // 0 means no heartbeat
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
func (spec *ConnectionSpec) Heartbeat(seconds int64) *ConnectionSpec {
	spec.heartbeat = seconds
	return spec
}

// ----------------------------------------------------------------------------
// Generic Conn handle and methods
// ----------------------------------------------------------------------------

// General control structure used by connections.
//
type connHdl struct {
	spec   *ConnectionSpec
	conn   net.Conn // may want to change this to TCPConn
	reader *bufio.Reader
}

// Creates and opens a new connection to server per ConnectionSpec.
// The new connection is wrapped by a new connHdl with its bufio.Reader
// delegating to the net.Conn's reader.
//
func newConnHdl(spec *ConnectionSpec) (hdl *connHdl, err Error) {
	here := "newConnHdl"

	if hdl = new(connHdl); hdl == nil {
		return nil, NewError(SYSTEM_ERR, fmt.Sprintf("%s(): failed to allocate connHdl", here))
	}
	addr := fmt.Sprintf("%s:%d", spec.host, spec.port)
	raddr, e := net.ResolveTCPAddr("tcp", addr)
	if e != nil {
		msg := fmt.Sprintf("%s(): failed to resolve remote address %s", here, addr)
		return nil, NewErrorWithCause(SYSTEM_ERR, msg, e)
	}
	conn, e := net.DialTCP(TCP, nil, raddr)
	switch {
	case e != nil:
		err = NewErrorWithCause(SYSTEM_ERR, fmt.Sprintf("%s(): could not open connection", here), e)
	case conn == nil:
		err = NewError(SYSTEM_ERR, fmt.Sprintf("%s(): net.Dial returned nil, nil (?)", here))
	default:
		configureConn(conn, spec)
		hdl.spec = spec
		hdl.conn = conn
		bufsize := 4096
		hdl.reader, e = bufio.NewReaderSize(conn, bufsize)
		if e != nil {
			msg := fmt.Sprintf("%s(): bufio.NewReaderSize (%d) error", here, bufsize)
			err = NewErrorWithCause(SYSTEM_ERR, msg, e)
		} else {
			if err = hdl.onConnect(); err == nil && debug() {
				fmt.Println("[Go-Redis] Opened SynchConnection connection to ", addr)
			}
		}
	}
	return hdl, err
}

func configureConn(conn *net.TCPConn, spec *ConnectionSpec) {
	// these two -- the most important -- are causing problems on my osx/64
	// where a "service unavailable" pops up in the async reads
	// but we absolutely need to be able to use timeouts.
	//			conn.SetReadTimeout(spec.rTimeout);
	//			conn.SetWriteTimeout(spec.wTimeout);
	conn.SetLinger(spec.lingerspec)
	conn.SetKeepAlive(spec.keepalive)
	conn.SetReadBuffer(spec.rBufSize)
	conn.SetWriteBuffer(spec.wBufSize)
}
// TODO: return redis.Error
func (c *connHdl) onConnect() (e Error) {
	if c.spec.password != DefaultRedisPassword {
		_, e = c.ServiceRequest(&AUTH, [][]byte{[]byte(c.spec.password)})
		if e != nil {
			return
		}
	}
	if c.spec.db != DefaultRedisDB {
		_, e = c.ServiceRequest(&SELECT, [][]byte{[]byte(fmt.Sprintf("%d", c.spec.db))})
		if e != nil {
			return
		}
	}
	return
}

func (c *connHdl) onDisconnect() Error {
	return nil // for now
}
// closes the connHdl's net.Conn connection.
// Is public so that connHdl struct can be used as SyncConnection (TODO: review that.)
//
func (hdl connHdl) Close() os.Error {
	err := hdl.conn.Close()
	if debug() {
		fmt.Println("[Go-Redis] Closed connection: ", hdl)
	}
	return err
}

// ----------------------------------------------------------------------------
// Connection SyncConnection
// ----------------------------------------------------------------------------

// Defines the service contract supported by synchronous (Request/Reply)
// connections.

type SyncConnection interface {
	ServiceRequest(cmd *Command, args [][]byte) (Response, Error)
	Close() os.Error
}

// Creates a new SyncConnection using the provided ConnectionSpec
func NewSyncConnection(spec *ConnectionSpec) (c SyncConnection, err Error) {
	//	connHdl, e := newConnHdl(spec);
	//	connHdl.onConnect();
	//	return connHdl, e;
	return newConnHdl(spec)
}

// Implementation of SyncConnection.ServiceRequest.
//
func (chdl *connHdl) ServiceRequest(cmd *Command, args [][]byte) (resp Response, err Error) {
	here := "connHdl.ServiceRequest"
	errmsg := ""
	ok := false
	buff, e := CreateRequestBytes(cmd, args) // 2<<<
	if e == nil {
		e = sendRequest(chdl.conn, buff)
		if e == nil {
			resp, e = GetResponse(chdl.reader, cmd)
			if e == nil {
				if resp.IsError() {
					redismsg := fmt.Sprintf(" [%s]: %s", cmd.Code, resp.GetMessage())
					err = NewRedisError(redismsg)
				}
				ok = true
			} else {
				errmsg = fmt.Sprintf("%s(%s): failed to get response", here, cmd.Code)
			}
		} else {
			errmsg = fmt.Sprintf("%s(%s): failed to send request", here, cmd.Code)
		}
	} else {
		errmsg = fmt.Sprintf("%s(%s): failed to create request buffer", here, cmd.Code)
	}

	if !ok {
		return resp, withError(NewErrorWithCause(SYSTEM_ERR, errmsg, e)) // log it on debug
	}

	return
}

// ----------------------------------------------------------------------------
// Asynchronous connections
// ----------------------------------------------------------------------------

type status_code byte

const (
	ok status_code = iota
	info
	warning
	error
	reqerr
	inierr
	snderr
	rcverr
)

type interrupt_code byte

const (
	_ interrupt_code = iota
	// connection process control
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
)

type workerStatus struct {
	id       int
	event    event_code
	taskinfo *taskStatus
	ctlchan  *workerCtl
}
type taskStatus struct {
	code  status_code
	error os.Error
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

	shutdown chan bool
}

// Creates a new asyncConnHdl with a new connHdl as its delegated 'super'.
// Note it does not start the processing goroutines for the channels.

func newAsyncConnHdl(spec *ConnectionSpec) (async *asyncConnHdl, err Error) {
	//	here := "newAsynConnHDL";
	connHdl, err := newConnHdl(spec)
	if err == nil && connHdl != nil {
		async = new(asyncConnHdl)
		if async != nil {
			async.super = connHdl
			var e os.Error
			async.writer, e = bufio.NewWriterSize(connHdl.conn, spec.wBufSize)
			if e == nil {
				async.pendingReqs = make(chan asyncReqPtr, spec.reqChanCap)
				async.pendingResps = make(chan asyncReqPtr, spec.rspChanCap)
				async.faults = make(chan asyncReqPtr, spec.reqChanCap) // not sure about sizing here ...

				async.reqProcCtl = make(workerCtl)
				async.rspProcCtl = make(workerCtl)
				async.heartbeatCtl = make(workerCtl)
				async.managerCtl = make(workerCtl)

				async.feedback = make(chan workerStatus)
				async.shutdown = make(chan bool, 1)

				return
			} else {
				connHdl.conn.Close()
				msg := fmt.Sprintf("NewWriterSize(%d) failed.", spec.wBufSize)
				err = NewErrorWithCause(SYSTEM_ERR, msg, e)
			}
		}
	}
	// fall through here on errors only
	if debug() {
		log.Println("Error creating asyncConnHdl: ", err)
		//		err =  os.NewError("Error creating asyncConnHdl");
	}
	return nil, err
}

// Creates and opens a new AsyncConnection and starts the goroutines for
// request and response processing
// TODO: NewXConnection methods need to return redis.Error due to initial connect
// interaction with redis (AUTH &| SELECT)
func NewAsynchConnection(spec *ConnectionSpec) (conn AsyncConnection, err Error) {
	var async *asyncConnHdl
	if async, err = newAsyncConnHdl(spec); err == nil {
		async.onConnect()
		async.startup()
	}
	return async, err
}

func (c *asyncConnHdl) onConnect() (e Error) { return }
func (c *asyncConnHdl) onDisconnect() (e Error) {
	return
}

// responsible for managing the various moving parts of the asyncConnHdl
func (c *asyncConnHdl) startup() {
	//	fmt.Println("connection manager -- begin");

	go c.worker(manager, "manager", managementTask, c.managerCtl, nil)
	c.managerCtl <- start

	go c.worker(heartbeatworker, "heartbeat", heartbeatTask, c.heartbeatCtl, c.feedback)
	c.heartbeatCtl <- start

	go c.worker(heartbeatworker, "request-processor", reqProcessingTask, c.reqProcCtl, c.feedback)
	c.reqProcCtl <- start

	go c.worker(heartbeatworker, "response-processor", rspProcessingTask, c.rspProcCtl, c.feedback)
	c.rspProcCtl <- start

	fmt.Println("Connection started ...")
}

// This could find a happy home in a generalized worker package ...
// TODO
func (c *asyncConnHdl) worker(id int, name string, task workerTask, ctl workerCtl, fb chan workerStatus) {
	var signal interrupt_code
	var tstat *taskStatus

	// todo: add startup hook for worker

await_signal:
	fmt.Println(name, "_worker: await_signal.")
	signal = <-ctl

on_interrupt:
	//		fmt.Println(name, "_worker: on_interrupt: ", signal);
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
			fmt.Println(name, "_worker: task error!")
			tstat = stat
			goto on_error
		} else if is != nil {
			signal = *is
			goto on_interrupt
		}
		goto work
	}

on_error:
	fmt.Println(name, "_worker: on_error!")
	// TODO: log it, send it, and go back to wait_start:
	log.Println(name, "_worker task raised error: ", tstat)
	fb <- workerStatus{id, faulted, tstat, &ctl}
	goto await_signal

before_stop:
	fmt.Println(name, "_worker: before_stop!")
	// TODO: add shutdown hook for worker

	fmt.Println(name, "_worker: STOPPED!")
}

// ----------------------------------------------------------------------------
// aync tasks
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
	log.Println("MGR: do task ...")
	select {
	case stat := <-c.feedback:
		log.Println("MGR: Feedback from one of my minions: ", stat)
		// do the shutdown for now -- TODO: try reconnect
		if stat.event == faulted {
			log.Println("MGR: Shutting down due to fault in ", stat.id)
			go func() { c.reqProcCtl <- stop }()
			go func() { c.rspProcCtl <- stop }()
			go func() { c.heartbeatCtl <- stop }()

			log.Println("MGR: Signal SHUTDOWN ... ")
			c.shutdown <- true
			// stop self // TODO: should manager really be a task or a FSM?
			c.managerCtl <- stop
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
	case <-NewTimer(ns1Sec * c.super.spec.heartbeat):
		response, e := async.QueueRequest(&PING, [][]byte{})
		if e != nil {
			return nil, &taskStatus{reqerr, e}
		}
		stat, re, ok := response.future.(FutureBool).TryGet(ns1Sec)
		if re != nil {
			log.Println("ERROR: Heartbeat recieved error response on PING")
			return nil, &taskStatus{error, re}
		} else if !ok {
			log.Println("Warning: Heartbeat timeout on get PING response.")
		} else {
			// flytrap
			if stat != true {
				log.Println("<BUG> Heartbeat recieved false stat on PING while response error was nil")
				return nil, &taskStatus{error, NewError(SYSTEM_ERR, "BUG false stat on PING w/out error")}
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
func rspProcessingTask(c *asyncConnHdl, ctl workerCtl) (sig *interrupt_code, te *taskStatus) {

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

	resp, e3 := GetResponse(reader, cmd)
	if e3 != nil {
		// system error
		log.Println("<TEMP DEBUG> Request sent to faults chan on error in GetResponse: ", e3)
		req.stat = rcverr
		req.error = NewErrorWithCause(SYSTEM_ERR, "GetResponse os.Error", e3)
		c.faults <- req
		return nil, &taskStatus{rcverr, e3}
	}
	SetFutureResult(req.future, cmd, resp)

	return nil, &ok_status
}

// Pending further tests, this addresses bug in earlier version
// and can be interrupted

func reqProcessingTask(c *asyncConnHdl, ctl workerCtl) (ic *interrupt_code, ts *taskStatus) {

	var err os.Error
	var errmsg string

	bytecnt := 0
	blen := 0
	bufsize := c.super.spec.wBufSize
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

func (c *asyncConnHdl) processAsyncRequest(req asyncReqPtr) (blen int, e os.Error) {
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
		// TODO: set stat on future & inform conn control and put it in fauls
		c.pendingReqs <- req
	}
	return
}


// ----------------------------------------------------------------------------
// AsyncConnection interface & Impl
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

func (c *asyncConnHdl) QueueRequest(cmd *Command, args [][]byte) (*PendingResponse, Error) {
	select {
	case <-c.shutdown:
		log.Println("<DEBUG> we're shutdown and not accepting any more requests ...")
		return nil, NewError(SYSTEM_ERR, "Connection is shutdown.")
	default:
	}

	future := CreateFuture(cmd)

	request := &asyncRequestInfo{0, 0, cmd, nil, future, nil}
	buff, e1 := CreateRequestBytes(cmd, args)
	if e1 == nil {
		request.outbuff = &buff
		c.pendingReqs <- request
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
// internal ops
// ----------------------------------------------------------------------------

// request id needs to be unique in context of associated connection
// only one goroutine calls this so no need to provide concurrency guards
func (c *asyncConnHdl) nextId() (id int64) {
	id = c.nextid
	c.nextid++
	return
}

// Either writes all the bytes or it fails and returns an error
//
func sendRequest(w io.Writer, data []byte) (e os.Error) {
	here := "connHdl.sendRequest"
	if w == nil {
		return withNewError(fmt.Sprintf("<BUG> in %s(): nil Writer", here))
	}

	n, e := w.Write(data)
	if e != nil {
		var msg string
		switch {
		case e == os.EAGAIN:
			// socket timeout -- don't handle that yet but may in future ..
			msg = fmt.Sprintf("%s(): timeout (os.EAGAIN) error on Write", here)
		default:
			// anything else
			msg = fmt.Sprintf("%s(): error on Write", here)
		}
		return withOsError(msg, e)
	}

	// doc isn't too clear but the underlying netFD may return n<len(data) AND
	// e == nil, but that's precisely what we're testing for here.
	// presumably we can try sending the remaining bytes but that is precisely
	// what netFD.Write is doing (and it couldn't) so ...
	if n < len(data) {
		msg := fmt.Sprintf("%s(): connection Write wrote %d bytes only.", here, n)
		return withNewError(msg)
	}
	return
}
