package python

/*
#include <stdlib.h>
#include <Python.h>
*/
import "C"
import "unsafe"

// need release
func RunString(str string, start int, globals, locals *Object) *Object {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	return newObject(C.PyRun_String(cstr, C.int(start), globals.v, locals.v))
}
