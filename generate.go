package mock

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// todo: add interface

type generator struct {
	rand *rand.Rand
}

// Chars contains gen default string avaliable chars
const Chars = `0123456789abcdefghijklmnopqrstuvwxyz!@#$%^&*(){}[]<>?,./\|:";'~`

func (g generator) fromValues(vals []interface{}) interface{} {
	return vals[g.int63n(int64(len(vals)))]
}

func (g generator) int63n(n int64) int64 {
	if n == 0 {
		return 0
	}
	return g.rand.Int63n(n)
}

func (g generator) int(tag Tag) int64 {
	return g.int63n(tag.Max-tag.Min) + tag.Min
}

func (g generator) uint(tag Tag) uint64 {
	return g.rand.Uint64()%uint64(tag.Max-tag.Min) + uint64(tag.Min)
}

func (g generator) float(tag Tag) float64 {
	return g.rand.Float64()*float64(tag.Max-tag.Min) + float64(tag.Min)
}

func (g generator) string(tag Tag) string {
	if len(tag.Values) > 0 {
		return g.fromValues(tag.Values).(string)
	}

	if isInTypeList(tag.Type) {
		switch tag.Type {
		case "date":
			return g.dateString()
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

func (g generator) dateString() string {
	return time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006")
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
