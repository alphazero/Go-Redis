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
package main

import (
	"flag"
	"fmt"
	"log"
	"redis"
)

// Test Setex command
func main() {

	// Parse command-line flags; needed to let flags used by Go-Redis be parsed.
	flag.Parse()

	spec := redis.DefaultSpec().Db(13).Password("go-redis")
	client, e := redis.NewSynchClientWithSpec(spec)
	if e != nil {
		log.Println("failed to create the client", e)
		return
	}

	key := "expire:test"

	if err := client.Setex(key, 120, []byte("dummy-key-value")); err != nil {
		log.Println("error on Setex", err)
		return
	}

	value, e := client.Get(key)
	if e != nil {
		log.Println("error on Get", e)
		return
	}

	ttl, e := client.Ttl(key)
	if e != nil {
		log.Println("error on Ttl", e)
		return
	}

	fmt.Printf("%s with value %s expires in %d seconds", key, value, ttl)

	return
}
