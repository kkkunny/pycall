package python

/*
#include <stdlib.h>
#include <Python.h>
*/
import "C"
import "unsafe"

var (
	UnicodeType = newTypeObject(&C.PyUnicode_Type)
)

// need release
func UnicodeFromString(u string) *Object {
	cu := C.CString(u)
	defer C.free(unsafe.Pointer(cu))
	return newObject(C.PyUnicode_FromString(cu))
}

func (self *Object) UnicodeAsUTF8() string {
	cstr := C.PyUnicode_AsUTF8(self.v)
	return C.GoString(cstr)
}
