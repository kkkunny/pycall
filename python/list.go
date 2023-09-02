package python

/*
#include <Python.h>
*/
import "C"

var (
	ListType = newTypeObject(&C.PyList_Type)
)

// need release
func ListNew(length int) *Object {
	return newObject(C.PyList_New(C.Py_ssize_t(length)))
}

func (self *Object) ListSetItem(index int, item *Object) bool {
	return C.PyList_SetItem(self.v, C.Py_ssize_t(index), item.v) == 0
}

func (self *Object) ListSize() int {
	return int(C.PyList_Size(self.v))
}

func (self *Object) ListGetItem(index int) *Object {
	return newObject(C.PyList_GetItem(self.v, C.Py_ssize_t(index)))
}
