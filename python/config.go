package python

/*
#include <stdlib.h>
#include <Python.h>
#include "config.h"
*/
import "C"
import (
	"unsafe"
)

type Config struct {
	v *C.PyConfig
}

func ConfigNewConfig() *Config {
	cfg := &Config{v: C.PyConfig_NewConfig()}
	cfg.configInitPythonConfig()
	return cfg
}

func (self *Config) configInitPythonConfig() {
	C.PyConfig_InitPythonConfig(self.v)
}

func (self *Config) ConfigClear() {
	C.PyConfig_Clear(self.v)
	C.free(unsafe.Pointer(self.v))
}

func (self *Config) ConfigSetProgramName(path string) *Status {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	return newStatus(C.PyConfig_SetProgramName(self.v, cpath))
}

func (self *Config) ConfigAddPath(path string) *Status {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	return newStatus(C.PyConfig_AddPath(self.v, cpath))
}

func (self *Config) ConfigSetExecutable(path string) *Status {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	return newStatus(C.PyConfig_SetExecutable(self.v, cpath))
}
