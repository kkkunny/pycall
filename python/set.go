package python

/*
#include <Python.h>
*/
import "C"

var (
	SetType = newTypeObject(&C.PySet_Type)
)

// need release
func SetNew() *Object {
	return newObject(C.PySet_New((*C.PyObject)(C.NULL)))
}

func (self *Object) SetAdd(key *Object) bool {
	return C.PySet_Add(self.v, key.v) == 0
}

func (self *Object) SetSize() int {
	return int(C.PySet_Size(self.v))
}
