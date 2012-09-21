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

type PubSubConnection interface {
	Subscriptions() map[string]*Subscription
	ServiceRequest(cmd *Command, args [][]byte) (ok bool, err Error)
}

type Subscription struct {
	activated chan bool
	Channel   chan []byte
	IsActive  bool
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
// panics on error (with error)
func newConnHdl(spec *ConnectionSpec) (hdl *connHdl) {
	loginfo := "newConnHdl"

	hdl = new(connHdl)
	// REVU - this is silly
	if hdl == nil {
		panic(fmt.Errorf("%s(): failed to allocate connHdl", loginfo))
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
			panic(fmt.Errorf("%s(): failed to resolve remote address %s", loginfo, addr))
		}
	}

	conn, e := net.Dial(mode, addr)
	switch {
	case e != nil:
		panic(fmt.Errorf("%s(): could not open connection", loginfo))
	case conn == nil:
		panic(fmt.Errorf("%s(): net.Dial returned nil, nil (?)", loginfo))
	default:
		configureConn(conn, spec)
		hdl.spec = spec
		hdl.conn = conn
		hdl.connected = true
		bufsize := 4096
		hdl.reader = bufio.NewReaderSize(conn, bufsize)
	}
	return
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
// panics on error (with error)
func (c *connHdl) connect() {
	if c.spec.password != DefaultRedisPassword {
		args := [][]byte{[]byte(c.spec.password)}
		if _, e := c.ServiceRequest(&AUTH, args); e != nil {
			panic(fmt.Errorf("<ERROR> Authentication failed - %s", e.Message()))
		}
	}
	if c.spec.db != DefaultRedisDB {
		args := [][]byte{[]byte(fmt.Sprintf("%d", c.spec.db))}
		if _, e := c.ServiceRequest(&SELECT, args); e != nil {
			panic(fmt.Errorf("<ERROR> REDIS_DB Select failed - %s", e.Message()))
		}
	}
	// REVU - pretty please TODO do the customized log
	//	log.Printf("<INFO> %s - CONNECTED", c)
	return
}

// disconnects from net connections and sets connected state to false
// panics on net error (with error)
func (hdl *connHdl) disconnect() {
	// silently ignore repeated calls to closed connections
	if hdl.connected {
		if e := hdl.conn.Close(); e != nil {
			panic(fmt.Errorf("on connHdl.Close()", e))
			//			return NewErrorWithCause(SYSTEM_ERR, "on connHdl.Close()", e)
		}
		hdl.connected = false
		// REVU - pretty please TODO do the customized log
		//		log.Printf("<INFO> %s - DISCONNECTED", hdl)
	}
}

// Creates a new SyncConnection using the provided ConnectionSpec.
// Note that this function will also connect to the specified redis server.
func NewSyncConnection(spec *ConnectionSpec) (c SyncConnection, err Error) {
	defer func() {
		if e := recover(); e != nil {
			connerr := e.(error)
			err = NewErrorWithCause(SYSTEM_ERR, "NewSyncConnection", connerr)
		}
	}()

	connHdl := newConnHdl(spec)
	connHdl.connect()
	c = connHdl
	return
}

// Implementation of SyncConnection.ServiceRequest.
func (c *connHdl) ServiceRequest(cmd *Command, args [][]byte) (resp Response, err Error) {
	loginfo := "connHdl.ServiceRequest"

	defer func() {
		if re := recover(); re != nil {
			// REVU - needs to be logged - TODO
			err = NewErrorWithCause(SYSTEM_ERR, "ServiceRequest", re.(error))
		}
	}()

	if !c.connected {
		panic(fmt.Errorf("Connection %s is alredy closed", c.String()))
	}

	if cmd == &QUIT {
		c.disconnect()
		return
	}

	buff := CreateRequestBytes(cmd, args)
	sendRequest(c.conn, buff) // panics

	// REVU - this demands resp to be non-nil even in case of io errors
	// TODO - look into this
	resp, e := GetResponse(c.reader, cmd)
	if e != nil {
		panic(fmt.Errorf("%s(%s) - failed to get response", loginfo, cmd.Code))
	}

	// handle Redis server ERR - don't panic
	if resp.IsError() {
		redismsg := fmt.Sprintf(" [%s]: %s", cmd.Code, resp.GetMessage())
		err = NewRedisError(redismsg)
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

	nextid int64

	pendingReqs  chan asyncReqPtr
	pendingResps chan asyncReqPtr
	faults       chan asyncReqPtr

	subscriptions map[string]*Subscription // REDIS_PUBSUB only

	managerCtl   workerCtl
	reqProcCtl   workerCtl
	rspProcCtl   workerCtl
	heartbeatCtl workerCtl // REDIS_DB only

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
//
// REVU - PubSub could be checked here
func newAsyncConnHdl(spec *ConnectionSpec) *asyncConnHdl {

	c := new(asyncConnHdl)

	// connection base
	connHdl := newConnHdl(spec) // panics
	c.super = connHdl

	// conn management
	c.managerCtl = make(workerCtl)
	c.feedback = make(chan workerStatus)
	c.shutdown = make(chan bool, 1)

	// request processing
	c.writer = bufio.NewWriterSize(connHdl.conn, spec.wBufSize)
	c.reqProcCtl = make(workerCtl)
	c.pendingReqs = make(chan asyncReqPtr, spec.reqChanCap) // REVU for PubSub

	// fault processing
	c.faults = make(chan asyncReqPtr, spec.reqChanCap) // REVU - not sure about sizing

	// response processing
	c.rspProcCtl = make(workerCtl)

	switch spec.protocol {
	case REDIS_DB:
		c.heartbeatCtl = make(workerCtl)
		c.pendingResps = make(chan asyncReqPtr, spec.rspChanCap)
	case REDIS_PUBSUB:
		c.subscriptions = make(map[string]*Subscription)
	}

	// REVU - this is state - TODO move to startup
	c.isShutdown = false

	return c
}

// Creates and opens a new AsyncConnection and starts the goroutines for
// request and response processing
// interaction with redis (AUTH &| SELECT)
func NewAsynchConnection(spec *ConnectionSpec) (conn AsyncConnection, err Error) {
	defer func() {
		if e := recover(); e != nil {
			connerr := e.(error)
			err = NewErrorWithCause(SYSTEM_ERR, "NewSyncConnection", connerr)
		}
	}()

	//	var async *asyncConnHdl
	async := newAsyncConnHdl(spec)
	async.connect()
	async.startup()

	conn = async
	return
}

func NewPubSubConnection(spec *ConnectionSpec) (conn PubSubConnection, err Error) {
	defer func() {
		if e := recover(); e != nil {
			connerr := e.(error)
			err = NewErrorWithCause(SYSTEM_ERR, "NewSyncConnection", connerr)
		}
	}()

	//	var async *asyncConnHdl
	spec.Protocol(REDIS_PUBSUB) // must be so set it regardless
	async := newAsyncConnHdl(spec)
	async.connect()
	async.startup()

	conn = async
	return
}

// ----------------------------------------------------------------------------
// asyncConnHdl support for AsyncConnection interface
// ----------------------------------------------------------------------------

func (c *asyncConnHdl) QueueRequest(cmd *Command, args [][]byte) (pending *PendingResponse, err Error) {

	defer func() {
		if re := recover(); re != nil {
			// REVU - needs to be logged - TODO
			err = NewErrorWithCause(SYSTEM_ERR, "QueueRequest", re.(error))
		}
	}()

	if c.isShutdown {
		panic(fmt.Errorf("Connection %s is alredy shutdown", c.String()))
	}
	select {
	case <-c.shutdown:
		c.isShutdown = true
		c.shutdown <- true // put it back REVU likely to be a bug under heavy load ..
		panic(fmt.Errorf("Connection %s is alredy shutdown", c.String()))
	default:
	}

	buff := CreateRequestBytes(cmd, args) // panics
	future := CreateFuture(cmd)
	request := &asyncRequestInfo{0, 0, cmd, &buff, future, nil}

	c.pendingReqs <- request

	pending = &PendingResponse{future}

	return
}

// ----------------------------------------------------------------------------
// asyncConnHdl support for PubSubConnection interface
// ----------------------------------------------------------------------------

// TODO - HERE
// PubSubConnection support (only)
// Accepts Redis commands (P)SUBSCRIBE and (P)UNSUBSCRIBE.
// Request is processed asynchronously but call semantics are sync/blocking.
func (c *asyncConnHdl) ServiceRequest(cmd *Command, args [][]byte) (ok bool, err Error) {

	defer func() {
		if re := recover(); re != nil {
			// REVU - needs to be logged - TODO
			err = NewErrorWithCause(SYSTEM_ERR, "QueueRequest", re.(error))
		}
	}()
	switch *cmd {
	case SUBSCRIBE, UNSUBSCRIBE: /* nop - ok */
	default:
		panic(fmt.Errorf("BUG - command %s is not applicable to PubSub", cmd))
	}

	if c.isShutdown {
		panic(fmt.Errorf("Connection %s is alredy shutdown", c.String()))
	}
	select {
	case <-c.shutdown:
		c.isShutdown = true
		c.shutdown <- true // put it back REVU likely to be a bug under heavy load ..
		panic(fmt.Errorf("Connection %s is alredy shutdown", c.String()))
	default:
	}

	buff := CreateRequestBytes(cmd, args) // panics
	for _, arg := range args {
		topic := string(arg)
		if s := c.subscriptions[topic]; s != nil {
			panic(fmt.Errorf("already subscribed to topic %s", topic))
		}
		subscription := &Subscription{
			IsActive:  false,
			activated: make(chan bool, 1),
			Channel:   make(chan []byte, 100), // TODO - from spec
		}
		c.subscriptions[topic] = subscription
		//  TODO - go routine for each sub to wait on the ack
		//  TODO - and we wait here for all
	}

	//	future := CreateFuture(cmd)
	//	request := &asyncRequestInfo{0, 0, cmd, &buff, future, nil}
	request := &asyncRequestInfo{0, 0, cmd, &buff, nil, nil}
	c.pendingReqs <- request

	// REVU - we can do timeout and mod the sig.
	ok = true

	return
}

func (c *asyncConnHdl) Subscriptions() map[string]*Subscription {
	return c.subscriptions
}

// ----------------------------------------------------------------------------
// asyncConnHdl internal ops
// ----------------------------------------------------------------------------

// Delegates to connHdl - suepr panics.
// See connHdl#connect
func (c *asyncConnHdl) connect() {
	c.super.connect()
}

// REVU - TODO opt 2 for Quit here
// panics
func (c *asyncConnHdl) disconnect() {

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

	go c.worker(requesthandler, "request-processor", reqProcessingTask, c.reqProcCtl, c.feedback)
	c.reqProcCtl <- start

	var rspProcTask workerTask
	switch protocol {
	case REDIS_DB:
		rspProcTask = dbRspProcessingTask
	case REDIS_PUBSUB:
		rspProcTask = msgProcessingTask
	}
	go c.worker(responsehandler, "response-processor", rspProcTask, c.rspProcCtl, c.feedback)
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

// REVU - 	this is wrong
//			Given that there is no symmetric REQ for RESPs (unlike db reqs/rsps)
// 			this needs to be a read on net with timeout so we need to mod
//			or enhance the protocol funcs to deal with net read timeout.
func msgProcessingTask(c *asyncConnHdl, ctl workerCtl) (sig *interrupt_code, te *taskStatus) {
	log.Println("<TEMP DEBUG> msgProcessingTask - S0 ")
	log.Printf("<TEMP DEBUG> msgProcessingTask - c.pendingResps:%s\n", c.pendingResps)
	var req asyncReqPtr
	select {
	case sig := <-ctl:
		// interrupted
		return &sig, &ok_status
	case req = <-c.pendingResps:
		c.pendingResps <- req
		// continue to process
	}

	log.Println("<TEMP DEBUG> msgProcessingTask - S1 ")
	reader := c.super.reader
	//	cmd := &SUBSCRIBE

	message, e := GetPubSubResponse(reader)
	if e != nil {
		// system error
		log.Println("<TEMP DEBUG> on error in msgProcessingTask: ", e)
		req.stat = rcverr
		req.error = NewErrorWithCause(SYSTEM_ERR, "GetResponse os.Error", e)
		c.faults <- req
		return nil, &taskStatus{rcverr, e}
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

// REVU - error return on this internal func is OK - see call site usage.
func (c *asyncConnHdl) processAsyncRequest(req asyncReqPtr) (blen int, e error) {
	//	req := <-c.pendingReqs;
	req.id = c.nextId()
	blen = len(*req.outbuff)

	defer func() {
		if re := recover(); re != nil {
			e = re.(error)
			log.Println("<BUG> lazy programmer >> ERROR in processRequest goroutine -req requeued for now")
			// TODO: set stat on future & inform conn control and put it in faulted list
			c.pendingReqs <- req
		}
	}()
	sendRequest(c.writer, *req.outbuff)

	req.outbuff = nil
	select {
	case c.pendingResps <- req:
	default:
		c.writer.Flush()
		c.pendingResps <- req
		blen = 0
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
