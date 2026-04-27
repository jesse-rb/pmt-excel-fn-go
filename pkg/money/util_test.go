package money

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToCents(t *testing.T) {
	cases := []float64{
		0,
		0.01,
		0.11,
		1.23,
		23.42,
	}

	expected := []int64{
		0,
		1,
		11,
		123,
		2342,
	}

	for i := range cases {
		actual := ToCents(cases[i])
		assert.Equal(t, expected[i], actual)
	}
}
