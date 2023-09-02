#include <stdlib.h>
#include <locale.h>
#include <string.h>
#include <Python.h>

wchar_t* string2wstring(const char* s){
    size_t size = strlen(s)+1;
    wchar_t* ws = malloc(sizeof(wchar_t)*size);
    mbstowcs(ws, s, size);
    return ws;
}

PyConfig* PyConfig_NewConfig(){
    return (PyConfig*)malloc(sizeof(PyConfig));
}

PyStatus PyConfig_SetProgramName(PyConfig* config, const char* path){
    wchar_t* wpath = string2wstring(path);
    PyStatus status = PyConfig_SetString(config, &config->program_name, wpath);
    free(wpath);
    return status;
}

PyStatus PyConfig_AddPath(PyConfig* config, const char* path){
    wchar_t* wpath = string2wstring(path);
    config->module_search_paths_set = 1;
    PyStatus status = PyWideStringList_Append(&config->module_search_paths, wpath);
    free(wpath);
    return status;
}

PyStatus PyConfig_SetExecutable(PyConfig* config, const char* path){
    wchar_t* wpath = string2wstring(path);
    PyStatus status = PyConfig_SetString(config, &config->executable, wpath);
    free(wpath);
    return status;
}