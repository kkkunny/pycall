package python

/*
#include <Python.h>
*/
import "C"

type Status struct {
	v C.PyStatus
}

func newStatus(v C.PyStatus) *Status {
	return &Status{v: v}
}

func (self *Status) StatusException() bool {
	return C.PyStatus_Exception(self.v) != 0
}
