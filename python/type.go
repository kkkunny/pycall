package python

/*
#include <Python.h>
*/
import "C"

type TypeObject struct {
	v *C.PyTypeObject
}

func newTypeObject(v *C.PyTypeObject) *TypeObject {
	if v == nil {
		return nil
	}
	return &TypeObject{v: v}
}

func (self *TypeObject) Equal(t *TypeObject) bool {
	return self.v == t.v
}

// need release
func (self *TypeObject) TypeGetQualName() *Object {
	return newObject(C.PyType_GetQualName(self.v))
}

func (self *TypeObject) String() string {
	obj := self.TypeGetQualName()
	defer obj.Decref()
	return obj.UnicodeAsUTF8()
}
