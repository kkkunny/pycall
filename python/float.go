package python

/*
#include <Python.h>
*/
import "C"

var (
	FloatType = newTypeObject(&C.PyFloat_Type)
)

// need release
func FloatFromDouble(v float64) *Object {
	return newObject(C.PyFloat_FromDouble(C.double(v)))
}

func (self *Object) FloatAsDouble() float64 {
	return float64(C.PyFloat_AsDouble(self.v))
}
