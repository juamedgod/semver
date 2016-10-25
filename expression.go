package semver

import (
	"fmt"
	"regexp"
	"strings"
)

var expressionRe = regexp.MustCompile(
	`(` + rangeExpr.String() + `)\s*` +
		`(?P<union>(\|\||\s*))\s*(?P<rest>.*)`)

type evaluable interface {
	// evaluate checks if the provided version v is matches the expression
	evaluate(v *Version) bool
}

// Expression defines a semver expression
type Expression interface {
	Matches(v *Version) bool
	String() string
}

type semverExpression struct {
	str string
	c   evaluable
}

func (e *semverExpression) String() string {
	return e.str
}

// Matches checks if the provided version v is accepted by the expression
func (e *semverExpression) Matches(v *Version) bool {
	return e.c.evaluate(v)
}

type exprCondition struct {
	Op        string
	Operator1 evaluable
	Operator2 evaluable
}

type trueCondition struct {
}

func (c *trueCondition) evaluate(v *Version) bool {
	return true
}

func (c *exprCondition) evaluate(v *Version) bool {
	if c.Op == "AND" {
		return c.Operator1.evaluate(v) && c.Operator2.evaluate(v)
	}
	return c.Operator1.evaluate(v) || c.Operator2.evaluate(v)
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
	var condition evaluable
	condition = &trueCondition{}
	op := "AND"
	for {
		matched, mapping := namedReEvaluate(expressionRe, text)
		if !matched {
			break
		}
		text = mapping["rest"]
		var ev evaluable
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
		return nil, fmt.Errorf(`Cannot parse expression %q`, str)
	} else if strings.TrimSpace(text) != "" {
		return &semverExpression{c: condition, str: str}, fmt.Errorf(`Extra characters found in expression: %q`, text)
	}
	return &semverExpression{c: condition, str: str}, nil
}
