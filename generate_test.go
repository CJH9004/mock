package mock

import (
	"fmt"
	"testing"
)

func TestGenInt63n(t *testing.T) {
	results := [10]int64{}
	for i := range results {
		results[i] = genInt63n(10)
	}
	fmt.Println(results)
}
