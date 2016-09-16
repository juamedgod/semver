package semver

import (
	"fmt"
	"regexp"
	"strings"
)

var expressionRe = regexp.MustCompile(
	`(` + rangeExpr.String() + `)\s*` +
		`(?P<union>(\|\||\s*))\s*(?P<rest>.*)`)

// Expression defines a semver expression
type Expression interface {
	// Evaluate checks if the provided version v is matches the expression
	Evaluate(v *Version) bool
}
type exprCondition struct {
	Op        string
	Operator1 Expression
	Operator2 Expression
}

type trueCondition struct {
}

func (c *trueCondition) Evaluate(v *Version) bool {
	return true
}
func (c *exprCondition) Evaluate(v *Version) bool {
	if c.Op == "AND" {
		return c.Operator1.Evaluate(v) && c.Operator2.Evaluate(v)
	}
	return c.Operator1.Evaluate(v) || c.Operator2.Evaluate(v)
}

// MustParseExpr parses a semver string
// It returns the expression if str is well formed and panics otherwise
func MustParseExpr(str string) (expr Expression) {
	expr, err := ParseExpr(str)
	if err != nil {
		panic(err)
	}
	return expr
}

// ParseExpr parses a semver string
// It returns the expression if str is well formed and a non-nil error otherwise
func ParseExpr(str string) (Expression, error) {
	text := str
	var condition Expression
	condition = &trueCondition{}
	op := "AND"
	for {
		matched, mapping := namedReEvaluate(expressionRe, text)
		if !matched {
			break
		}
		text = mapping["rest"]
		var ev Expression
		switch {
		case mapping["range"] != "":
			ev = MustParseRange(mapping["range"])
		}

		condition = &exprCondition{Op: op, Operator1: condition, Operator2: ev}
		if mapping["union"] == `||` {
			op = "OR"
		} else {
			op = "AND"
		}
	}
	if condition == nil {
		return condition, fmt.Errorf(`Cannot parse expression %q`, str)
	} else if strings.TrimSpace(text) != "" {
		return condition, fmt.Errorf(`Extra characters found in expression: %q`, text)
	}
	return condition, nil
}
