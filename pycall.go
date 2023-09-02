package pycall

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strings"

	"github.com/kkkunny/pycall/python"
)

func Initialize(bin string, program string, path ...string) error {
	cfg := python.ConfigNewConfig()
	defer cfg.ConfigClear()

	if cfg.ConfigSetProgramName(program).StatusException() {
		return fmt.Errorf("set python program name error")
	}
	for _, p := range path {
		if cfg.ConfigAddPath(p).StatusException() {
			return fmt.Errorf("add python path error")
		}
	}
	if cfg.ConfigSetExecutable(bin).StatusException() {
		return fmt.Errorf("set python sxecutable error")
	}
	if python.InitializeFromConfig(cfg).StatusException() {
		return fmt.Errorf("initialize python error")
	}
	return nil
}

func Finalize() error {
	if !python.FinalizeEx() {
		return fmt.Errorf("finalize python error")
	}
	return nil
}

// GetDefaultExecutable 获取默认可执行文件
func GetDefaultExecutable() (string, error) {
	path, err := exec.LookPath("python3")
	if err == nil {
		return path, nil
	}
	if !errors.Is(err, exec.ErrNotFound) {
		return "", err
	}
	return exec.LookPath("python")
}

// GetModuleSearchPaths 获取包搜索路径
func GetModuleSearchPaths(bin string) ([]string, error) {
	cmd := exec.Command(bin, "-c", "import sys; print(sys.path)")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	re, err := regexp.Compile(`'(.*?)'`)
	if err != nil {
		return nil, err
	}
	reResult := re.FindAllStringSubmatch(string(output), -1)
	if len(reResult) == 0 {
		return nil, nil
	}
	paths := make([]string, 0, len(reResult))
	for _, r := range reResult {
		if strings.TrimSpace(r[1]) == "" {
			continue
		}
		paths = append(paths, r[1])
	}
	return paths, err
}

func InitializeDefault(path ...string) error {
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	bin, err := GetDefaultExecutable()
	if err != nil {
		return err
	}
	paths, err := GetModuleSearchPaths(bin)
	if err != nil {
		return err
	}
	paths = append(paths, root)
	return Initialize(bin, root, append(paths, path...)...)
}

// Recover 从异常中恢复
func Recover() error {
	defer python.ErrClear()
	vptype, vpvalue, vptraceback := python.ErrFetch()
	if vptype == nil || vpvalue == nil || vptraceback == nil {
		return nil
	}
	python.ErrNormalizeException(vptype, vpvalue, vptraceback)
	return fmt.Errorf(vpvalue.String())
}

// Go2PythonType Go类型转Python类型
func Go2PythonType(t reflect.Type) (*python.TypeObject, error) {
	switch t.Kind() {
	case reflect.Bool:
		return python.BoolType, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return python.LongType, nil
	case reflect.Float32, reflect.Float64:
		return python.FloatType, nil
	case reflect.String:
		return python.UnicodeType, nil
	case reflect.Slice:
		return python.ListType, nil
	case reflect.Map:
		if t.Elem().AssignableTo(reflect.StructOf(nil)) {
			return python.SetType, nil
		}
		return python.DictType, nil
	case reflect.Struct:
		return python.TupleType, nil
	default:
		return nil, fmt.Errorf("can not covert type '%s' in golang to type in python", t)
	}
}

// Go2PythonObject Go值转Python值
func Go2PythonObject(v any, pt *python.TypeObject) (*python.Object, error) {
	vt := reflect.TypeOf(v)

	if pt == nil {
		var err error
		pt, err = Go2PythonType(vt)
		if err != nil {
			return nil, err
		}
	}

	switch {
	case pt.Equal(python.BoolType):
		if vt.Kind() != reflect.Bool {
			return nil, fmt.Errorf("can not covert type '%s' in golang to type '%s' in python", vt, python.BoolType)
		}
		vv := v.(bool)
		if vv {
			return python.True, nil
		} else {
			return python.False, nil
		}
	case pt.Equal(python.LongType):
		switch vt.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			vv := reflect.ValueOf(v).Int()
			return python.LongFromLongLong(vv), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			vv := reflect.ValueOf(v).Uint()
			return python.LongFromUnsignedLongLong(vv), nil
		default:
			return nil, fmt.Errorf("can not covert type '%s' in golang to type '%s' in python", vt, python.LongType)
		}
	case pt.Equal(python.FloatType):
		switch vt.Kind() {
		case reflect.Float32, reflect.Float64:
			vv := reflect.ValueOf(v).Float()
			return python.FloatFromDouble(vv), nil
		default:
			return nil, fmt.Errorf("can not covert type '%s' in golang to type '%s' in python", vt, python.FloatType)
		}
	case pt.Equal(python.UnicodeType):
		if vt.Kind() != reflect.String {
			return nil, fmt.Errorf("can not covert type '%s' in golang to type '%s' in python", vt, python.UnicodeType)
		}
		vv := v.(string)
		return python.UnicodeFromString(vv), nil
	case pt.Equal(python.ListType):
		if vt.Kind() != reflect.Slice {
			return nil, fmt.Errorf("can not covert type '%s' in golang to type '%s' in python", vt, python.ListType)
		}
		vv := reflect.ValueOf(v)
		length := vv.Len()
		obj := python.ListNew(length)
		for i := 0; i < length; i++ {
			elemObj, err := Go2PythonObject(vv.Index(i).Interface(), nil)
			if err != nil {
				return nil, err
			}
			defer elemObj.Decref()
			if !obj.ListSetItem(i, elemObj) {
				return nil, fmt.Errorf("python list set error")
			}
		}
		return obj, nil
	case pt.Equal(python.DictType):
		if vt.Kind() != reflect.Map {
			return nil, fmt.Errorf("can not covert type '%s' in golang to type '%s' in python", vt, python.DictType)
		}
		vv := reflect.ValueOf(v)
		iter := vv.MapRange()
		obj := python.DictNew()
		for iter.Next() {
			keyObj, err := Go2PythonObject(iter.Key().Interface(), nil)
			if err != nil {
				return nil, err
			}
			defer keyObj.Decref()
			valObj, err := Go2PythonObject(iter.Value().Interface(), nil)
			if err != nil {
				return nil, err
			}
			defer valObj.Decref()
			if !obj.DictSetItem(keyObj, valObj) {
				return nil, fmt.Errorf("python dict set error")
			}
		}
		return obj, nil
	case pt.Equal(python.SetType):
		if vt.Kind() != reflect.Map || !vt.Elem().AssignableTo(reflect.StructOf(nil)) {
			return nil, fmt.Errorf("can not covert type '%s' in golang to type '%s' in python", vt, python.SetType)
		}
		vv := reflect.ValueOf(v)
		keys := vv.MapKeys()
		obj := python.SetNew()
		for _, key := range keys {
			elemObj, err := Go2PythonObject(key.Interface(), nil)
			if err != nil {
				return nil, err
			}
			defer elemObj.Decref()
			if !obj.SetAdd(elemObj) {
				return nil, fmt.Errorf("python set add error")
			}
		}
		return obj, nil
	case pt.Equal(python.TupleType):
		if vt.Kind() != reflect.Struct {
			return nil, fmt.Errorf("can not covert type '%s' in golang to type '%s' in python", vt, python.TupleType)
		}
		vv := reflect.ValueOf(v)
		length := vv.NumField()
		elemObjs := make([]*python.Object, 0, length)
		for i := 0; i < length; i++ {
			ft, fv := vt.Field(i), vv.Field(i)
			if !ft.IsExported() {
				continue
			}
			elemObj, err := Go2PythonObject(fv.Interface(), nil)
			if err != nil {
				return nil, err
			}
			defer elemObj.Decref()
			elemObjs = append(elemObjs, elemObj)
		}
		obj := python.TupleNew(len(elemObjs))
		for i, elemObj := range elemObjs {
			if !obj.TupleSetItem(i, elemObj) {
				return nil, fmt.Errorf("python tuple set error")
			}
		}
		return obj, nil
	default:
		return nil, fmt.Errorf("can not covert type '%s' in golang to type '%s' in python", vt, pt)
	}
}

// Python2GoInterface Python值转Go值
func Python2GoInterface(v *python.Object, gt reflect.Type) (any, error) {
	if v == nil {
		return nil, nil
	}
	vt := v.Type()

	switch {
	case vt.Equal(python.BoolType):
		if gt.Kind() != reflect.Bool {
			return nil, fmt.Errorf("can not covert type '%s' in python to type '%s' in golang", python.BoolType, gt)
		}
		return v.ObjectIsTrue(), nil
	case vt.Equal(python.LongType):
		switch gt.Kind() {
		case reflect.Int:
			return int(v.LongAsLongLong()), nil
		case reflect.Int8:
			return int8(v.LongAsLongLong()), nil
		case reflect.Int16:
			return int16(v.LongAsLongLong()), nil
		case reflect.Int32:
			return int32(v.LongAsLongLong()), nil
		case reflect.Int64:
			return v.LongAsLongLong(), nil
		case reflect.Uint:
			return uint(v.LongAsUnsignedLongLong()), nil
		case reflect.Uint8:
			return uint8(v.LongAsUnsignedLongLong()), nil
		case reflect.Uint16:
			return uint16(v.LongAsUnsignedLongLong()), nil
		case reflect.Uint32:
			return uint32(v.LongAsUnsignedLongLong()), nil
		case reflect.Uint64:
			return v.LongAsUnsignedLongLong(), nil
		default:
			return nil, fmt.Errorf("can not covert type '%s' in python to type '%s' in golang", python.BoolType, gt)
		}
	case vt.Equal(python.FloatType):
		switch gt.Kind() {
		case reflect.Float32:
			return float32(v.FloatAsDouble()), nil
		case reflect.Float64:
			return v.FloatAsDouble(), nil
		default:
			return nil, fmt.Errorf("can not covert type '%s' in python to type '%s' in golang", python.FloatType, gt)
		}
	case vt.Equal(python.UnicodeType):
		if gt.Kind() != reflect.String {
			return nil, fmt.Errorf("can not covert type '%s' in python to type '%s' in golang", python.UnicodeType, gt)
		}
		return v.UnicodeAsUTF8(), nil
	case vt.Equal(python.ListType):
		if gt.Kind() != reflect.Slice {
			return nil, fmt.Errorf("can not covert type '%s' in python to type '%s' in golang", python.ListType, gt)
		}
		et := gt.Elem()
		length := v.ListSize()
		iter := v.ObjectGetIter()
		defer iter.Decref()
		obj := reflect.MakeSlice(gt, length, length)
		f := func(i int, item *python.Object) error {
			defer item.Decref()
			elem, err := Python2GoInterface(item, et)
			if err != nil {
				return err
			}
			obj.Index(i).Set(reflect.ValueOf(elem))
			return nil
		}
		var i int
		for item := iter.IterNext(); item != nil; item = iter.IterNext() {
			err := f(i, item)
			if err != nil {
				return nil, err
			}
			i++
		}
		return obj.Interface(), nil
	case vt.Equal(python.DictType):
		if gt.Kind() != reflect.Map {
			return nil, fmt.Errorf("can not covert type '%s' in python to type '%s' in golang", python.DictType, gt)
		}
		kt, vt := gt.Key(), gt.Elem()
		iter := v.DictItems().ObjectGetIter()
		defer iter.Decref()
		obj := reflect.MakeMap(gt)
		f := func(item *python.Object) error {
			defer item.Decref()
			key, err := Python2GoInterface(item.TupleGetItem(0), kt)
			if err != nil {
				return err
			}
			value, err := Python2GoInterface(item.TupleGetItem(1), vt)
			if err != nil {
				return err
			}
			obj.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value))
			return nil
		}
		for item := iter.IterNext(); item != nil; item = iter.IterNext() {
			err := f(item)
			if err != nil {
				return nil, err
			}
		}
		return obj.Interface(), nil
	case vt.Equal(python.SetType):
		if gt.Kind() != reflect.Map || !gt.Elem().AssignableTo(reflect.StructOf(nil)) {
			return nil, fmt.Errorf("can not covert type '%s' in python to type '%s' in golang", python.SetType, gt)
		}
		et := gt.Key()
		iter := v.ObjectGetIter()
		defer iter.Decref()
		obj := reflect.MakeMap(gt)
		f := func(item *python.Object) error {
			defer item.Decref()
			elem, err := Python2GoInterface(item, et)
			if err != nil {
				return err
			}
			obj.SetMapIndex(reflect.ValueOf(elem), reflect.ValueOf(struct{}{}))
			return nil
		}
		for item := iter.IterNext(); item != nil; item = iter.IterNext() {
			err := f(item)
			if err != nil {
				return nil, err
			}
		}
		return obj.Interface(), nil
	case vt.Equal(python.TupleType):
		if gt.Kind() != reflect.Struct {
			return nil, fmt.Errorf("can not covert type '%s' in python to type '%s' in golang", python.TupleType, gt)
		}

		length := gt.NumField()
		elemTypes := make([]reflect.StructField, 0, length)
		for i := 0; i < length; i++ {
			ft := gt.Field(i)
			if !ft.IsExported() {
				continue
			}
			elemTypes = append(elemTypes, ft)
		}

		if len(elemTypes) != v.TupleSize() {
			return nil, fmt.Errorf("can not covert type '%s' in python to type '%s' in golang", python.TupleType, gt)
		}

		iter := v.ObjectGetIter()
		defer iter.Decref()
		obj := reflect.New(gt)
		f := func(i int, item *python.Object) error {
			defer item.Decref()
			field := elemTypes[i]
			elem, err := Python2GoInterface(item, field.Type)
			if err != nil {
				return err
			}
			obj.Elem().FieldByName(field.Name).Set(reflect.ValueOf(elem))
			return nil
		}
		var i int
		for item := iter.IterNext(); item != nil; item = iter.IterNext() {
			err := f(i, item)
			if err != nil {
				return nil, err
			}
			i++
		}
		return obj.Elem().Interface(), nil
	case vt.Equal(python.TupleType):
		if gt.Kind() != reflect.Struct {
			return nil, fmt.Errorf("can not covert type '%s' in python to type '%s' in golang", python.TupleType, gt)
		}

		length := gt.NumField()
		elemTypes := make([]reflect.StructField, 0, length)
		for i := 0; i < length; i++ {
			ft := gt.Field(i)
			if !ft.IsExported() {
				continue
			}
			elemTypes = append(elemTypes, ft)
		}

		if len(elemTypes) != v.TupleSize() {
			return nil, fmt.Errorf("can not covert type '%s' in python to type '%s' in golang", python.TupleType, gt)
		}

		iter := v.ObjectGetIter()
		defer iter.Decref()
		obj := reflect.New(gt)
		f := func(i int, item *python.Object) error {
			defer item.Decref()
			field := elemTypes[i]
			elem, err := Python2GoInterface(item, field.Type)
			if err != nil {
				return err
			}
			obj.Elem().FieldByName(field.Name).Set(reflect.ValueOf(elem))
			return nil
		}
		var i int
		for item := iter.IterNext(); item != nil; item = iter.IterNext() {
			err := f(i, item)
			if err != nil {
				return nil, err
			}
			i++
		}
		return obj.Elem().Interface(), nil
	default:
		return nil, fmt.Errorf("can not covert type '%s' in python to type '%s' in golang", vt, gt)
	}
}

// Python2GoValue Python值转Go值
func Python2GoValue[T any](v *python.Object) (T, error) {
	var empty T
	value, err := Python2GoInterface(v, reflect.TypeOf(empty))
	if err != nil {
		return empty, err
	}
	return value.(T), nil
}

// GetFunction 获取函数
// need release
func GetFunction[T any](module string, name string) (T, error) {
	var empty T
	ft := reflect.TypeOf(empty)
	if ft.Kind() != reflect.Func {
		return empty, fmt.Errorf("expect a function type for generic T")
	}

	moduleObj := python.ImportImportModule(module)
	if moduleObj == nil {
		return empty, fmt.Errorf("not exist module '%s'", module)
	}
	defer moduleObj.Decref()
	funcObj := moduleObj.ObjectGetAttrString(name)
	if funcObj == nil {
		return empty, fmt.Errorf("not exist function '%s'", name)
	}
	defer funcObj.Decref()

	f := reflect.MakeFunc(ft, func(args []reflect.Value) (results []reflect.Value) {
		retNum := ft.NumOut()
		results = make([]reflect.Value, retNum)
		for i := 0; i < retNum; i++ {
			results[i] = reflect.Zero(ft.Out(i))
		}

		defer func() {
			if err := Recover(); err != nil {
				results[retNum-1] = reflect.ValueOf(err)
			}
		}()

		pythonArgs := make([]*python.Object, len(args))
		for i, a := range args {
			op, err := Go2PythonObject(a.Interface(), nil)
			if err != nil {
				results[retNum-1] = reflect.ValueOf(err)
				return results
			}
			pythonArgs[i] = op
		}
		retObj := funcObj.ObjectCallObject(pythonArgs...)
		if retObj == nil {
			return results
		}
		defer retObj.Decref()

		if !retObj.ObjectTypeCheck(python.TupleType) {
			if retNum != 2 {
				results[retNum-1] = reflect.ValueOf(fmt.Errorf("param error"))
				return results
			}
			ret, err := Python2GoInterface(retObj, ft.Out(0))
			if err != nil {
				results[1] = reflect.ValueOf(err)
				return results
			}
			results[0] = reflect.ValueOf(ret)
			return results
		}
		fieldType := make([]reflect.StructField, retNum-1)
		for i := 0; i < retNum-1; i++ {
			fieldType[i] = reflect.StructField{
				Name: fmt.Sprintf("R%d", i),
				Type: ft.Out(i),
			}
		}
		if retObj.TupleSize() != retNum-1 {
			results[retNum-1] = reflect.ValueOf(fmt.Errorf("param error"))
			return results
		}
		retStruct, err := Python2GoInterface(retObj, reflect.StructOf(fieldType))
		if err != nil {
			results[retNum-1] = reflect.ValueOf(err)
			return results
		}
		retStructObj := reflect.ValueOf(retStruct)
		for i := 0; i < retNum-1; i++ {
			results[i] = retStructObj.Field(i)
		}
		return results
	})
	return f.Interface().(T), nil
}
