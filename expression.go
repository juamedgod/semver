package semver

import (
	"fmt"
	"regexp"
	"strings"
)

var expressionRe = regexp.MustCompile(
	`(` + rangeExpr.String() + `)\s*` +
		`(?P<union>(\|\||\s*))\s*(?P<rest>.*)`)

type Expression interface {
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

func MustParseExpr(str string) (expr Expression) {
	expr, err := ParseExpr(str)
	if err != nil {
		panic(err)
	}
	return expr
}
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
