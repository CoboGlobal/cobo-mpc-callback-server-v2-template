package utils

// IsEqualSet checks if two input lists have the same set of items, and
// the same counts of duplicated items (if any).
//
//	IsEqualSet([1,2,3,3], [3,3,2,1]) == true
//	IsEqualSet([1,2,3,3], [1,2,2,3]) == false
func IsEqualSet[T comparable](s1, s2 []T) bool {
	return len(ElementDiff(s1, s2)) == 0
}

// ElementDiff counts the differences of occurrence of elements. Example:
//
//	ElementDiff([a,a,b], [a,b,c,c]) == {a:1, c:-2}
//	ElementDiff([a,b,c], [a,b,c]) == {empty map}
func ElementDiff[T comparable](listA, listB []T) map[T]int {
	m := make(map[T]int)
	for _, item := range listA {
		m[item] += 1
	}
	for _, item := range listB {
		m[item] -= 1
		if m[item] == 0 {
			delete(m, item)
		}
	}
	return m
}

// IsSubset checks if listA is a subset of listB.
func IsSubset[T comparable](listA, listB []T) bool {
	diff := ElementDiff(listA, listB)
	for _, count := range diff {
		if count > 0 {
			return false
		}
	}
	return true
}
