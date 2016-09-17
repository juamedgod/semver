package semver

// Valid checks if a provided string is a valid version
func Valid(str string) bool {
	if _, err := ParseVersion(str); err == nil {
		return true
	}
	return false
}

// ONLY VERSIONS FOR NOW (to be implemented: ranges and GlobVersion )

// Greater checks if element e1 is greater than e2
func Greater(v1 Comparable, v2 Comparable) bool {
	return v1.(*Version).Greater(v2.(*Version))
}

// Less checks if element e1 is less than e2
func Less(v1 Comparable, v2 Comparable) bool {
	return v1.(*Version).Less(v2.(*Version))
}

// GreaterOrEqual checks if element e1 is greater or equal than e2
func GreaterOrEqual(v1 Comparable, v2 Comparable) bool {
	return v1.(*Version).GreaterOrEqual(v2.(*Version))
}

// LessOrEqual checks if element e1 is less or equal to e2
func LessOrEqual(v1 Comparable, v2 Comparable) bool {
	return v1.(*Version).LessOrEqual(v2.(*Version))
}

// Equal checks if element e1 is equal to e2
func Equal(v1 Comparable, v2 Comparable) bool {
	return v1.(*Version).Equal(v2.(*Version))
}
