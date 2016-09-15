package semver

import "testing"

func TestValid(t *testing.T) {
	//	for verStr, expected := range map[string]bool{} {
	//	if op.Evaluate(MustParseVersion(v)) != result {
	//	t.Errorf("Expected ^%v of %v to evaluate to %v", rangeStr, v, result)
	//	}
	//}
}

func TestLess(t *testing.T) {
	for rangeStr, data := range map[string]map[string]bool{
		">1.3.0": {
			"1.3.0": true,
		},
		"1.3.0": {
			"1.3.0": false,
		},
	} {
		for v, result := range data {
			if Less(MustParseVersion(v), NewRange(rangeStr)) != result {
				t.Errorf("Expected Less(%q, %q) == %q", result)
			}
		}
	}
}
