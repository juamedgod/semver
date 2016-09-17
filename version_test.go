package semver

import "testing"

var versionTestBattery = map[string]*Version{
	"1":                NewVersion(1, 0, 0),
	"1.3":              NewVersion(1, 3, 0),
	"1.2.1":            NewVersion(1, 2, 1),
	"  4.3.2  ":        NewVersion(4, 3, 2),
	"1.2.3-alpha.3":    NewVersion(1, 2, 3, "alpha.3"),
	"1.4.5+2342":       NewVersion(1, 4, 5, "", "2342"),
	"1.2.3-rc.4+34.a2": NewVersion(1, 2, 3, "rc.4", "34.a2"),
	"0.0.0":            NewVersion(0, 0, 0),
	"..":               nil,
	"a.1.a":            nil,
	"a1.2.3a":          nil,
	"3.5.2  xxx":       nil,
	"foo":              nil,
	"":                 nil,
}

func TestParseVersion(t *testing.T) {
	for verStr, expected := range versionTestBattery {
		v, err := ParseVersion(verStr)
		if expected == nil {
			if err == nil || v != nil {
				t.Errorf("Expected %q to not be parseable", verStr)
			}
		} else {
			if expected.Major != v.Major ||
				expected.Minor != v.Minor ||
				expected.Patch != v.Patch ||
				expected.PreRelease != v.PreRelease ||
				expected.Build != v.Build {
				t.Errorf("Parsed version %q does not match expected %q", v, expected)
			}

		}
	}
}
