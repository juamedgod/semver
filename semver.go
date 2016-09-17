package semver

import "fmt"

// Valid checks if a provided string is a valid version
func Valid(str string) bool {
	if _, err := ParseVersion(str); err == nil {
		return true
	}
	return false
}

func toVersion(e interface{}) (*Version, error) {
	switch v := e.(type) {
	case *Version:
		return v, nil
	case string:
		return ParseVersion(v)
	default:
		return nil, fmt.Errorf("unknown element type")
	}
}

func toExpression(e interface{}) (Expression, error) {
	switch v := e.(type) {
	case Expression:
		return v, nil
	case string:
		return ParseExpr(v)
	default:
		return nil, fmt.Errorf("unknown element type: %T", v)
	}
}

// Satisfies receives a version (*Version or its string representation) and an expression (Expression or its string representation)
// and returns whether the version satisfies the expression or not
func Satisfies(version interface{}, expr interface{}) (bool, error) {
	var err error
	var v *Version
	var e Expression
	if v, err = toVersion(version); err != nil {
		return false, fmt.Errorf("Cannot parse version: %v", err)
	}
	if e, err = toExpression(expr); err != nil {
		return false, fmt.Errorf("Cannot parse expression: %v", err)
	}
	return e.Matches(v), nil
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
