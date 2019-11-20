package mock

import (
	"fmt"
	"testing"
)

func TestMock(t *testing.T) {
	type SA struct {
		A string `mock:"A"`
	}
	type SB struct {
		B string `mock:"B"`
	}
	type SC struct {
		a string // should ignore
		SA
		B   SB     `mock:"SB"`
		C   string `mock:"C"`
		Set string `mock:"C"`
		D   []struct {
			LA string `mock:"LA"`
		} `mock:"D"`
		E [3]struct {
			LA string `mock:"LA"`
		} `mock:"E"`
		Map map[string]struct {
			LA string `mock:"LA"`
		} `mock:"Map"`
	}
	sc := SC{Set: "set"}
	Mock(&sc)
	fmt.Println(sc)
}
