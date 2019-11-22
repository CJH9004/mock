package mock

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
)

// GenFunc is costomized mock func
type GenFunc func() interface{}

// ValidFunc is costomized valid func
type ValidFunc func(interface{}) bool

// GenFuncs is costomized mock funcs map
type GenFuncs map[string]GenFunc

// ValidFuncs is costomized valid funcs map
type ValidFuncs map[string]ValidFunc

// Mocker mock the data
type Mocker interface {
	Mock(tags string, data interface{}) error
	Valid(tags string, data interface{}) (bool, error)
	SetGenFuncs(fns GenFuncs)
	SetValidFuncs(fns ValidFuncs)
	SetTags(map[string]string)
	SetFormats(map[string]string)
}

type mocker struct {
	genFuncs   GenFuncs
	validFuncs ValidFuncs
	tags       map[string]string
	formats    map[string]string
	gen        generator
	err        error
}

// Options store the ortions of Mocker
type Options struct {
	GenFuncs   GenFuncs
	ValidFuncs ValidFuncs
	Tags       map[string]string
	Formats    map[string]string
}

// New return a Mocker
func New(seed int64, options *Options) Mocker {
	if options == nil {
		options = &Options{}
	}
	return &mocker{
		genFuncs:   options.GenFuncs,
		validFuncs: options.ValidFuncs,
		tags:       options.Tags,
		formats:    options.Formats,
		gen:        generator{rand: rand.New(rand.NewSource(seed))},
	}
}

func (m *mocker) SetGenFuncs(fns GenFuncs) {
	m.genFuncs = fns
}

func (m *mocker) SetValidFuncs(fns ValidFuncs) {
	m.validFuncs = fns
}

func (m *mocker) SetTags(tags map[string]string) {
	m.tags = tags
}

func (m *mocker) SetFormats(formats map[string]string) {
	m.formats = formats
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
	if v.Type().Kind() == reflect.Ptr {
		m.mock(tags, v.Elem())
	}
	t := m.parseTag(v.Type().Name(), tags)
	if fn, ok := m.genFuncs[t.GenFunc]; ok {
		v.Set(reflect.ValueOf(fn()))
		return
	}
	switch v.Type().Kind() {
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
	if t, err = ParseTag(typ, tags); err != nil {
		m.err = err
	}
	if tag, ok := m.tags[t.Tag]; ok {
		t = m.parseTag(typ, tag)
	}
	if v, ok := m.tags[t.Key]; ok {
		t.Key = v
	}
	if v, ok := m.tags[t.Elem]; ok {
		t.Elem = v
	}
	if v, ok := m.formats[t.Format]; ok {
		t.Format = v
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
		m.mock(t.Elem, v.Index(i))
	}
}

func (m *mocker) mockArray(t Tag, v reflect.Value) {
	for i := 0; i < v.Len(); i++ {
		m.mock(t.Elem, v.Index(i))
	}
}

func (m *mocker) mockMap(t Tag, v reflect.Value) {
	if v.Type().Key().Kind() != reflect.String {
		m.err = fmt.Errorf("Unsupported map key type: %s", v.Type().Key().Kind())
		return
	}

	length := m.gen.int(t)
	v.Set(reflect.MakeMapWithSize(v.Type(), int(length)))
	for i := 0; i < int(length); i++ {
		keyTag := "type(word)"
		if t.Key != "" {
			keyTag = t.Key
		}
		key := reflect.ValueOf(m.gen.string(m.parseTag("string", keyTag)))
		value := reflect.New(v.Type().Elem())
		m.mock(t.Elem, value.Elem())
		v.SetMapIndex(key, value.Elem())
	}
}
