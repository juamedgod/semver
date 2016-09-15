package semver

func Valid(str string) bool {
	if _, err := ParseVersion(str); err == nil {
		return true
	}
	return false
}

func Greater(e1 Comparable, e2 Comparable) bool {
	return e1.UpperLimit().Greater(e2.UpperLimit())
}

func Less(e1 Comparable, e2 Comparable) bool {
	return e1.LowerLimit().Less(e2.LowerLimit())
}

func GreaterOrEqual(e1 Comparable, e2 Comparable) bool {
	return e1.UpperLimit().GreaterOrEqual(e2.UpperLimit())
}

func LessOrEqual(e1 Comparable, e2 Comparable) bool {
	return e1.LowerLimit().LessOrEqual(e2.LowerLimit())
}

func Equal(e1 Comparable, e2 Comparable) bool {
	return e1.LowerLimit().Equal(e2.LowerLimit()) && e1.UpperLimit().Equal(e2.UpperLimit())
}
