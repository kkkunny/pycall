PyConfig* PyConfig_NewConfig();
PyStatus PyConfig_SetProgramName(PyConfig* config, const char* path);
PyStatus PyConfig_AddPath(PyConfig* config, const char* path);
PyStatus PyConfig_SetExecutable(PyConfig* config, const char* path);