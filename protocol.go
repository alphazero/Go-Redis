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
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
)

// ----------------------------------------------------------------------------
// Wire
// ----------------------------------------------------------------------------

// protocol's special bytes
const (
	CR_BYTE    byte = byte('\r')
	LF_BYTE         = byte('\n')
	SPACE_BYTE      = byte(' ')
	ERR_BYTE        = byte('-')
	OK_BYTE         = byte('+')
	COUNT_BYTE      = byte('*')
	SIZE_BYTE       = byte('$')
	NUM_BYTE        = byte(':')
	TRUE_BYTE       = byte('1')
)

type ctlbytes []byte

var CRLF ctlbytes = ctlbytes{CR_BYTE, LF_BYTE}
var WHITESPACE ctlbytes = ctlbytes{SPACE_BYTE}

// ----------------------------------------------------------------------------
// Services
// ----------------------------------------------------------------------------

// Creates the byte buffer that corresponds to the specified Command and
// provided command arguments.
//
// TODO: tedious but need to check for errors on all buffer writes ..
//
func CreateRequestBytes(cmd *Command, args [][]byte) ([]byte, error) {

	cmd_bytes := []byte(cmd.Code)

	buffer := bytes.NewBufferString("")
	buffer.WriteByte(COUNT_BYTE)
	buffer.Write([]byte(strconv.Itoa(len(args) + 1)))
	buffer.Write(CRLF)
	buffer.WriteByte(SIZE_BYTE)
	buffer.Write([]byte(strconv.Itoa(len(cmd_bytes))))
	buffer.Write(CRLF)
	buffer.Write(cmd_bytes)
	buffer.Write(CRLF)

	for _, s := range args {
		buffer.WriteByte(SIZE_BYTE)
		buffer.Write([]byte(strconv.Itoa(len(s))))
		buffer.Write(CRLF)
		buffer.Write(s)
		buffer.Write(CRLF)
	}

	return buffer.Bytes(), nil
}

// Creates a specific Future type for the given Redis command
// and returns it as a generic reference.
func CreateFuture(cmd *Command) (future interface{}) {
	switch cmd.RespType {
	case BOOLEAN:
		future = newFutureBool()
	case BULK:
		future = newFutureBytes()
	case MULTI_BULK:
		future = newFutureBytesArray()
	case NUMBER:
		future = newFutureInt64()
	case STATUS:
		future = newFutureBool()
	case STRING:
		future = newFutureString()
	case VIRTUAL:
		// REVU - treating virtual futures as FutureBools (always true)
		future = newFutureBool()
	}
	return
}

// Sets the type specific result value from the response for the future reference
// based on the command type.
func SetFutureResult(future interface{}, cmd *Command, r Response) {
	if r.IsError() {
		future.(FutureResult).onError(NewRedisError(r.GetMessage()))
	} else {
		switch cmd.RespType {
		case BOOLEAN:
			future.(FutureBool).set(r.GetBooleanValue())
		case BULK:
			future.(FutureBytes).set(r.GetBulkData())
		case MULTI_BULK:
			future.(FutureBytesArray).set(r.GetMultiBulkData())
		case NUMBER:
			future.(FutureInt64).set(r.GetNumberValue())
		case STATUS:
			future.(FutureBool).set(true)
		case STRING:
			future.(FutureString).set(r.GetStringValue())
		case VIRTUAL:
			// REVU - OK to treat virtual commands as FutureBool
			future.(FutureBool).set(true)
		}
	}
}

// ----------------------------------------------------------------------------
// request processing
// ----------------------------------------------------------------------------

// Either writes all the bytes or it fails and returns an error
//
func sendRequest(w io.Writer, data []byte) (e error) {
	loginfo := "sendRequest"
	if w == nil {
		return withNewError(fmt.Sprintf("<BUG> in %s(): nil Writer", loginfo))
	}

	n, e := w.Write(data)
	if e != nil {
		msg := fmt.Sprintf("%s(): connection Write wrote %d bytes only.", loginfo, n)
		return withNewError(msg)
	}

	// doc isn't too clear but the underlying netFD may return n<len(data) AND
	// e == nil, but that's precisely what we're checking.
	// presumably we can try sending the remaining bytes but that is precisely
	// what netFD.Write is doing (and it couldn't) so ...
	if n < len(data) {
		msg := fmt.Sprintf("%s(): connection Write wrote %d bytes only.", loginfo, n)
		return withNewError(msg)
	}
	return
}

// ----------------------------------------------------------------------------
// Response
// ----------------------------------------------------------------------------

type Response interface {
	IsError() bool
	GetMessage() string
	GetBooleanValue() bool
	GetNumberValue() int64
	GetStringValue() string
	GetBulkData() []byte
	GetMultiBulkData() [][]byte
}
type _response struct {
	isError       bool
	msg           string
	boolval       bool
	numval        int64
	stringval     string
	bulkdata      []byte
	multibulkdata [][]byte
}

func (r *_response) IsError() bool          { return r.isError }
func (r *_response) GetMessage() string     { return r.msg }
func (r *_response) GetBooleanValue() bool  { return r.boolval }
func (r *_response) GetNumberValue() int64  { return r.numval }
func (r *_response) GetStringValue() string { return r.stringval }
func (r *_response) GetBulkData() []byte    { return r.bulkdata }
func (r *_response) GetMultiBulkData() [][]byte {
	return r.multibulkdata
}

// ----------------------------------------------------------------------------
// response processing
// ----------------------------------------------------------------------------

// Gets the response to the command.
// Any errors (whether runtime or bugs) are returned as os.Error.
// The returned response (regardless of flavor) may have (application level)
// errors as sent from Redis server.
//
func GetResponse(reader *bufio.Reader, cmd *Command) (resp Response, err error) {

	defer func() {
		e := recover()
		if e != nil {
			err = e.(error)
		}
	}()

	buf := readToCRLF(reader)

	// Redis error
	if buf[0] == ERR_BYTE {
		resp = &_response{msg: string(buf[1:]), isError: true}
		return
	}

	switch cmd.RespType {
	case STATUS:
		resp = &_response{msg: string(buf[1:])}
		return
	case STRING:
		assertCtlByte(buf, OK_BYTE, "STRING")
		resp = &_response{stringval: string(buf[1:])}
		return
	case BOOLEAN:
		assertCtlByte(buf, NUM_BYTE, "BOOLEAN")
		resp = &_response{boolval: buf[1] == TRUE_BYTE}
		return
	case NUMBER:
		assertCtlByte(buf, NUM_BYTE, "NUMBER")
		n, e := strconv.ParseInt(string(buf[1:]), 10, 64)
		assertNotError(e, "in GetResponse - parse error in NUMBER response")
		resp = &_response{numval: n}
		return
	case VIRTUAL:
		resp = &_response{boolval: true}
		return
	case BULK:
		assertCtlByte(buf, SIZE_BYTE, "BULK")
		size, e := strconv.Atoi(string(buf[1:]))
		assertNotError(e, "in GetResponse - parse error in BULK size")
		resp = &_response{bulkdata: readBulkData(reader, size)}
		return
	case MULTI_BULK:
		assertCtlByte(buf, COUNT_BYTE, "MULTI_BULK")
		cnt, e := strconv.Atoi(string(buf[1:]))
		assertNotError(e, "in GetResponse - parse error in MULTIBULK cnt")
		resp = &_response{multibulkdata: readMultiBulkData(reader, cnt)}
		return
	}

	panic(fmt.Errorf("BUG - GetResponse - this should not have been reached"))
}

func assertCtlByte(buf []byte, b byte, info string) {
	if buf[0] != b {
		panic(fmt.Errorf("control byte for %s is not '%s' as expected - got '%s'", info, string(b), string(buf[0])))
	}
}

func assertNotError(e error, info string) {
	if e != nil {
		panic(fmt.Errorf("%s - error: %s", info, e))
	}
}

// ----------------------------------------------------------------------------
// PubSub Message
// ----------------------------------------------------------------------------

// REVU - export?
type PubSubMType int

// REVU - export?
const (
	SUBSCRIBE_ACK PubSubMType = iota
	UNSUBSCRIBE_ACK
	MESSAGE
)

//// REVU - necessary?
//// REVU - export?
//type PubSubMessage interface {
//	Subscription() string
//	Type() PubSubMType
//	Body() []byte
//	Info() int64
//}

// Conforms to the payload as received from wire.
// If mtype is mtype_message, then body will contain the message body, and
// sinfo will be -1.
// otherwise, it is expected that sinfo will contain subscription-info,
// e.g. number of subscribed channels, and data will be nil.
type message struct {
	mtype   PubSubMType
	channel string
	body    []byte
	sinfo   int64
}

func newMessage(channel string, body []byte) *message {
	m := message{}
	m.mtype = MESSAGE
	m.channel = channel
	m.body = body
	m.sinfo = -1
	return &m
}

func newPubSubAck(mtype PubSubMType, channel string, scnt int64) *message {
	m := message{}
	m.mtype = mtype
	m.channel = channel
	m.body = nil
	m.sinfo = scnt
	return &m
}

func newSubcribeAck(channel string, scnt int64) *message {
	return newPubSubAck(SUBSCRIBE_ACK, channel, scnt)
}

func newUnsubcribeAck(channel string, scnt int64) *message {
	return newPubSubAck(UNSUBSCRIBE_ACK, channel, scnt)
}

// ----------------------------------------------------------------------------
// PubSub message processing
// ----------------------------------------------------------------------------

func getPubSubResponse(r *bufio.Reader, cmd *Command) (msg *message, err error) {
	defer func() {
		e := recover()
		if e != nil {
			err = e.(error)
		}
	}()

	buf := readToCRLF(r)
	assertCtlByte(buf, COUNT_BYTE, "PubSub Sequence")

	num, e := strconv.ParseInt(string(buf[1:len(buf)]), 10, 64)
	assertNotError(e, "in getPubSubResponse - ParseInt")
	if num != 3 {
		panic(fmt.Errorf("<BUG> Expecting *3 for len in response - got %d", num))
	}

	header := readMultiBulkData(r, 2)

	mtype := string(header[0])
	subid := string(header[1])

	buf = readToCRLF(r)

	n, e := strconv.ParseInt(string(buf[1:]), 10, 64)
	assertNotError(e, "in getPubSubResponse - pubsub msg seq 3 line - number parse error")

	switch mtype {
	case "subscribe", "psubscribe":
		assertCtlByte(buf, NUM_BYTE, "UNSUBSCRIBE_ACK")
		msg = newPubSubAck(UNSUBSCRIBE_ACK, subid, n)
	case "unsubscribe", "punsubscribe":
		assertCtlByte(buf, NUM_BYTE, "SUBSCRIBE_ACK")
		msg = newPubSubAck(SUBSCRIBE_ACK, subid, n)
	case "message":
		assertCtlByte(buf, SIZE_BYTE, "MESSAGE")
		msg = newMessage(subid, readBulkData(r, int(n)))
	}

	return
}

// ----------------------------------------------------------------------------
// Protocol i/o
// ----------------------------------------------------------------------------

// reads all bytes upto CR-LF.  (Will eat those last two bytes)
// return the line []byte up to CR-LF
// error returned is NOT ("-ERR ...").  If there is a Redis error
// that is in the line buffer returned
//
// panics on error

func readToCRLF(r *bufio.Reader) []byte {
	//	var buf []byte
	buf, err := r.ReadBytes(CR_BYTE)
	if err != nil {
		panic(fmt.Errorf("readToCRLF - %s", err))
	}

	var b byte
	b, err = r.ReadByte()
	if err != nil {
		panic(fmt.Errorf("readToCRLF - %s", err))
	}
	if b != LF_BYTE {
		err = errors.New("<BUG> Expecting a Linefeed byte here!")
	}
	return buf[0 : len(buf)-1]
}

func readBulkData(r *bufio.Reader, n int) (data []byte) {
	if n >= 0 {
		buffsize := n + 2
		data = make([]byte, buffsize)
		if _, e := io.ReadFull(r, data); e != nil {
			panic(fmt.Errorf("readBulkData - ReadFull - %s", e))
		} else {
			if data[n] != CR_BYTE || data[n+1] != LF_BYTE {
				panic(fmt.Errorf("terminal was not CRLF as expected"))
			}
			data = data[:n]
		}
	}
	return
}

func readMultiBulkData(conn *bufio.Reader, num int) [][]byte {
	data := make([][]byte, num)
	for i := 0; i < num; i++ {
		buf := readToCRLF(conn)
		if buf[0] != SIZE_BYTE {
			panic(fmt.Errorf("readMultiBulkData - Expecting a SIZE_BYTE"))
		}

		size, e := strconv.Atoi(string(buf[1:]))
		if e != nil {
			panic(fmt.Errorf("readMultiBulkData - Atoi parse error - %s", e))
		}
		data[i] = readBulkData(conn, size)
	}
	return data
}
