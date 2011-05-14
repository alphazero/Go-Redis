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
	"reflect"
)

// ----------------------------------------------------------------------------
//	reflect contains a few helper functions to facilitate working with
//	the reflect artifacts of Go.  Probably deserves to be in its own
//	(independent) little package as using Go reflection is quite tedious.
// ----------------------------------------------------------------------------

func ToByteSliceArray(v reflect.Value) (bsa [][]byte, ok bool) {
	n := v.NumField()
	bsa = make([][]byte, n)
	for i := 0; i < n; i++ {
		bsa[i], ok = GetByteArrayAtIndex(v, i)
		if !ok {
			if debug() {
			}
			break
		}
	}
	return
}
// TODO: document me
//
func GetByteArrayAtIndex(v reflect.Value, i int) (arr []byte, ok bool) {
	field := v.Field(i)
	return GetByteArray(field)
}

// TODO: document me
//
func GetByteArray(v reflect.Value) (arr []byte, ok bool) {
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		aosv := reflect.Value(v)
		arr = make([]byte, aosv.Len())
		for i := 0; i < aosv.Len(); i++ {
			arr[i] = uint8(aosv.Index(i).Uint())
		}
		return arr, true
	}
	return
}

// TODO: document me
//
func GetByteArrayLen(v reflect.Value) (len int, ok bool) {
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		aosv := reflect.Value(v)
		return aosv.Len(), true
	}
	return
}
