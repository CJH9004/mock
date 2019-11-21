package mock

import (
	"fmt"
	"testing"
)

func TestMock(t *testing.T) {
	type SA struct {
		A string `mock:"min(10) max(10)"`
	}
	type SB struct {
		B string `mock:"email"`
	}
	type SC struct {
		a string // should ignore
		SA
		B   SB
		C   string `mock:"type(date)"`
		Set string `mock:"type(email)"`
		D   []struct {
			LA string `mock:"type(phone)"`
		} `mock:"min(3) max(10)"`
		E [3]struct {
			LA string `mock:"type(url)"`
		}
		Map map[string]struct {
			LA string `mock:"type(ipv4)"`
		} `mock:"min(0) max(4)"`
	}
	sc := SC{Set: "set"}
	mocker := New(nil, nil)
	err := mocker.Mock("", &sc)
	fmt.Println(sc)
	fmt.Println(err)
}
