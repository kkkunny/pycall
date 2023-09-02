package python

/*
#include <stdlib.h>
#include <Python.h>
*/
import "C"
import "unsafe"

type Object struct {
	v *C.PyObject
}

func newObject(v *C.PyObject) *Object {
	if v == nil || C.Py_IsNone(v) == 1 {
		return nil
	}
	return &Object{v: v}
}

// need release
func ObjectNew(t *TypeObject) *Object {
	return newObject(C._PyObject_New(t.v))
}

func (self *Object) Type() *TypeObject {
	return newTypeObject(C.Py_TYPE(self.v))
}

// need release
func (self *Object) ObjectCallObject(args ...*Object) *Object {
	if len(args) == 0 {
		return newObject(C.PyObject_CallObject(self.v, (*C.PyObject)(C.NULL)))
	}
	argObj := TupleNew(len(args))
	for i, o := range args {
		argObj.TupleSetItem(i, o)
	}
	return newObject(C.PyObject_CallObject(self.v, argObj.v))
}

// need release
func (self *Object) ObjectGetAttrString(attrName string) *Object {
	cattrName := C.CString(attrName)
	defer C.free(unsafe.Pointer(cattrName))
	return newObject(C.PyObject_GetAttrString(self.v, cattrName))
}

func (self *Object) ObjectTypeCheck(t *TypeObject) bool {
	return C.PyObject_TypeCheck(self.v, t.v) != 0
}

// need release
func (self *Object) ObjectStr() *Object {
	return newObject(C.PyObject_Str(self.v))
}

func (self *Object) String() string {
	obj := self.ObjectStr()
	defer obj.Decref()
	return obj.UnicodeAsUTF8()
}

func (self *Object) ObjectIsTrue() bool {
	return C.PyObject_IsTrue(self.v) == 1
}

func (self *Object) Decref() {
	DECREF(self)
}

// need release
func (self *Object) ObjectGetIter() *Object {
	return newObject(C.PyObject_GetIter(self.v))
}
