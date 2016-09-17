package semver

import "testing"

var versionComparissons = map[vPair]int{
	np("1.3", "1.1"):      1,
	np("2.4.2", "2.4.2"):  0,
	np("4.1", "4"):        1,
	np("4.1.1", "4.1"):    1,
	np("3.1.3", "3.1.20"): -1,
	np("0", "1"):          -1,
	np("0.1", "0.1"):      0,
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

func np(v1, v2 string) vPair {
	return vPair{v1: MustParseVersion(v1), v2: MustParseVersion(v2)}
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

// func TestLess(t *testing.T) {
// 	for rangeStr, data := range map[string]map[string]bool{
// 		">1.3.0": {
// 			"1.3.0": true,
// 		},
// 		"1.3.0": {
// 			"1.3.0": false,
// 		},
// 	} {
// 		for v, result := range data {
// 			if Less(MustParseVersion(v), NewRange(rangeStr)) != result {
// 				t.Errorf("Expected Less(%q, %q) == %q", result)
// 			}
// 		}
// 	}
// }
