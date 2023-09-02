package python

// #cgo CFLAGS: -I/usr/include/python3.11
// #cgo LDFLAGS: -lpython3.11
/*
#include <stdlib.h>
#include <Python.h>
*/
import "C"

var (
	EvalInput   = int(C.Py_eval_input)
	FileInput   = int(C.Py_file_input)
	SingleInput = int(C.Py_single_input)
)

func Initialize() {
	C.Py_Initialize()
}

func InitializeFromConfig(cfg *Config) *Status {
	return newStatus(C.Py_InitializeFromConfig(cfg.v))
}

func IsInitialized() bool {
	return C.Py_IsInitialized() != 0
}

func FinalizeEx() bool {
	return C.Py_FinalizeEx() == 0
}

func GetVersion() string {
	return C.GoString(C.Py_GetVersion())
}

func DECREF(o *Object) {
	C.Py_DECREF(o.v)
}
