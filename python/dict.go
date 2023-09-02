package python

/*
#include <Python.h>
*/
import "C"

var (
	DictType = newTypeObject(&C.PyDict_Type)
)

// need release
func DictNew() *Object {
	return newObject(C.PyDict_New())
}

func (self *Object) DictSetItem(key, val *Object) bool {
	return C.PyDict_SetItem(self.v, key.v, val.v) == 0
}

// need release
func (self *Object) DictItems() *Object {
	return newObject(C.PyDict_Items(self.v))
}
