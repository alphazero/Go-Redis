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
	Package reflect_utils provides helper functions to faciliate working with
	the reflect artifacts of Go.  
*/
package reflect_utils

import (
	"reflect";
)

func GetByteArrayAtIndex(v *reflect.StructValue, i int) (arr []byte, ok bool){
	field := v.Field(i);
	return GetByteArray (field);
}
func GetByteArray (v reflect.Value) (arr []byte, ok bool) {
	switch v := v.(type) {
	case reflect.ArrayOrSliceValue:
		aosv := reflect.ArrayOrSliceValue(v);
		arr = make([]byte, aosv.Len());
		for i:=0; i<aosv.Len();i++ {
//			a := aosv.Elem(i).(*reflect.Uint8Value);
			arr[i] = aosv.Elem(i).(*reflect.Uint8Value).Get(); 
		} 
		return arr, true;
	}
	return;
}
