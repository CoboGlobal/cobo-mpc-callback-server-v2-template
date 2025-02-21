package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEqualSet(t *testing.T) {
	type testcase struct {
		name   string
		s1     []int
		s2     []int
		expect bool
	}

	cases := []testcase{
		{
			name:   "non-dup equal",
			s1:     []int{1, 2, 3, 4},
			s2:     []int{4, 3, 2, 1},
			expect: true,
		},
		{
			name:   "dup equal",
			s1:     []int{1, 2, 2, 3, 3, 3, 4},
			s2:     []int{1, 2, 3, 4, 2, 3, 3},
			expect: true,
		},
		{
			name:   "non-dup not equal A",
			s1:     []int{1, 2, 3, 4},
			s2:     []int{1, 2, 3, 5},
			expect: false,
		},
		{
			name:   "non-dup not equal B",
			s1:     []int{1, 2, 3, 4},
			s2:     []int{1, 2, 3},
			expect: false,
		},
		{
			name:   "non-dup not equal C",
			s1:     []int{1, 2, 3},
			s2:     []int{1, 2, 3, 4},
			expect: false,
		},
		{
			name:   "dup not equal A",
			s1:     []int{1, 2, 3},
			s2:     []int{1, 2, 3, 3},
			expect: false,
		},
		{
			name:   "dup not equal B",
			s1:     []int{1, 2, 3, 2},
			s2:     []int{1, 2, 3, 3},
			expect: false,
		},
		{
			name:   "dup not equal C",
			s1:     []int{1, 2, 3, 3},
			s2:     []int{1, 2, 3},
			expect: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := IsEqualSet(c.s1, c.s2)
			assert.Equal(t, c.expect, got)
		})
	}
}

func TestIsSubset(t *testing.T) {
	type testcase struct {
		name   string
		s1     []int
		s2     []int
		expect bool
	}

	cases := []testcase{
		{
			name:   "empty A",
			s1:     []int{},
			s2:     []int{1, 2, 3},
			expect: true,
		},
		{
			name:   "empty B",
			s1:     []int{1, 2, 3},
			s2:     []int{},
			expect: false,
		},
		{
			name:   "equal A",
			s1:     []int{1, 2, 3},
			s2:     []int{1, 2, 3},
			expect: true,
		},
		{
			name:   "non-dup subset A",
			s1:     []int{1, 2},
			s2:     []int{1, 2, 3},
			expect: true,
		},
		{
			name:   "non-dup subset B",
			s1:     []int{1, 2, 3},
			s2:     []int{1, 2},
			expect: false,
		},
		{
			name:   "non-dup subset C",
			s1:     []int{1, 2, 3},
			s2:     []int{1, 2, 4},
			expect: false,
		},
		{
			name:   "dup subset A",
			s1:     []int{1, 2, 2},
			s2:     []int{1, 2, 3},
			expect: false,
		},
		{
			name:   "dup subset B",
			s1:     []int{1, 2, 3},
			s2:     []int{1, 2, 2},
			expect: false,
		},
		{
			name:   "dup subset C",
			s1:     []int{1, 2, 3, 3},
			s2:     []int{1, 2, 2, 4},
			expect: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := IsSubset(c.s1, c.s2)
			assert.Equal(t, c.expect, got)
		})
	}
}
