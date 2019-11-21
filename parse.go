package mock

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// todo: error handle; add tag funcs: key(tag) elem(tag) range(1,2) format(Mon Jan 2 15:04:05 -0700 MST 2006) ; 

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
	MockFunc  string
	ValidFunc string
}

func parseTag(typ, tags string) (t Tag, err error) {
	re := regexp.MustCompile(`(min|max|type|value|mock|valid)\((.+?)\)`)
	fields := re.FindAllStringSubmatch(tags, -1)
	t.Min = 1
	t.Max = 10
	for _, f := range fields {
		switch f[1] {
		case "min":
			if t.Min, err = strconv.ParseInt(f[2], 10, 64); err != nil {
				return t, fmt.Errorf("the param %s in min func is not a number", f[2])
			}
			if t.Min < 0 && (typ == "string" || strings.HasPrefix(typ, "uint")) {
				return t, fmt.Errorf("the fileld type %s need a positive number but the param in min func got %s", typ, f[2])
			}
		case "max":
			if t.Max, err = strconv.ParseInt(f[2], 10, 64); err != nil {
				return t, fmt.Errorf("the param %s in max func is not a number", f[2])
			}
		case "type":
			if !isInTypeList(f[2]) {
				return t, fmt.Errorf("the type of %s is invalid", f[2])
			}
			if f[2] == "date" && typ != "int64" && typ != "string" {
				return t, fmt.Errorf("the type of date need field type of int64 or string")
			}
			if f[2] != "date" && typ != "string" {
				return t, fmt.Errorf("the type of %s need field type string", f[2])
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
						return t, fmt.Errorf("the param %s in values func is not a int but the filed type is %s", v, typ)
					}
					t.Values = append(t.Values, n)
				}
			case strings.HasPrefix(typ, "uint"):
				for _, v := range vals {
					var n uint64
					if n, err = strconv.ParseUint(v, 10, 64); err != nil {
						return t, fmt.Errorf("the param %s in values func is not a uint but the filed type is %s", v, typ)
					}
					t.Values = append(t.Values, n)
				}
			case strings.HasPrefix(typ, "float"):
				for _, v := range vals {
					var n float64
					if n, err = strconv.ParseFloat(v, 64); err != nil {
						return t, fmt.Errorf("the param %s in values func is not a float but the filed type is %s", v, typ)
					}
					t.Values = append(t.Values, n)
				}
			}
		case "mock":
			t.MockFunc = f[2]
		case "valid":
			t.ValidFunc = f[2]
		}
	}
	if t.Max < t.Min {
		return t, fmt.Errorf("the param %d in max func is less than the param %d in min func", t.Max, t.Min)
	}
	if t.Min < 0 && (t.Type == "word" || t.Type == "sentence") {
		return t, fmt.Errorf("the type %s need set min more than 1 but got %d", t.Type, t.Min)
	}
	return t, nil
}
