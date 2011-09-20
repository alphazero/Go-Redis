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
   TODO:  do me ..
*/
package redis

import (
	"os"
	"fmt"
	"log"
)
// ----------------------------------------------------------------------------
// ERRORS
// ----------------------------------------------------------------------------

// redis.Errors can be either user or system category.
//  
// REDIS_ERR are sent by the Redis server in response to user commands. (For example
// "operation against key holding the wrong value).  If a call to the client method
// has a non nil error with this cateogry, it means you did something wrong.  If you
// want to know the details, see Error.GetMessage() which returns the precise error
// sent from Redis.
//
// If error is of category SYSTEM_ERR, it either means there were errors during exeuction
// (which may be due to your system environment, networking, memory, etc.) OR, there is
// a bug in the package.   Unlike REDIS_ERR cateogry, SYSTEM_ERR category errors may have
// and underlying cause (os.Error) which can be obtained by Error.Cause() method.

type ErrorCategory uint8;
const (
	_ ErrorCategory = iota;
	REDIS_ERR;
	SYSTEM_ERR;
)

// Defines the interfce to get details of an Error.
// 
type Error interface {
	Cause() os.Error;
	Category() ErrorCategory;
	Message() string;
	String() string;
}
type redisError struct {
	msg 	  string;
	category  ErrorCategory;
	cause     os.Error;
}
func (e redisError) Cause() os.Error { return e.cause; }
func (e redisError) Category() ErrorCategory { return e.category; }
func (e redisError) Message() string { return e.msg; }
func (e redisError) String() string {
	var errCat string;
	var causeDetails string = "";
	switch e.category {
	case REDIS_ERR: 
		errCat = "REDIS-ERROR";
	case SYSTEM_ERR: 
		errCat = "SYSTEM-ERROR";
		if e.Cause() != nil { 
			causeDetails = fmt.Sprintf(" [cause: %s]", e.Cause()); 
		}
	}
	return fmt.Sprintf("[go-redis|%d|%s]: %s %s", e.category, errCat, e.msg, causeDetails);
}

// Creates an Error of REDIS_ERR category with the message.
//
func NewRedisError(msg string) Error {
	e := new (redisError);
	e.msg = msg;
	e.category = REDIS_ERR;
	return e;
}

// Creates an Error of specified category with the message.
//
func NewError(t ErrorCategory, msg string) Error {
	e := new (redisError);
	e.msg = msg;
	e.category = t;
	return e;
}

// Creates an Error of specified category with the message and cause.
//
func NewErrorWithCause(cat ErrorCategory, msg string, cause os.Error) Error {
	e := new (redisError);
	e.msg = msg;
	e.category = cat;
	e.cause = cause;
	return e;
}

// ----------------
// impl. utils
//
// a utility function for various components
//
func withError (e Error) Error {
	if debug() { log.Println(e); }
	return e;
}
func withNewError (m string) os.Error {
	e:= os.NewError(m);
	if debug() { log.Println(e); }
	return e;
}
func withOsError (m string, cause os.Error) os.Error {
	return withNewError(fmt.Sprintf("%s [cause: %s]", m, cause));
}
