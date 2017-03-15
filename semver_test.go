package semver

import "testing"

const HonorPreRelease = true

var versionComparissons = map[vPair]int{
	np("1.3", "1.1"):                          1,
	np("2.4.2", "2.4.2"):                      0,
	np("4.1", "4"):                            1,
	np("4.1.1", "4.1"):                        1,
	np("3.1.3", "3.1.20"):                     -1,
	np("0", "1"):                              -1,
	np("0.1", "0.1"):                          0,
	np("1.3.0-0", "1.3.0-1"):                  0,
	np("1.3.0-0", "1.3.0-1", HonorPreRelease): -1,
	//	np("1.2.4", "1.*"):    -1,
}

func TestValid(t *testing.T) {
	for verStr, v := range versionTestBattery {
		isValid := v != nil
		if Valid(verStr) != isValid {
			t.Errorf("Expected Valid(%q) to be %v", verStr, isValid)
		}
	}
}

type vPair struct {
	v1 Comparable
	v2 Comparable
}

func np(v1, v2 string, args ...bool) vPair {
	semver1 := MustParseVersion(v1)
	semver2 := MustParseVersion(v2)

	if len(args) > 0 {
		semver1.HonorPreRelease(args[0])
		semver2.HonorPreRelease(args[0])
	}

	return vPair{v1: semver1, v2: semver2}
}

func TestGreater(t *testing.T) {
	for pair, comparisson := range versionComparissons {
		greater := comparisson > 0
		if Greater(pair.v1, pair.v2) != greater {
			t.Errorf("Expected Greater(%q, %q) to be %v", pair.v1, pair.v2, greater)
		}
		// If they are not equal, we reverse it to also test the other combination
		if comparisson != 0 {
			if Greater(pair.v2, pair.v1) != !greater {
				t.Errorf("Expected Greater(%q, %q) to be %v", pair.v2, pair.v1, !greater)
			}
		}
	}
}

func TestGreaterOrEqual(t *testing.T) {
	for pair, comparisson := range versionComparissons {
		greaterOrEqual := comparisson >= 0
		if GreaterOrEqual(pair.v1, pair.v2) != greaterOrEqual {
			t.Errorf("Expected GreaterOrEqual(%q, %q) to be %v", pair.v1, pair.v2, greaterOrEqual)
		}
		// If they are not equal, we reverse it to also test the other combination
		if comparisson != 0 {
			if GreaterOrEqual(pair.v2, pair.v1) != !greaterOrEqual {
				t.Errorf("Expected GreaterOrEqual(%q, %q) to be %v", pair.v2, pair.v1, !greaterOrEqual)
			}
		}
	}
}

func TestLess(t *testing.T) {
	for pair, comparisson := range versionComparissons {
		less := comparisson < 0
		if Less(pair.v1, pair.v2) != less {
			t.Errorf("Expected Less(%q, %q) to be %v", pair.v1, pair.v2, less)
		}
		// If they are not equal, we reverse it to also test the other combination
		if comparisson != 0 {
			if Less(pair.v2, pair.v1) != !less {
				t.Errorf("Expected Less(%q, %q) to be %v", pair.v2, pair.v1, !less)
			}
		}
	}
}

func TestLessOrEqual(t *testing.T) {
	for pair, comparisson := range versionComparissons {
		lessOrEqual := comparisson <= 0
		if LessOrEqual(pair.v1, pair.v2) != lessOrEqual {
			t.Errorf("Expected LessOrEqual(%q, %q) to be %v", pair.v1, pair.v2, lessOrEqual)
		}
		// If they are not equal, we reverse it to also test the other combination
		if comparisson != 0 {
			if LessOrEqual(pair.v2, pair.v1) != !lessOrEqual {
				t.Errorf("Expected LessOrEqual(%q, %q) to be %v", pair.v2, pair.v1, !lessOrEqual)
			}
		}
	}
}

func TestEqual(t *testing.T) {
	for pair, comparisson := range versionComparissons {
		equal := comparisson == 0
		if Equal(pair.v1, pair.v2) != equal {
			t.Errorf("Expected Equal(%q, %q) to be %v", pair.v1, pair.v2, equal)
		}
		if Equal(pair.v2, pair.v1) != equal {
			t.Errorf("Expected Equal(%q, %q) to be %v", pair.v2, pair.v1, equal)
		}
	}
}

func TestSatisfies(t *testing.T) {
	for _, battery := range []map[string]map[string]bool{
		rangeTestBattery,
		exprTestBattery,
	} {
		for exprStr, data := range battery {
			exprObj := MustParseExpr(exprStr)
			for vStr, expected := range data {
				vObj := MustParseVersion(vStr)
				results := make([]bool, 4)
				errors := make([]error, 4)
				results[0], errors[0] = Satisfies(vObj, exprObj)
				results[1], errors[1] = Satisfies(vStr, exprStr)
				results[2], errors[2] = Satisfies(vStr, exprObj)
				results[3], errors[3] = Satisfies(vObj, exprStr)
				for i := 0; i < len(results); i++ {
					if errors[i] != nil {
						t.Errorf("Expected Satisfies(%q, %q) (combination %d) to succeed but got %v", exprStr, vStr, i, errors[i])
					}
					if results[i] != expected {
						t.Errorf("Expected Satisfies(%q, %q) (combination %d) to be %v", exprStr, vStr, i, expected)
					}
				}
			}
		}
	}

}
