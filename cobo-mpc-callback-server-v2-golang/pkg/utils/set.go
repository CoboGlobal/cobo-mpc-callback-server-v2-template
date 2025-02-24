package utils

import mapset "github.com/deckarep/golang-set/v2"

func IsEqualSet[T comparable](s1, s2 []T) bool {
	if len(s1) == 0 || len(s2) == 0 {
		return false
	}

	set1 := mapset.NewSet[T]()
	set2 := mapset.NewSet[T]()

	for _, v := range s1 {
		set1.Add(v)
	}
	for _, v := range s2 {
		set2.Add(v)
	}

	return set1.Equal(set2)
}

func IsSubset[T comparable](s1, s2 []T) bool {
	if len(s1) == 0 || len(s2) == 0 {
		return false
	}

	set1 := mapset.NewSet[T]()
	set2 := mapset.NewSet[T]()

	for _, v := range s1 {
		set1.Add(v)
	}
	for _, v := range s2 {
		set2.Add(v)
	}

	return set1.IsSubset(set2)
}
