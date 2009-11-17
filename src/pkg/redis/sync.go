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
