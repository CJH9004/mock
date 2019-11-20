package mock

import (
	"errors"
	"reflect"
)

type MockFunc func() interface{}

type ValidFunc func(interface{}) bool

type MockFuncs map[string]MockFunc

type ValidFuncs map[string]ValidFunc

type Mocker interface {
	Mock(data interface{}) error
	Valid(data interface{}) (bool, error)
	SetMockFuncs(fns map[string]func() interface{})
	SetValidFuncs(fns map[string]func(interface{}) bool)
}

type mocker struct {
	mockFuncs  MockFuncs
	validFuncs ValidFuncs
}

func (m *mocker) SetMockFuncs(fns MockFuncs) {
	m.mockFuncs = fns
}

func (m *mocker) SetValidFuncs(fns ValidFuncs) {
	m.validFuncs = fns
}

func (m *mocker) Mock(data interface{}) error {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Ptr {
		return errors.New("not a pointer")
	}
	if v.Elem().Kind() != reflect.Struct {
		return errors.New("not a pointer of a struct")
	}
	m.mock("", v.Elem())
	return nil
}

func (m *mocker) Valid(data interface{}) (bool, error) {
	return true, nil
}

func (m *mocker) mock(tags string, v reflect.Value) {
	
	switch v.Type().Kind() {
	case reflect.Ptr:
		m.mock("", v.Elem())
	case reflect.Struct:
		m.mockStruct(v)
	case reflect.Slice:
		m.mockSlice(tags, v)
	case reflect.Array:
		m.mockArray(tags, v)
	case reflect.Map:
		m.mockMap(tags, v)
	default:
		m.mockField(tags, v)
	}
}

func (m *mocker) mockStruct(v reflect.Value) {
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		vf := v.Field(i)
		tf := t.Field(i)
		tags := tf.Tag.Get("mock")
		if !v.Field(i).CanSet() || tags == "-" || !vf.IsZero() {
			continue
		}
		m.mock(tags, vf)
	}
}

func (m *mocker) mockField(tags string, v reflect.Value) {
	switch v.Type().Kind() {
	case reflect.String:
		v.SetString(m.mockString(tags))
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		v.SetInt(m.mockInt(tags))
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		v.SetUint(m.mockUint(tags))
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		v.SetFloat(m.mockFloat(tags))
		// default:
		// 	log.Println("Unsupported type:", v.Type().Kind())
	}
}

func (m *mocker) mockString(tags string) string {
	return "asdf"
}

func (m *mocker) mockInt(tags string) int64 {
	return -10
}

func (m *mocker) mockUint(tags string) uint64 {
	return 10
}

func (m *mocker) mockFloat(tags string) float64 {
	return 10.101
}

func (m *mocker) mockSlice(tags string, v reflect.Value) {
	v.Set(reflect.MakeSlice(v.Type(), 3, 3))
	for i := 0; i < v.Len(); i++ {
		m.mock(tags, v.Index(i))
	}
}

func (m *mocker) mockArray(tags string, v reflect.Value) {
	for i := 0; i < v.Len(); i++ {
		m.mock(tags, v.Index(i))
	}
}

func (m *mocker) mockMap(tags string, v reflect.Value) {
	if v.Type().Key().Kind() != reflect.String {
		// log.Fatal("Unsupported map key type:", v.Type().Key().Kind())
		return
	}

	v.Set(reflect.MakeMapWithSize(v.Type(), 3))
	for i := 0; i < 3; i++ {
		key := reflect.ValueOf(m.mockString(tags))
		value := reflect.New(v.Type().Elem())
		m.mock(tags, value)
		v.SetMapIndex(key, value.Elem())
	}
}
