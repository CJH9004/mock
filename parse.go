package mock

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// todo: error handle; add tag funcs: key(tag) elem(tag) range(1,2) format(Mon Jan 2 15:04:05 -0700 MST 2006) ;

// ParamError descripe the tag func param error
type ParamError struct {
	Name   string      // tag name
	Expect string      // except type
	Param  interface{} // actual param
}

func (e ParamError) Error() string {
	return fmt.Sprintf("%s except a %s, got %v", e.Name, e.Expect, e.Param)
}

// NewParamError construct a ParamError
func NewParamError(name, expect string, param interface{}) error {
	return ParamError{
		Name:   name,
		Expect: expect,
		Param:  param,
	}
}

// ConflictError descripe the tag func conflict error
type ConflictError struct {
	Name1  string      // tag name1
	Name2  string      // tag name2
	Detail string      // what conflict
	Param1 interface{} // actual param1
	Param2 interface{} // actual param2
}

func (e ConflictError) Error() string {
	return fmt.Sprintf("%s(%v) conflict with %s(%v), %s", e.Name1, e.Param1, e.Name2, e.Param2, e.Detail)
}

// NewConflictError construct a ConflictError
func NewConflictError(name1 string, param1 interface{}, name2 string, param2 interface{}, detail string) error {
	return ConflictError{
		Name1:  name1,
		Name2:  name2,
		Param1: param1,
		Param2: param2,
		Detail: detail,
	}
}

// TypeList is the avalid type
var TypeList = []string{"eamil", "date", "phone", "url", "ipv4", "domain", "word", "sentence"}

func isInTypeList(s string) bool {
	for _, v := range TypeList {
		if v == s {
			return true
		}
	}
	return false
}

// Tag store the fileds
type Tag struct {
	Type      string
	Values    []interface{}
	Min       int64 // default 1
	Max       int64 // default 10
	Key       string
	Elem      string
	Format    string
	Tag       string
	GenFunc   string
	ValidFunc string
}

// DefaultTag return a tag with default value
func DefaultTag() Tag {
	return Tag{
		Min: 1,
		Max: 10,
	}
}

// ParseTag parse string to Tag
func ParseTag(typ, tags string) (t Tag, err error) {
	t = DefaultTag()
	if tags == "" {
		return t, nil
	}

	re := regexp.MustCompile(`(range|type|value|mock|valid|key|elem|format|tag)\((.+?)\)`)
	fields := re.FindAllStringSubmatch(tags, -1)
	for _, f := range fields {
		switch f[1] {
		case "range":
			vals := strings.Split(f[2], ",")
			if len(vals) == 1 {
				if v, err := strconv.ParseInt(strings.TrimSpace(vals[0]), 10, 64); err == nil {
					if v < t.Min {
						t.Min = v
						t.Max = 1
					} else {
						t.Max = v
					}
				} else {
					return DefaultTag(), NewParamError(f[1], "number", vals[0])
				}
			} else if len(vals) == 2 {
				if t.Min, err = strconv.ParseInt(strings.TrimSpace(vals[0]), 10, 64); err != nil {
					return DefaultTag(), NewParamError(f[1], "number", vals[0])
				}
				if t.Max, err = strconv.ParseInt(strings.TrimSpace(vals[1]), 10, 64); err != nil {
					return DefaultTag(), NewParamError(f[1], "number", vals[1])
				}
				if t.Min > t.Max {
					return DefaultTag(), NewParamError(f[1], "min <= max", fmt.Sprintf("min: %d > max: %d", t.Min, t.Max))
				}
			} else {
				return DefaultTag(), NewParamError(f[1], "one or two number", len(vals))
			}
		case "type":
			if !isInTypeList(f[2]) {
				return DefaultTag(), NewParamError(f[1], strings.Join(TypeList, "/"), f[2])
			}
			if f[2] == "date" && typ != "int64" && typ != "string" {
				return DefaultTag(), NewConflictError("fieldType", typ, f[1], f[2], "date need field type int64 or string")
			}
			if f[2] != "date" && typ != "string" {
				return DefaultTag(), NewConflictError("fieldType", typ, f[1], f[2], fmt.Sprintf("%s need field type string", f[2]))
			}
			t.Type = f[2]
		case "value":
			vals := strings.Split(f[2], ",")
			t.Values = make([]interface{}, 0, len(vals))
			for i, v := range vals {
				vals[i] = strings.TrimSpace(v)
			}
			switch {
			case typ == "string":
				for _, v := range vals {
					t.Values = append(t.Values, v)
				}
			case strings.HasPrefix(typ, "int"):
				for _, v := range vals {
					var n int64
					if n, err = strconv.ParseInt(v, 10, 64); err != nil {
						return DefaultTag(), NewConflictError("fieldType", typ, f[1], v, err.Error())
					}
					t.Values = append(t.Values, n)
				}
			case strings.HasPrefix(typ, "uint"):
				for _, v := range vals {
					var n uint64
					if n, err = strconv.ParseUint(v, 10, 64); err != nil {
						return DefaultTag(), NewConflictError("fieldType", typ, f[1], v, err.Error())
					}
					t.Values = append(t.Values, n)
				}
			case strings.HasPrefix(typ, "float"):
				for _, v := range vals {
					var n float64
					if n, err = strconv.ParseFloat(v, 64); err != nil {
						return DefaultTag(), NewConflictError("fieldType", typ, f[1], v, err.Error())
					}
					t.Values = append(t.Values, n)
				}
			case typ == "bool":
				for _, v := range vals {
					if v != "true" && v != "false" {
						return DefaultTag(), NewConflictError("fieldType", typ, f[1], v, err.Error())
					}
					n := false
					if v == "true" {
						n = true
					}
					t.Values = append(t.Values, n)
				}
			}
		case "mock":
			t.GenFunc = f[2]
		case "valid":
			t.ValidFunc = f[2]
		case "key":
			t.Key = f[2]
		case "elem":
			t.Elem = f[2]
		case "format":
			t.Format = f[2]
		case "tag":
			t.Tag = f[2]
		}
	}
	if t.Min < 0 && (t.Type == "word" || t.Type == "sentence") {
		return DefaultTag(), NewConflictError("type", t.Type, "min", t.Min, "word and sentence length need greater than 1")
	}
	if t.Min < 0 && (typ == "slice" || typ == "array" || typ == "map") {
		return DefaultTag(), NewConflictError("fieldType", typ, "min", t.Min, "slice/array/map length need greater than 1")
	}
	return t, nil
}
