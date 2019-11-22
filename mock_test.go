package mock

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMockInt(t *testing.T) {
	m := New(time.Now().UnixNano(), nil)
	var err error
	count := 10

	// default [1, 10)
	var n int64
	for i := 0; i < count; i++ {
		n = 0
		err = m.Mock("", &n)
		assert.Nil(t, err)
		assert.True(t, n >= 1)
		assert.True(t, n < 10)
	}

	// [-10, 1)
	for i := 0; i < count; i++ {
		n = 0
		err = m.Mock("range(-10)", &n)
		assert.Nil(t, err)
		assert.True(t, n >= -10)
		assert.True(t, n < 1)
	}

	// [-1, -1)
	for i := 0; i < count; i++ {
		n = 0
		err = m.Mock("range(-1, -1)", &n)
		assert.Nil(t, err)
		assert.Equal(t, int64(-1), n)
	}

	// date
	n = 0
	err = m.Mock("type(date) format(ms)", &n)
	assert.Nil(t, err)
	assert.True(t, n > 1000000000000)
	assert.True(t, n < 10000000000000)

	// int8
	var n8 int8
	for i := 0; i < count; i++ {
		n8 = 0
		err = m.Mock("range(-1, -1)", &n8)
		assert.Nil(t, err)
		assert.Equal(t, int8(-1), n8)
	}
}

func TestMockUint(t *testing.T) {
	m := New(time.Now().UnixNano(), nil)
	var err error
	count := 10

	var n uint
	for i := 0; i < count; i++ {
		n = 0
		err = m.Mock("", &n)
		assert.Nil(t, err)
		assert.True(t, n >= 1)
		assert.True(t, n < 10)
	}
}

func TestMockFloat(t *testing.T) {
	m := New(time.Now().UnixNano(), nil)
	var err error
	count := 10

	var n float64
	for i := 0; i < count; i++ {
		n = 0
		err = m.Mock("", &n)
		assert.Nil(t, err)
		assert.True(t, n >= 1)
		assert.True(t, n < 10)
	}
}

func TestMockString(t *testing.T) {
	m := New(time.Now().UnixNano(), nil)
	var err error
	count := 10

	var n string
	for i := 0; i < count; i++ {
		n = ""
		err = m.Mock("", &n)
		assert.Nil(t, err)
		assert.True(t, len(n) >= 1)
		assert.True(t, len(n) < 10)
	}

	for _, typ := range TypeList {
		n = ""
		err = m.Mock(fmt.Sprintf("type(%s)", typ), &n)
		assert.Nil(t, err)
		assert.True(t, len(n) > 0)
	}

	n = ""
	err = m.Mock("type(word) range(-1)", &n)
	assert.NotNil(t, err)
}

func TestMockSlice(t *testing.T) {
	m := New(time.Now().UnixNano(), nil)
	var err error
	count := 10

	var n []string
	for i := 0; i < count; i++ {
		n = []string{}
		err = m.Mock("", &n)
		assert.Nil(t, err)
		assert.True(t, len(n) >= 1)
		assert.True(t, len(n) < 10)
		for _, v := range n {
			assert.True(t, len(v) >= 1)
			assert.True(t, len(v) < 10)
		}
	}
}

func TestMockMap(t *testing.T) {
	m := New(time.Now().UnixNano(), nil)
	var err error
	count := 10

	var n map[string]float64
	for i := 0; i < count; i++ {
		n = map[string]float64{}
		err = m.Mock("", &n)
		assert.Nil(t, err)
		assert.True(t, len(n) >= 1)
		assert.True(t, len(n) < 10)
		for k, v := range n {
			assert.True(t, len(k) >= 1)
			assert.True(t, len(k) < 10)
			assert.True(t, v >= 1)
			assert.True(t, v < 10)
		}
	}
}

func TestMockArray(t *testing.T) {
	m := New(time.Now().UnixNano(), nil)
	var err error

	var n [10]int
	err = m.Mock("", &n)
	assert.Nil(t, err)
	for _, v := range n {
		assert.True(t, v >= 1)
		assert.True(t, v < 10)
	}
}

func TestMockStruct(t *testing.T) {
	m := New(time.Now().UnixNano(), nil)
	var err error

	type N struct {
		Set     string
		Default string
		Slice   []string
		Embed   struct {
			A int
		}
	}
	n := N{Set: "asdf"}
	err = m.Mock("", &n)
	assert.Nil(t, err)
	assert.Equal(t, "asdf", n.Set)
	assert.True(t, len(n.Default) >= 1)
	assert.True(t, len(n.Default) < 10)
	assert.True(t, len(n.Slice) >= 1)
	assert.True(t, len(n.Slice) < 10)
	assert.True(t, n.Embed.A >= 1)
	assert.True(t, n.Embed.A < 10)
}

func TestCustomizedGenFuncs(t *testing.T) {
	m := New(time.Now().UnixNano(), nil)
	m.SetGenFuncs(GenFuncs{
		"genInt": func() interface{} { return 101010 },
	})
	var err error

	var n int
	err = m.Mock("mock(genInt)", &n)
	assert.Nil(t, err)
	assert.Equal(t, 101010, n)
}

func TestEmbedTags(t *testing.T) {
	m := New(time.Now().UnixNano(), nil)
	m.SetTags(map[string]string{
		"bigInt": "range(123456, 123457)",
	})
	var err error

	var n int
	err = m.Mock("tag(bigInt)", &n)
	assert.Nil(t, err)
	assert.Equal(t, 123456, n)
}

func TestKeyAndElem(t *testing.T) {
	m := New(time.Now().UnixNano(), nil)
	m.SetTags(map[string]string{
		"word12":     "type(word) range(12, 12)",
		"sentence20": "type(sentence) range(20, 20)",
	})
	var err error

	var n map[string]string
	err = m.Mock("range(10,10) key(word12) elem(sentence20)", &n)
	assert.Nil(t, err)
	assert.Equal(t, 10, len(n))
	for k, v := range n {
		assert.Equal(t, 12, len(k))
		assert.Equal(t, 20, len(strings.Fields(v)))
	}
}

func TestFormats(t *testing.T) {
	m := New(time.Now().UnixNano(), nil)
	m.SetFormats(map[string]string{
		"dash": "2006-01-02",
	})
	var err error

	var n string
	err = m.Mock("type(date) format(dash)", &n)
	assert.Nil(t, err)
	assert.Equal(t, time.Now().Format("2006-01-02"), n)
}
