package semver

// Valid checks if a provided string is a valid version
func Valid(str string) bool {
	if _, err := ParseVersion(str); err == nil {
		return true
	}
	return false
}

// Greater checks if element e1 is greater than e2
func Greater(e1 Comparable, e2 Comparable) bool {
	return e1.UpperLimit().Greater(e2.UpperLimit())
}

// Less checks if element e1 is less than e2
func Less(e1 Comparable, e2 Comparable) bool {
	return e1.LowerLimit().Less(e2.LowerLimit())
}

// GreaterOrEqual checks if element e1 is greater or equal than e2
func GreaterOrEqual(e1 Comparable, e2 Comparable) bool {
	return e1.UpperLimit().GreaterOrEqual(e2.UpperLimit())
}

// LessOrEqual checks if element e1 is less or equal to e2
func LessOrEqual(e1 Comparable, e2 Comparable) bool {
	return e1.LowerLimit().LessOrEqual(e2.LowerLimit())
}

// Equal checks if element e1 is equal to e2
func Equal(e1 Comparable, e2 Comparable) bool {
	return e1.LowerLimit().Equal(e2.LowerLimit()) && e1.UpperLimit().Equal(e2.UpperLimit())
}
