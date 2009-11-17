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
	"time";
)

// ----------------------------------------------------------------------------
// synchronization utilities.
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


func NewTimer (ns int64) (signal <-chan int64) {
    if ns <= 0 {
        return nil
    }
    c := make(chan int64);
    go func() {
    	t := time.Nanoseconds();
    	e := time.Sleep(ns);
    	if e != nil { 
    		t = 0 - (time.Nanoseconds() - t);
    	}
    	else {
    		t = time.Nanoseconds() - t;
    	}
    	c<- t;
    }();
    return c;
}

// Signaling
//
// Signal interface defines the semantics of simple signaling between
// a sending and awaiting party, with timeout support.
//
type Signal interface {
	// Used to send the signal to the waiting party
	Send();
	
	// Used by the waiting party.  This call will block until
	// the Send() method has been invoked.
	Wait();
	
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
	WaitFor (timeout int64) (timedout bool, interrupted bool);
}

// signal wraps a channel and implements the Signal interface.
//
type signal struct {
	c chan byte;
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
//
//		tout, nsinterrupt := s.WaitFor (t);
//		if tout {
//			out.Printf("Timedout waiting for task.  interrupted: %v\n", nsinterrupt);
//
//			out.Printf("Will wait until the signal is sent ...\n");
//
//			// will block indefinitely until signal is sent
//			s.Wait();  
//
//			out.Printf("... alright - its done\n");
//		}
//		else {
//			out.Printf("Have signal task is completed!\n");
//		}
//	}
//


func NewSignal () Signal {
	c := make(chan byte);
	return &signal{c};
}

// implementation of Signal.Wait()
//
func (s *signal) Wait () {
	<-s.c;
	return;
}

// implementation of Signal.WaitFor(int64)
//
func (s *signal) WaitFor (timeout int64) (timedout bool, interrupted bool){
	timer := NewTimer(timeout);
	select {
		case <-s.c: 
		case to := <-timer:
			if to < 0 { interrupted = true; }
			else { timedout = true; }
	}
	return;
} 

// implementation of Signal.Send()
//
func (s *signal) Send () {
	s.c<-1;
}
