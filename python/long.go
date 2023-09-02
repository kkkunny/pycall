package python

/*
#include <Python.h>
*/
import "C"

var (
	LongType = newTypeObject(&C.PyLong_Type)
)

// need release
func LongFromLongLong(v int64) *Object {
	return newObject(C.PyLong_FromLongLong(C.longlong(v)))
}

// need release
func LongFromUnsignedLongLong(v uint64) *Object {
	return newObject(C.PyLong_FromUnsignedLongLong(C.ulonglong(v)))
}

func (self *Object) LongAsLongLong() int64 {
	return int64(C.PyLong_AsLongLong(self.v))
}

func (self *Object) LongAsUnsignedLongLong() uint64 {
	return uint64(C.PyLong_AsUnsignedLongLong(self.v))
}
