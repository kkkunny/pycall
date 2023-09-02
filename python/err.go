package python

/*
#include <Python.h>
*/
import "C"

func ErrOccurred() *Object {
	return newObject(C.PyErr_Occurred())
}

func ErrClear() {
	C.PyErr_Clear()
}

func ErrFetch() (*Object, *Object, *Object) {
	var vptype, vpvalue, vptraceback *C.PyObject
	C.PyErr_Fetch(&vptype, &vpvalue, &vptraceback)
	return newObject(vptype), newObject(vpvalue), newObject(vptraceback)
}

func ErrNormalizeException(exc, val, tb *Object) {
	C.PyErr_NormalizeException(&exc.v, &val.v, &tb.v)
}
