package semver

import (
	"fmt"
	"testing"
)

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
	"~1.3.4": {
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
	"~1.3": {
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
	"~1": {
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
		"1.3.1": false,
		"2.3":   false,
	},
	"1.3": {
		"1.3.0": true,
		"1.3":   true,
		"1.3.1": false,
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

func TestRangeOperator(t *testing.T) {
	for rangeStr, data := range rangeTestBattery {
		r := MustParseRange(rangeStr)
		re := r.RegExp()
		for v, result := range data {
			if r.Evaluate(MustParseVersion(v)) != result {
				t.Errorf("Expected %v of %v to evaluate to %v", rangeStr, v, result)
			}
			nv := MustParseVersion(v)

			if re[0].MatchString(fmt.Sprintf(`%d.%d.%d`, nv.Major, nv.Minor, nv.Patch)) &&
				re[1].MatchString(fmt.Sprintf(`%d.%d.%d`, nv.Major, nv.Minor, nv.Patch)) != result {
				fmt.Println(re[0], re[1])
				fmt.Println("WTF", fmt.Sprintf(`%d.%d.%d`, nv.Major, nv.Minor, nv.Patch), rangeStr, result)
			}
		}
	}
}

func TestRangeOperatorRe(t *testing.T) {
	for rangeStr, data := range rangeTestBattery {
		r := MustParseRange(rangeStr)
		re := r.RegExp()
		for v, result := range data {
			nv := MustParseVersion(v)
			if re[0].MatchString(fmt.Sprintf(`%d.%d.%d`, nv.Major, nv.Minor, nv.Patch)) &&
				re[1].MatchString(fmt.Sprintf(`%d.%d.%d`, nv.Major, nv.Minor, nv.Patch)) != result {
				t.Errorf("Expected %v of %v to evaluate to %v", rangeStr, v, result)
			}
		}
	}
}
