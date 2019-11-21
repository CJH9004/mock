package mock

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
)

// todo: 嵌套tag，将对tag的处理转移到generator

// MockFunc is costomized mock func
type MockFunc func() interface{}

// ValidFunc is costomized valid func
type ValidFunc func(interface{}) bool

// MockFuncs is costomized mock funcs map
type MockFuncs map[string]MockFunc

// ValidFuncs is costomized valid funcs map
type ValidFuncs map[string]ValidFunc

// Mocker mock the data
type Mocker interface {
	Mock(tags string, data interface{}) error
	Valid(tags string, data interface{}) (bool, error)
	SetMockFuncs(fns MockFuncs)
	SetValidFuncs(fns ValidFuncs)
	// Errors() []error
}

type mocker struct {
	mockFuncs  MockFuncs
	validFuncs ValidFuncs
	gen        generator
	err        error
}

// New return a Mocker
func New(mockFuncs MockFuncs, validFuncs ValidFuncs, seed int64) Mocker {
	return &mocker{
		mockFuncs:  mockFuncs,
		validFuncs: validFuncs,
		gen:        generator{rand: rand.New(rand.NewSource(seed))},
	}
}

func (m *mocker) SetMockFuncs(fns MockFuncs) {
	m.mockFuncs = fns
}

func (m *mocker) SetValidFuncs(fns ValidFuncs) {
	m.validFuncs = fns
}

func (m *mocker) Mock(tags string, data interface{}) error {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Ptr {
		return errors.New("not a pointer")
	}
	m.mock(tags, v.Elem())
	return m.err
}

func (m *mocker) Valid(tags string, data interface{}) (bool, error) {
	return true, nil
}

func (m *mocker) mock(tags string, v reflect.Value) {
	t := m.parseTag(v.Type().Name(), tags)
	if fn, ok := m.mockFuncs[t.MockFunc]; ok {
		v.Set(reflect.ValueOf(fn()))
	}
	switch v.Type().Kind() {
	case reflect.Ptr:
		m.mock(tags, v.Elem())
	case reflect.Struct:
		m.mockStruct(t, v)
	case reflect.Slice:
		m.mockSlice(t, v)
	case reflect.Array:
		m.mockArray(t, v)
	case reflect.Map:
		m.mockMap(t, v)
	default:
		m.mockField(t, v)
	}
}

func (m *mocker) parseTag(typ, tags string) Tag {
	var t Tag
	var err error
	if t, err = parseTag(typ, tags); err != nil {
		m.err = err
	}
	return t
}

func (m *mocker) mockStruct(tag Tag, v reflect.Value) {
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

func (m *mocker) mockField(t Tag, v reflect.Value) {
	switch v.Type().Kind() {
	case reflect.String:
		m.mockString(t, v)
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		m.mockInt(t, v)
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		m.mockUint(t, v)
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		m.mockFloat(t, v)
		// default:
		// 	log.Println("Unsupported type:", v.Type().Kind())
	}
}

func (m *mocker) mockString(t Tag, v reflect.Value) {
	v.SetString(m.gen.string(t))
}

func (m *mocker) mockInt(t Tag, v reflect.Value) {
	v.SetInt(m.gen.int(t))
}

func (m *mocker) mockUint(t Tag, v reflect.Value) {
	v.SetUint(m.gen.uint(t))
}

func (m *mocker) mockFloat(t Tag, v reflect.Value) {
	v.SetFloat(m.gen.float(t))
}

func (m *mocker) mockSlice(t Tag, v reflect.Value) {
	length := m.gen.int(t)
	v.Set(reflect.MakeSlice(v.Type(), int(length), int(length)))
	for i := 0; i < v.Len(); i++ {
		m.mock("", v.Index(i))
	}
}

func (m *mocker) mockArray(t Tag, v reflect.Value) {
	for i := 0; i < v.Len(); i++ {
		m.mock("", v.Index(i))
	}
}

func (m *mocker) mockMap(t Tag, v reflect.Value) {
	if v.Type().Key().Kind() != reflect.String {
		m.err = fmt.Errorf("Unsupported map key type: %s", v.Type().Key().Kind())
		return
	}

	length := m.gen.int(t)
	v.Set(reflect.MakeMapWithSize(v.Type(), int(length)))
	for i := 0; i < 3; i++ {
		key := reflect.ValueOf(m.gen.string(m.parseTag("string", "type(word)")))
		value := reflect.New(v.Type().Elem())
		m.mock("", value)
		v.SetMapIndex(key, value.Elem())
	}
}
