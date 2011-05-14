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
	"time"
	"log"
)

// ----------------------------------------------------------------------------
// synchronization utilities.
// ----------------------------------------------------------------------------

// result and result channel
//
// result defines a basic struct that either holds a generic (interface{}) reference
// or an Error reference. It is used to generically send and receive future results
// through channels.

type result struct {
	v interface{}
	e Error
}

// creates a result struct using the provided references
// and sends it on the specified channel.

func send(c chan result, v interface{}, e Error) {
	c <- result{v, e}
}

// blocks on the channel until a result is received.
// If the received result reference's e (error) field is not null,
// it will return it.

func receive(c chan result) (v interface{}, error Error) {
	fv := <-c
	if fv.e != nil {
		error = fv.e
	} else if fv.v == nil {
		// ? should we allow nil results?
	} else {
		v = fv.v
	}
	return
}
// using a timer blocks on the channel until a result is received.
// or timeout period expires.
// if timedout, returns ok==false.
// otherwise,
// If the received result reference's e (error) field is not null,
// it will return it.
//
// For now presumably a nil result is OK and if error is nil, the return
// value v is the intended result, even if nil.
// TODO: think this through a bit more

func tryReceive(c chan result, ns int64) (v interface{}, error Error, ok bool) {
	timer := NewTimer(ns)
	select {
	case fv := <-c:
		ok = true
		if fv.e != nil {
			error = fv.e
		} else if fv.v == nil {
			// ? should we allow nil results?
		} else {
			v = fv.v
		}
	case to := <-timer:
		if debug() {
			log.Println("resultchan.TryGet() -- timedout waiting for futurevaluechan | timeout after ", to)
		}
	}
	return
}

// ------------------
// Future? interfaces very much in line with the Future<?> of Java.
// These variants all expose the same set of semantics in a type-safe manner.
//
// We're only exposing the getters on these future objects as references to
// these interfaces are returned to redis users.  Same considerations also
// inform the decision to limit the exposure of the newFuture? methods to the
// package.
//
// Also note that while the current implementation does not enforce this, the
// Future? references can only be used until a value, or an error is obtained.
// If a value is obtained from a Future? reference, any further calls to Get()
// will block indefinitely.  It is OK, of course, to use TryGet(..) repeatedly
// until it returns with a true 'ok' return out param.

// FutureResult
//
// A generic future.  All type-safe Futures support this interface
//
type FutureResult interface {
	onError(Error)
}

// FutureBytes (for []byte)
//
type FutureBytes interface {
	//	onError (Error);
	set([]byte)
	Get() (vale []byte, error Error)
	TryGet(timeout int64) (value []byte, error Error, ok bool)
}
type _byteslicefuture chan result

func newFutureBytes() FutureBytes            { return make(_byteslicefuture, 1) }
func (fvc _byteslicefuture) onError(e Error) { send(fvc, nil, e) }
func (fvc _byteslicefuture) set(v []byte)    { send(fvc, v, nil) }
func (fvc _byteslicefuture) Get() (v []byte, error Error) {
	gv, err := receive(fvc)
	if err != nil {
		return nil, err
	}
	return gv.([]byte), err
}
func (fvc _byteslicefuture) TryGet(ns int64) (v []byte, error Error, ok bool) {
	gv, err, ok := tryReceive(fvc, ns)
	if !ok {
		return nil, nil, ok
	}
	if err != nil {
		return nil, err, ok
	}
	return gv.([]byte), err, ok
}

// FutureBytesArray (for [][]byte)
//
type FutureBytesArray interface {
	//	onError (Error);
	set([][]byte)
	Get() (vale [][]byte, error Error)
	TryGet(timeout int64) (value [][]byte, error Error, ok bool)
}
type _bytearrayslicefuture chan result

func newFutureBytesArray() FutureBytesArray { return make(_bytearrayslicefuture, 1) }
func (fvc _bytearrayslicefuture) onError(e Error) {
	send(fvc, nil, e)
}
func (fvc _bytearrayslicefuture) set(v [][]byte) {
	send(fvc, v, nil)
}
func (fvc _bytearrayslicefuture) Get() (v [][]byte, error Error) {
	gv, err := receive(fvc)
	if err != nil {
		return nil, err
	}
	return gv.([][]byte), err
}
func (fvc _bytearrayslicefuture) TryGet(ns int64) (v [][]byte, error Error, ok bool) {
	gv, err, ok := tryReceive(fvc, ns)
	if !ok {
		return nil, nil, ok
	}
	if err != nil {
		return nil, err, ok
	}
	return gv.([][]byte), err, ok
}

// FutureBool
//
type FutureBool interface {
	//	onError (Error);
	set(bool)
	Get() (val bool, error Error)
	TryGet(timeout int64) (value bool, error Error, ok bool)
}
type _boolfuture chan result

func newFutureBool() FutureBool         { return make(_boolfuture, 1) }
func (fvc _boolfuture) onError(e Error) { send(fvc, nil, e) }
func (fvc _boolfuture) set(v bool)      { send(fvc, v, nil) }
func (fvc _boolfuture) Get() (v bool, error Error) {
	gv, err := receive(fvc)
	if err != nil {
		return false, err
	}
	return gv.(bool), err
}
func (fvc _boolfuture) TryGet(ns int64) (v bool, error Error, ok bool) {
	gv, err, ok := tryReceive(fvc, ns)
	if !ok {
		return false, nil, ok
	}
	if err != nil {
		return false, err, ok
	}
	return gv.(bool), err, ok
}

// ------------------
// FutureString
//
type FutureString interface {
	//	onError (execErr Error);
	set(v string)
	Get() (string, Error)
	TryGet(timeout int64) (value string, error Error, ok bool)
}
type _futurestring chan result

func newFutureString() FutureString       { return make(_futurestring, 1) }
func (fvc _futurestring) onError(e Error) { send(fvc, nil, e) }
func (fvc _futurestring) set(v string)    { send(fvc, v, nil) }
func (fvc _futurestring) Get() (v string, error Error) {
	gv, err := receive(fvc)
	if err != nil {
		return "", err
	}
	return gv.(string), err
}
func (fvc _futurestring) TryGet(ns int64) (v string, error Error, ok bool) {
	gv, err, ok := tryReceive(fvc, ns)
	if !ok {
		return "", nil, ok
	}
	if err != nil {
		return "", err, ok
	}
	return gv.(string), err, ok
}

// ------------------
// FutureInt64
//
type FutureInt64 interface {
	//	onError (execErr Error);
	set(v int64)
	Get() (int64, Error)
	TryGet(timeout int64) (value int64, error Error, ok bool)
}
type _futureint64 chan result

func newFutureInt64() FutureInt64        { return make(_futureint64, 1) }
func (fvc _futureint64) onError(e Error) { send(fvc, nil, e) }
func (fvc _futureint64) set(v int64)     { send(fvc, v, nil) }
func (fvc _futureint64) Get() (v int64, error Error) {
	gv, err := receive(fvc)
	if err != nil {
		return -1, err
	}
	return gv.(int64), err
}
func (fvc _futureint64) TryGet(ns int64) (v int64, error Error, ok bool) {
	gv, err, ok := tryReceive(fvc, ns)
	if !ok {
		return -1, nil, ok
	}
	if err != nil {
		return -1, err, ok
	}
	return gv.(int64), err, ok
}
// ----------------------------------------------------------------------------
// Future types that wrap generic Redis response types - not entirely happy about
// this ...

// ------------------
// FutureFloat64
//
type FutureFloat64 interface {
	Get() (float64, Error)
	TryGet(timeout int64) (v float64, error Error, ok bool)
}
type _futurefloat64 struct {
	future FutureBytes
}

func newFutureFloat64(future FutureBytes) FutureFloat64 {
	return _futurefloat64{future}
}
func (fvc _futurefloat64) Get() (v float64, error Error) {
	gv, err := fvc.future.Get()
	if err != nil {
		return 0, err
	}
	v, err = Btof64(gv)
	return v, nil
}
func (fvc _futurefloat64) TryGet(ns int64) (v float64, error Error, ok bool) {
	gv, err, ok := fvc.future.TryGet(ns)
	if !ok {
		return 0, nil, ok
	}
	if err != nil {
		return 0, err, ok
	}
	v, err = Btof64(gv)
	return v, nil, ok
}


// ------------------
// FutureKeys
//
type FutureKeys interface {
	Get() ([]string, Error)
	TryGet(timeout int64) (keys []string, error Error, ok bool)
}
type _futurekeys struct {
	future FutureBytes
}

func newFutureKeys(future FutureBytes) FutureKeys {
	return _futurekeys{future}
}
func (fvc _futurekeys) Get() (v []string, error Error) {
	gv, err := fvc.future.Get()
	if err != nil {
		return nil, err
	}
	v = convAndSplit(gv)
	return v, nil
}
func (fvc _futurekeys) TryGet(ns int64) (v []string, error Error, ok bool) {
	gv, err, ok := fvc.future.TryGet(ns)
	if !ok {
		return nil, nil, ok
	}
	if err != nil {
		return nil, err, ok
	}
	v = convAndSplit(gv)
	return v, nil, ok
}

// ------------------
// FutureInfo
//
type FutureInfo interface {
	Get() (map[string]string, Error)
	TryGet(timeout int64) (keys map[string]string, error Error, ok bool)
}
type _futureinfo struct {
	future FutureBytes
}

func newFutureInfo(future FutureBytes) FutureInfo {
	return _futureinfo{future}
}
func (fvc _futureinfo) Get() (v map[string]string, error Error) {
	gv, err := fvc.future.Get()
	if err != nil {
		return nil, err
	}
	v = parseInfo(gv)
	return v, nil
}
func (fvc _futureinfo) TryGet(ns int64) (v map[string]string, error Error, ok bool) {
	gv, err, ok := fvc.future.TryGet(ns)
	if !ok {
		return nil, nil, ok
	}
	if err != nil {
		return nil, err, ok
	}
	v = parseInfo(gv)
	return v, nil, ok
}

// ------------------
// FutureKeyType
//
type FutureKeyType interface {
	Get() (KeyType, Error)
	TryGet(timeout int64) (keys KeyType, error Error, ok bool)
}
type _futurekeytype struct {
	future FutureString
}

func newFutureKeyType(future FutureString) FutureKeyType {
	return _futurekeytype{future}
}
func (fvc _futurekeytype) Get() (v KeyType, error Error) {
	gv, err := fvc.future.Get()
	if err != nil {
		return RT_NONE, err
	}
	v = GetKeyType(gv)
	return v, nil
}
func (fvc _futurekeytype) TryGet(ns int64) (v KeyType, error Error, ok bool) {
	gv, err, ok := fvc.future.TryGet(ns)
	if !ok {
		return RT_NONE, nil, ok
	}
	if err != nil {
		return RT_NONE, err, ok
	}
	v = GetKeyType(gv)
	return v, nil, ok
}

// ----------------------------------------------------------------------------
// Timer
//
// start a new timer that will signal on the returned
// channel when the specified ns (timeout in nanoseconds)
// have passsed.  If ns < 0, function returns immediately
// with nil.  Otherwise, the caller can select on the channel
// and will recieve an item after timeout.  If the timer
// itself was interrupted during sleep, the value in channel
// will be 0-time-elapsed.  Otherwise, for normal operation,
// it will return time elapsed in ns (which hopefully is very
// close to the specified ns.
//
// Example:
//
//	tasksignal := DoSomethingWhileIWait ();  // could take a while..
//
//	timeout := redis.NewTimer(1000*800);
//
//	select {
//		case <-tasksignal:
//			out.Printf("Task completed!\n");
//		case to := <-timeout:
//			out.Printf("Timedout waiting for task.  %d\n", to);
//	}


func NewTimer(ns int64) (signal <-chan int64) {
	if ns <= 0 {
		return nil
	}
	c := make(chan int64)
	go func() {
		t := time.Nanoseconds()
		e := time.Sleep(ns)
		if e != nil {
			t = 0 - (time.Nanoseconds() - t)
		} else {
			t = time.Nanoseconds() - t
		}
		c <- t
	}()
	return c
}

// ----------------------------------------------------------------------------
// Signaling
//
// Signal interface defines the semantics of simple signaling between
// a sending and awaiting party, with timeout support.
//
type Signal interface {
	// Used to send the signal to the waiting party
	Send()

	// Used by the waiting party.  This call will block until
	// the Send() method has been invoked.
	Wait()

	// Used by the waiting party.  This call will block until
	// either the Send() method has been invoked, or, an interrupt
	// occurs, or, the timeout duration passes.
	//
	// out param timedout is true if the period expired before
	// signal was received.
	//
	// out param interrupted is true if an interrupt occurred.
	//
	// timedout and interrupted are mutually exclusive.
	//
	WaitFor(timeout int64) (timedout bool, interrupted bool)
}

// signal wraps a channel and implements the Signal interface.
//
type signal struct {
	c chan byte
}

// Creates a new Signal
//
// Usage exmple:
//
//  The sending party -- here it also creates the signal but that
// can happen elsewhere and passed to it.
//
//	func DoSomethingAndSignalOnCompletion (ns int64) (redis.Signal) {
//		s := redis.NewSignal();
//   	go func () {
//			out.Printf("I'm going to sleep for %d nseconds ...\n", ns);
//			time.Sleep(ns);
//			out.Printf("the sleeper has awakened!\n");
//			s.Send();
//		}();
//		return s;
//	}
//
// elsewhere, the waiting party gets a signal (here by making a call to
// the above func) and then first waits using
//
//	func useSignal(t int64) {
//
//		// returns a signal
//		s := DoSomethingSignalOnCompletion(1000*1000);
//
//		// wait on signal or timeout
//		tout, nsinterrupt := s.WaitFor (t);
//		if tout {
//			out.Printf("Timedout waiting for task.  interrupted: %v\n", nsinterrupt);
//		}
//		else {
//			out.Printf("Have signal task is completed!\n");
//		}
//	}
//

func NewSignal() Signal {
	c := make(chan byte)
	return &signal{c}
}

// implementation of Signal.Wait()
//
func (s *signal) Wait() {
	<-s.c
	return
}

// implementation of Signal.WaitFor(int64)
//
func (s *signal) WaitFor(timeout int64) (timedout bool, interrupted bool) {
	timer := NewTimer(timeout)
	select {
	case <-s.c:
	case to := <-timer:
		if to < 0 {
			interrupted = true
		} else {
			timedout = true
		}
	}
	return
}

// implementation of Signal.Send()
//
func (s *signal) Send() { s.c <- 1 }
