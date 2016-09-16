package semver

import "testing"

var exprTestBattery = map[string]map[string]bool{
	"1.x || >=2.5.0 || 5.0.0 - 7.2.3": {
		"1.2.3":   true,
		"3.0":     true,
		"5.45.23": true,
		"2.5.0":   true,
		"7.2.3":   true,
		"5.0.0":   true,
		"1.0.0":   true,
		"0.9.3":   false,
		"2.4.99":  false,
	},
	">=1.2.7 <1.3.0": {
		"1.2.7":  true,
		"1.2.8":  true,
		"1.2.99": true,
		"1.1.0":  false,
		"1.3.0":  false,
		"1.2.6":  false,
	},
	"1.2.7 || >=1.2.9 <2.0.0": {
		"1.2.7": true,
		"1.2.9": true,
		"1.4.6": true,
		"1.2.8": false,
		"2.0.0": false,
	},
}

func TestPerseExpr(t *testing.T) {
	for _, battery := range []map[string]map[string]bool{
		// Ranges are still expressions
		rangeTestBattery,
		exprTestBattery,
	} {
		for exprStr, data := range battery {
			e := MustParseExpr(exprStr)
			for vStr, result := range data {
				v := MustParseVersion(vStr)
				if e.Evaluate(v) != result {
					t.Errorf("Expected %q of %v to evaluate to %v", exprStr, v, result)
				}
			}
		}
	}
}
