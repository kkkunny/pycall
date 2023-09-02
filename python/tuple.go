package python

/*
#include <Python.h>
*/
import "C"

var (
	TupleType = newTypeObject(&C.PyTuple_Type)
)

// need release
func TupleNew(length int) *Object {
	return newObject(C.PyTuple_New(C.Py_ssize_t(length)))
}

func (self *Object) TupleSetItem(pos int, o *Object) bool {
	return C.PyTuple_SetItem(self.v, C.Py_ssize_t(pos), o.v) == 0
}

func (self *Object) TupleSize() int {
	return int(C.PyTuple_Size(self.v))
}

func (self *Object) TupleGetItem(pos int) *Object {
	return newObject(C.PyTuple_GetItem(self.v, C.Py_ssize_t(pos)))
}

// need release
func (self *Object) TupleAsList() *Object {
	length := self.TupleSize()
	obj := ListNew(length)
	for i := 0; i < length; i++ {
		obj.ListSetItem(i, self.TupleGetItem(i))
	}
	return obj
}
