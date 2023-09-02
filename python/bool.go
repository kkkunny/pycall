package python

/*
#include <Python.h>
*/
import "C"

var (
	BoolType = newTypeObject(&C.PyBool_Type)
	False    = newObject(C.Py_False)
	True     = newObject(C.Py_True)
)
