package mock

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Generator gen int uint float and string
type Generator interface {
	Int(string) (int64, error)
	Uint(string) (uint64, error)
	Float(string) (float64, error)
	String(string) (string, error)
}

type generator struct {
	rand *rand.Rand
}

// NewGen return a Generator
func NewGen(rand *rand.Rand) Generator {
	return generator{
		rand: rand,
	}
}

func (g generator) Bool(ts string) (ret bool, err error) {
	t, err := ParseTag("bool", ts)
	if err != nil {
		return false, err
	}
	return g.bool(t), nil
}

func (g generator) Int(ts string) (ret int64, err error) {
	t, err := ParseTag("int", ts)
	if err != nil {
		return 0, err
	}
	return g.int(t), nil
}

func (g generator) Uint(ts string) (ret uint64, err error) {
	t, err := ParseTag("uint", ts)
	if err != nil {
		return 0, err
	}
	return g.uint(t), nil
}

func (g generator) Float(ts string) (ret float64, err error) {
	t, err := ParseTag("float", ts)
	if err != nil {
		return 0, err
	}
	return g.float(t), nil
}

func (g generator) String(ts string) (ret string, err error) {
	t, err := ParseTag("string", ts)
	if err != nil {
		return "", err
	}
	return g.string(t), nil
}

// Chars contains gen default string avaliable chars
const Chars = `0123456789abcdefghijklmnopqrstuvwxyz!@#$%^&*(){}[]<>?,./\|:";'~`

// TimeFormat is default time format
const TimeFormat = "Mon Jan 2 15:04:05 -0700 MST 2006"

func (g generator) fromValues(vals []interface{}) interface{} {
	return vals[g.int63n(int64(len(vals)))]
}

func (g generator) int63n(n int64) int64 {
	if n == 0 {
		return 0
	}
	return g.rand.Int63n(n)
}

func (g generator) bool(tag Tag) bool {
	if len(tag.Values) > 0 {
		return g.fromValues(tag.Values).(bool)
	}

	return g.fromValues([]interface{}{true, false}).(bool)
}

func (g generator) int(tag Tag) int64 {
	if len(tag.Values) > 0 {
		return g.fromValues(tag.Values).(int64)
	}

	if tag.Type == "date" {
		return g.dateUnix(tag)
	}
	return g.int63n(tag.Max-tag.Min) + tag.Min
}

func (g generator) uint(tag Tag) uint64 {
	if len(tag.Values) > 0 {
		return g.fromValues(tag.Values).(uint64)
	}

	return g.rand.Uint64()%uint64(tag.Max-tag.Min) + uint64(tag.Min)
}

func (g generator) float(tag Tag) float64 {
	if len(tag.Values) > 0 {
		return g.fromValues(tag.Values).(float64)
	}

	return g.rand.Float64()*float64(tag.Max-tag.Min) + float64(tag.Min)
}

func (g generator) string(tag Tag) string {
	if len(tag.Values) > 0 {
		return g.fromValues(tag.Values).(string)
	}

	if isInTypeList(tag.Type) {
		switch tag.Type {
		case "date":
			return g.dateString(tag)
		case "email":
			return g.eamil()
		case "phone":
			return g.phone()
		case "url":
			return g.url()
		case "ipv4":
			return g.ipv4()
		case "domain":
			return g.domain()
		case "word":
			return g.word(tag.Min, tag.Max)
		case "sentence":
			return g.sentence(tag.Min, tag.Max)
		}
	}

	b := make([]byte, g.int63n(tag.Max-tag.Min)+tag.Min)
	for i := range b {
		b[i] = Chars[g.int63n(int64(len(Chars)))]
	}
	return string(b)
}

func (g generator) eamil() string {
	return g.word(0, 10) + "@" + g.word(0, 10) + "." + g.word(0, 10)
}

func (g generator) dateString(tag Tag) string {
	format := TimeFormat
	if tag.Format != "" {
		format = tag.Format
	}
	return time.Now().Format(format)
}

func (g generator) dateUnix(tag Tag) int64 {
	switch tag.Format {
	case "ns":
		return time.Now().UnixNano()
	case "ms":
		return time.Now().UnixNano() / 1000000
	default:
		return time.Now().Unix()
	}
}

func (g generator) phone() string {
	chars := "0123456789"
	b := make([]byte, 10)
	for i := range b {
		b[i] = chars[g.int63n(int64(len(chars)))]
	}
	return "1" + string(b)
}

func (g generator) url() string {
	return "http://" + g.domain() + "/" + g.word(1, 10)
}

func (g generator) ipv4() string {
	return fmt.Sprintf("%d.%d.%d.%d", g.int63n(256), g.int63n(256), g.int63n(256), g.int63n(256))
}

func (g generator) domain() string {
	return "www." + g.word(1, 10) + "." + g.word(1, 10)
}

func (g generator) word(min, max int64) string {
	chars := "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, g.int63n(max-min)+min)
	for i := range b {
		b[i] = chars[g.int63n(int64(len(chars)))]
	}
	return string(b)
}

func (g generator) sentence(min, max int64) string {
	words := make([]string, g.int63n(max-min)+min)
	for i := range words {
		words[i] = g.word(1, 10)
	}
	words[0] = strings.Title(words[0])
	return strings.Join(words, " ") + "."
}
