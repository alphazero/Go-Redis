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
	"errors"
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

type ErrorCategory uint8

const (
	_ ErrorCategory = iota
	REDIS_ERR
	SYSTEM_ERR
)

func (ec ErrorCategory) String() (s string) {
	switch ec {
	case SYSTEM_ERR:
		s = "SYSTEM_ERROR"
	case REDIS_ERR:
		s = "REDIS_ERROR"
	}
	return
}

// Defines the interfce to get details of an Error.
//
type Error interface {
	// underlying error that caused this error - or nil
	Cause() error

	// support std. Go Error interface
	Error() string

	// category of error -- never nil
	Category() ErrorCategory

	// the error message - if a Redis error these are exactly
	// as provided by the server e.g. "ERR - invalid password"
	// possibly "" if SYSTEM-ERR
	Message() string

	// conv to String
	String() string

	// check if error is a Redis error e.g. "ERR - invalid password"
	// a convenience method
	IsRedisError() bool
}
type redisError struct {
	msg      string
	category ErrorCategory
	cause    error
}

func (e redisError) IsRedisError() bool            { return e.category == REDIS_ERR}
func (e redisError) Cause() error            { return e.cause }
func (e redisError) Error() string           { return e.String() }
func (e redisError) Category() ErrorCategory { return e.category }
func (e redisError) Message() string         { return e.msg }
func (e redisError) String() string {
	var errCat string
	var causeDetails string = ""
	switch e.category {
	case REDIS_ERR:
		errCat = "REDIS-ERROR"
	case SYSTEM_ERR:
		errCat = "SYSTEM-ERROR"
		if e.Cause() != nil {
			causeDetails = fmt.Sprintf(" [cause: %s]", e.Cause())
		}
	}
	return fmt.Sprintf("[go-redis|%d|%s]: %s %s", e.category, errCat, e.msg, causeDetails)
}

func _newRedisError(cat ErrorCategory, msg string) redisError {
	var e redisError
	e.msg = msg
	e.category = cat
	return e
}

// Creates an Error of REDIS_ERR category with the message.
//
func NewRedisError(msg string) Error {
	return _newRedisError(REDIS_ERR, msg)
}

// Creates an Error of REDIS_ERR category with the message.
//
func NewSystemError(msg string) Error {
	return _newRedisError(SYSTEM_ERR, msg)
}

// Creates an Error of specified category with the message.
//
func NewError(cat ErrorCategory, msg string) Error {
	return _newRedisError(cat, msg)
}

// Creates an Error of specified category with the message and cause.
//
func NewErrorWithCause(cat ErrorCategory, msg string, cause error) Error {
	e := _newRedisError(cat, msg)
	e.cause = cause
	return e
}


// utility function emits log if _debug flag /debug() is true
// Error is returned.
// usage:
//      foo, e := FooBar()
//      if e != nil {
//          return withError(e)
//      }
func withError(e Error) Error {
	if debug() {
		log.Println(e)
	}
	return e
}

// creates a new generic (Go) error
// and emits log if _debug flag /debug() is true
// Error is returned.
// usage:
//      v := FooBar()
//      if v != expected {
//          return withNewError("value v is unexpected")
//      }
func withNewError(m string) error {
	e := errors.New(m)
	if debug() {
		log.Println(e)
	}
	return e
}

// creates a new redis.Error of category SYSTEM_ERR
// and emits log if _debug flag /debug() is true
// Error is returned.
// usage:
//      _, e := SomeLibraryOrGoCall()
//      if e != nil {
//          return withNewError("value v is unexpected", e)
//      }
func withOsError(m string, cause error) error {
	return withNewError(fmt.Sprintf("%s [cause: %s]", m, cause))
}
