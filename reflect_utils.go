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
