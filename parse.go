package mock

import (
	"strings"
	"fmt"
)

// mock:"min(0) max(10) type('string') value(1,3,4) required mock(func) valid(func)" valid:""

var TypeList = []string{"eamil", "date", "phone", "url", "ipv4", "domain"}

type Tags struct {
	Type      string
	Values    []interface{}
	Min       int64 // default 1
	Max       int64 // default 10
	MockFunc  MockFunc
	ValidFunc ValidFunc
}

func parseTags(tags string) (t Tags) {
	for _, v := range strings.Fields(tags) {
		fmt.Println(v)
	}
	return
}
