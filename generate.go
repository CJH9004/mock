package mock

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Chars contains gen default string avaliable chars
const Chars = `0123456789abcdefghijklmnopqrstuvwxyz!@#$%^&*(){}[]<>?,./\|:";'~`

func getFromValues(vals []interface{}) interface{} {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return vals[r.Intn(len(vals))]
}

func genInt63n(n int64) int64 {
	if n == 0 {
		return 0
	}
	time.Sleep(1 * time.Nanosecond)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int63n(n)
}

func genInt(tag Tag) int64 {
	return genInt63n(tag.Max-tag.Min) + tag.Min
}

func genUint(tag Tag) uint64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Uint64()%uint64(tag.Max-tag.Min) + uint64(tag.Min)
}

func genFloat(tag Tag) float64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Float64()*float64(tag.Max-tag.Min) + float64(tag.Min)
}

func genString(tag Tag) string {
	if len(tag.Values) > 0 {
		return getFromValues(tag.Values).(string)
	}

	if isInTypeList(tag.Type) {
		switch tag.Type {
		case "date":
			return genDateString()
		case "email":
			return genEamil()
		case "phone":
			return genPhone()
		case "url":
			return genURL()
		case "ipv4":
			return genIPv4()
		case "domain":
			return genDomain()
		case "word":
			return genWord(tag.Min, tag.Max)
		case "sentence":
			return genSentence(tag.Min, tag.Max)
		}
	}

	b := make([]byte, genInt63n(tag.Max-tag.Min)+tag.Min)
	for i := range b {
		b[i] = Chars[genInt63n(int64(len(Chars)))]
	}
	return string(b)
}

func genEamil() string {
	return genWord(0, 10) + "@" + genWord(0, 10) + "." + genWord(0, 10)
}

func genDateString() string {
	return time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006")
}

func genPhone() string {
	chars := "0123456789"
	b := make([]byte, 10)
	for i := range b {
		b[i] = chars[genInt63n(int64(len(chars)))]
	}
	return "1" + string(b)
}

func genURL() string {
	return "http://" + genDomain() + "/" + genWord(1, 10)
}

func genIPv4() string {
	return fmt.Sprintf("%d.%d.%d.%d", genInt63n(256), genInt63n(256), genInt63n(256), genInt63n(256))
}

func genDomain() string {
	return "www." + genWord(1, 10) + "." + genWord(1, 10)
}

func genWord(min, max int64) string {
	chars := "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, genInt63n(max-min)+min)
	for i := range b {
		b[i] = chars[genInt63n(int64(len(chars)))]
	}
	return string(b)
}

func genSentence(min, max int64) string {
	words := make([]string, genInt63n(max-min)+min)
	for i := range words {
		words[i] = genWord(1, 10)
	}
	words[0] = strings.Title(words[0])
	return strings.Join(words, " ") + "."
}
