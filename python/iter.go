package python

/*
#include <Python.h>
*/
import "C"

func (self *Object) IterNext() *Object {
	return newObject(C.PyIter_Next(self.v))
}
