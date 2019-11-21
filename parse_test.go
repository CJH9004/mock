package mock

import (
	"fmt"
	"testing"
)

func TestParseTag(t *testing.T) {
	s := "min(0) max(100) mock(mockFunc) valid(validFunc) type(date) value(a, bc)"
	fmt.Println(parseTag("string", s))
}
