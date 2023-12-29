package random

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewRandomString(t *testing.T) {
	tests := map[string]struct {
		size int
	}{
		"size=1":   {size: 1},
		"size=5":   {size: 5},
		"size=10":  {size: 10},
		"size=100": {size: 100},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			rndStr1 := NewRandomString(testCase.size)
			rndStr2 := NewRandomString(testCase.size)

			require.Len(t, rndStr1, testCase.size)
			require.Len(t, rndStr2, testCase.size)
			require.NotEqual(t, rndStr1, rndStr2)
		})
	}
}
