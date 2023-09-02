package python

/*
#include <stdlib.h>
#include <Python.h>
*/
import "C"
import "unsafe"

// need release
func ImportImportModule(name string) *Object {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return newObject(C.PyImport_ImportModule(cname))
}
