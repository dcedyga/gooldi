package concurrency

import (
	"reflect"
	"strings"
)

// StructHasField if an interface is either a struct or a pointer to a struct
// and has the defined member field
func HasField(Iface interface{}, FieldName string) bool {
	if Iface != nil {
		ValueIface := reflect.ValueOf(Iface)
	
		if ValueIface.Type().Kind() != reflect.Ptr {
			ValueIface = reflect.New(reflect.TypeOf(Iface))
		}
	
		Field := ValueIface.Elem().FieldByName(FieldName)
		if Field.IsValid() {
			return true
		}
	}
	
	return false
}

func SetFieldInt64Val(field string,input interface {},value int64){
	t := reflect.TypeOf(input)
	if (strings.HasPrefix(t.Name(), "*")){
		val := reflect.ValueOf(input).Elem().FieldByName(field)
		ptr := val.Addr().Interface().(*int64)
		*ptr = value
	}
}

func GetFieldInt64Val(field string,input interface {})int64{
	t := reflect.TypeOf(input)
	if (strings.HasPrefix(t.Name(), "*")){
		val := reflect.ValueOf(input).Elem().FieldByName(field)
		ptr := val.Addr().Interface().(*int64)
		return *ptr
	}
	return -1
}