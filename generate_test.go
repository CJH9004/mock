package mock

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestFloat(t *testing.T) {
	gen := generator{rand: rand.New(rand.NewSource(time.Now().UnixNano()))}
	count := 1000000
	hit := 0
	for i := 0; i < count; i++ {
		x := gen.float(Tag{Min: 0, Max: 1})
		y := gen.float(Tag{Min: 0, Max: 1})
		if x*x+y*y < 1 {
			hit++
		}
	}
	res := float64(hit) / float64(count) * 4
	if math.Abs(res-math.Pi) > 0.01 {
		t.Fatalf("excepted %f, got %f", math.Pi, res)
	}
}

func TestInt63n(t *testing.T) {
	gen := generator{rand: rand.New(rand.NewSource(time.Now().UnixNano()))}
	count := int64(1000000)
	hit := 0
	for i := int64(0); i < count; i++ {
		n := gen.int63n(count)
		if n < count/2 {
			hit++
		}
	}
	res := float64(hit) / float64(count)
	if math.Abs(res-0.5) > 0.01 {
		t.Fatalf("excepted %f, got %f", 0.5, res)
	}
}
