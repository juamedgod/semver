package semver

import (
	"fmt"
	"testing"
)

// The tilde is not actually included here, so we can prefix the alternative forms: ~ or ~>
var tildeTests = map[string]map[string]bool{
	"1.3.4": {
		"0.0.1":    false,
		"1.3.4":    true,
		"1.1.0":    false,
		"1.3.3":    false,
		"1.3.5":    true,
		"1.3.9999": true,
		"1.4.0":    false,
		"1.6.7":    false,
		"1.13.0":   false,
		"2.0.0":    false,
	},
	"1.3": {
		"0.0.1":    false,
		"1.3.0":    true,
		"1.3.4":    true,
		"1.1.0":    false,
		"1.3.3":    true,
		"1.3.5":    true,
		"1.3.9999": true,
		"1.4.0":    false,
		"1.6.7":    false,
		"1.13.0":   false,
		"2.0.0":    false,
	},
	"1": {
		"0.0.1":    false,
		"1.3.0":    true,
		"1.3.4":    true,
		"1.1.0":    true,
		"1.3.3":    true,
		"1.3.5":    true,
		"1.3.9999": true,
		"1.4.0":    true,
		"1.6.7":    true,
		"1.13.0":   true,
		"2.0.0":    false,
	},
}

var rangeTestBattery = map[string]map[string]bool{
	"^1.3.4": {
		"1.3.4":  true,
		"1.1.0":  false,
		"1.3.3":  false,
		"1.3.5":  true,
		"1.6.7":  true,
		"1.13.0": true,
		"2.0.0":  false,
	},
	"^0.3.1": {
		"0.3.0":    false,
		"0.2.1":    false,
		"0.3.1":    true,
		"1.2.3":    false,
		"0.4.0":    false,
		"0.3.9999": true,
	},
	"^0.0.2": {
		"1.0.1":   false,
		"0.1.0":   false,
		"0.0.2":   true,
		"0.0.999": false,
	},
	"=1.3.4": {
		"1.3.4": true,
		"1.3.0": false,
		"1.3":   false,
		"1.3.1": false,
		"2.4.2": false,
	},
	"=1.3": {
		"1.3.0": true,
		"1.3":   true,
		"1.3.1": true,
		"2.3":   false,
	},
	"1.3": {
		"1.3.0": true,
		"1.3":   true,
		"1.3.1": true,
		"2.3":   false,
	},
	"1.3.2 - 1.4.1": {
		"1.3.0": false,
		"1.3":   false,
		"1.3.2": true,
		"1.3.6": true,
		"1.4.0": true,
		"1.4.1": true,
		"1.4.2": false,
		"1.5.0": false,
	},
	"1.3.2 - 1.4": {
		"1.3.0":  false,
		"1.3":    false,
		"1.3.2":  true,
		"1.3.6":  true,
		"1.4.0":  true,
		"1.4.1":  true,
		"1.4.2":  true,
		"1.4.45": true,
		"1.5.0":  false,
	},
	"1.3.2 - 1": {
		"1.3.0":  false,
		"1.3":    false,
		"1.3.2":  true,
		"1.3.6":  true,
		"1.4.0":  true,
		"1.4.1":  true,
		"1.4.2":  true,
		"1.4.45": true,
		"1.5.0":  true,
		"2.0.0":  false,
	},
	"1.3.*": {
		"1.3.0":    true,
		"1.3.1":    true,
		"1.3.9999": true,
		"1.3":      true,
	},
	"*.*.*": {
		"1.0.2": true,
		"0.0.0": true,
		"0":     true,
	},
	"<*.*.*": {
		"1.0.2": false,
		"0.0.0": false,
		"0":     false,
	},
	"<1.*.*": {
		"2.0.0": false,
		"1.0.2": false,
		"1.0.0": false,
		"0.0.0": true,
		"0":     true,
	},
	">=1.2.7": {
		"1.2.7": true,
		"1.2.8": true,
		"2.5.3": true,
		"1.3.9": true,
		"1.2.6": false,
		"1.1.0": false,
	},
	">*.*.*": {
		"1.0.2": false,
		"0.0.0": false,
		"0":     false,
	},
	"^1.*.*": {
		"1.0.2": true,
		"0.0.0": false,
		"0":     false,
	},
	"*": {
		"1.3.4": true,
		"3.4":   true,
		"0":     true,
		"0.0.1": true,
	},
}

func testRange(t *testing.T, rangeStr string, battery map[string]bool) {
	r := MustParseRange(rangeStr)
	re := r.RegExp()
	for v, result := range battery {
		nv := MustParseVersion(v)
		if r.Contains(nv) != result {
			t.Errorf("Expected %v of %v to evaluate to %v", rangeStr, v, result)
		}
		if re[0].MatchString(nv.String()) &&
			re[1].MatchString(nv.String()) != result {
			t.Errorf("Expected %v of %v to evaluate to %v (Using regexp checks)", rangeStr, v, result)
		}
	}
}

func TestRangeOperator(t *testing.T) {
	for rangeStr, data := range rangeTestBattery {
		testRange(t, rangeStr, data)
	}
	for rangeStr, data := range tildeTests {
		for _, op := range []string{"~", "~>"} {
			testRange(t, fmt.Sprintf("%s%s", op, rangeStr), data)
		}
	}
}
