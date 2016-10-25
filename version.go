package semver

import (
	"fmt"
	"regexp"
	"encoding/json"
)

var idStr = `[a-zA-Z0-9-]+(\.[a-zA-Z0-9\.-]*[^\.])?`
var versionRe = regexp.MustCompile(
	fmt.Sprintf(
		`^`+
			// Major
			`\s*[v=]*(?P<major>\d+)`+
			// Followed by an optional .Minor.Patch-preRelease+build
			`(\.`+
			// Minor
			`(?P<minor>\d+)`+
			// Followed by an optional .Patch-preRelease+build
			`(\.`+
			// Patch
			`(?P<patch>\d+)`+
			// -preRelease
			`(-(?P<preRelease>%s))?`+
			// build
			`(\+(?P<build>%s))?`+
			`)?`+
			`)?`+
			`\s*$`,
		idStr, idStr,
	),
)

var relaxedVersionRe = regexp.MustCompile(
	fmt.Sprintf(
		// Major
		`\s*[=v]*(?P<major>\d[\da-zA-Z]*)`+
			// Followed by an optional .Minor.Patch-preRelease+build
			`([\._-]`+
			// Minor
			`(?P<minor>\d[\da-zA-Z]*)?`+
			// Followed by an optional .Patch-preRelease+build
			`([\._-]`+
			// Patch
			`(?P<patch>\d[\da-zA-Z]*)?`+
			// -preRelease
			`(-?(?P<preRelease>%s))?`+
			// build
			`(\+(?P<build>%s))?`+
			`)?`+
			`)?`, idStr, idStr,
	),
)

// Comparable defines the interface required to be compared
type Comparable interface {
}

// Version describes a semver
type Version struct {
	Major      int64
	Minor      int64
	Patch      int64
	PreRelease string
	Build      string

	majorPresent bool
	minorPresent bool
	patchPresent bool
}

func (v *Version) MarshalJSON() (data []byte, err error) {
	return json.Marshal(v.String())
}

func (v *Version) String() string {
	if v == nil {
		return ""
	}
	s := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.PreRelease != "" {
		s += `-` + v.PreRelease
	}
	if v.Build != "" {
		s += `+` + v.Build
	}
	return s
}

// NewVersion returns a version from the specified major, minor and patch components
func NewVersion(major int64, minor int64, patch int64, extra ...string) *Version {
	var preRelease, build string
	if len(extra) > 0 {
		preRelease = extra[0]
		if len(extra) > 1 {
			build = extra[1]
		}
	}
	return &Version{
		Major: major, Minor: minor, Patch: patch,
		majorPresent: true, minorPresent: true, patchPresent: true,
		PreRelease: preRelease, Build: build,
	}
}

func (v *Version) next() *Version {
	if v.patchPresent {
		return NewVersion(v.Major, v.Minor, v.Patch+1)
	}
	if v.minorPresent {
		return NewVersion(v.Major, v.Minor+1, 0)
	}
	return NewVersion(v.Major+1, 0, 0)
}
func (v *Version) split() []int64 {
	return []int64{v.Major, v.Minor, v.Patch}
}
func (v *Version) compare(v2 *Version) int {
	sv1 := v.split()
	sv2 := v2.split()

	for i := 0; i < 3; i++ {
		e1 := sv1[i]
		e2 := sv2[i]
		if e1 > e2 {
			return 1
		}
		if e1 < e2 {
			return -1
		}
	}
	return 0
}

func (v *Version) equal(v2 *Version) bool {
	if v2 == nil {
		return false
	}
	return v.compare(v2) == 0
}
func (v *Version) less(v2 *Version) bool {
	if v2 == nil {
		return true
	}
	return v.compare(v2) < 0
}

// LessOrEqual checks if v is less or equal than the provided version v2
func (v *Version) LessOrEqual(v2 Comparable) bool {
	return v.Less(v2) || v.Equal(v2)
}
func (v *Version) greater(v2 *Version) bool {
	if v2 == nil {
		return true
	}
	return v.compare(v2) > 0
}

// GreaterOrEqual checks if v is greater or equal than the provided version v2
func (v *Version) GreaterOrEqual(v2 Comparable) bool {
	return v.Greater(v2) || v.Equal(v2)
}

// Equal checks if v is equal to the provided version v2
func (v *Version) Equal(v2 Comparable) bool {
	return v.compareWith(v2, func(pos int, e1 int64, e2 int64, any bool) interface{} {
		if e1 != e2 && !any {
			return false
		}
		return nil
	}, v.equal)
}

// Less checks if v is less than the provided version v2
func (v *Version) Less(v2 Comparable) bool {

	if plainVersion, ok := v2.(*GlobVersion); (ok && plainVersion == nil) || v2 == nil {
		return true
	}
	return v.compareWith(v2, func(pos int, e1 int64, e2 int64, any bool) interface{} {
		if any {
			return false
		}
		if e1 < e2 {
			return true
		}
		if e1 > e2 {
			return false
		} else if e1 == e2 {
			if pos == 2 {
				return false
			}
		}
		return nil
	}, v.less)
}
func compareWithGlobVersion(v1 *Version, v2 *GlobVersion, fn func(pos int, e1 int64, e2 int64, any bool) interface{}) bool {
	for i, set := range []struct {
		e2, e1 int64
		any    bool
	}{
		{e2: v2.Major, e1: v1.Major, any: v2.anyMajor},
		{e2: v2.Minor, e1: v1.Minor, any: v2.anyMinor},
		{e2: v2.Patch, e1: v1.Patch, any: v2.anyPatch},
	} {
		v := fn(i, set.e1, set.e2, set.any)
		if res, ok := v.(bool); ok {
			return res
		}
	}
	return true
}

func (v *Version) compareWith(item Comparable, globComparer func(pos int, e1 int64, e2 int64, any bool) interface{}, regComparer func(v2 *Version) bool) bool {
	switch v2 := item.(type) {
	case *GlobVersion:
		if v2 == nil {

			return false
		}
		if v2.IsFixed() {
			return regComparer(v2.Version)
		}
		return compareWithGlobVersion(v, v2, globComparer)
	case *Version:
		return regComparer(v2)
	default:
		panic(fmt.Errorf(`Unsupported element type %T`, item))
	}
}

// Greater checks if v is greater than the provided version v2
func (v *Version) Greater(v2 Comparable) bool {
	if plainVersion, ok := v2.(*GlobVersion); (ok && plainVersion == nil) || v2 == nil {
		return true
	}
	return v.compareWith(v2, func(pos int, e1 int64, e2 int64, any bool) interface{} {
		if any {
			return false
		}
		if e1 < e2 {
			return false
		} else if e1 > e2 {
			return true
		} else if pos == 2 {
			return false
		}
		return nil
	}, v.greater)
}

func parseVersion(str string, re *regexp.Regexp) (map[string]string, error) {
	result := re.FindStringSubmatch(str)
	mapping := make(map[string]string)

	if len(result) == 0 {
		return mapping, fmt.Errorf("malformed version string")
	}
	for i, name := range re.SubexpNames() {
		if name == "" {
			continue
		}
		mapping[name] = result[i]
	}
	return mapping, nil
}

// ParseVersion parse a semver str and returns a Version.
// Returns an error if it fails to parse.
func ParseVersion(str string) (*Version, error) {
	mapping, err := parseVersion(str, versionRe)
	if err != nil {
		return nil, err
	}
	return &Version{
		Major:        toInt(mapping["major"]),
		Minor:        toInt(mapping["minor"]),
		Patch:        toInt(mapping["patch"]),
		PreRelease:   mapping["preRelease"],
		Build:        mapping["build"],
		majorPresent: mapping["major"] != "",
		minorPresent: mapping["minor"] != "",
		patchPresent: mapping["patch"] != "",
	}, nil
}

// ParsePermissiveVersion relaxes the Semver requirements to be able to parse non
// standard versions
func ParsePermissiveVersion(str string) (v *Version, err error) {
	if v, err = ParseVersion(str); err == nil {
		return v, err
	}

	mapping, err := parseVersion(str, relaxedVersionRe)
	if err != nil {
		return nil, err
	}
	var major, minor, patch int64

	fmt.Sscanf(mapping["major"], "%d", &major)
	fmt.Sscanf(mapping["minor"], "%d", &minor)
	fmt.Sscanf(mapping["patch"], "%d", &patch)
	return &Version{
		Major:        major,
		Minor:        minor,
		Patch:        patch,
		PreRelease:   mapping["preRelease"],
		Build:        mapping["build"],
		majorPresent: mapping["major"] != "",
		minorPresent: mapping["minor"] != "",
		patchPresent: mapping["patch"] != "",
	}, nil
}

// MustParseVersion parse a semver str and returns a Version.
// It panics if it fails to parse it.
func MustParseVersion(str string) *Version {
	if v, err := ParseVersion(str); err != nil {
		panic(err)
	} else {
		return v
	}
}

// UpperLimit returns the upper limit of a version (itself).
// Allows treating a version as a range for comparisson
func (v *Version) UpperLimit() *Version {
	return v
}

// LowerLimit returns the lower limit of a version (itself).
// Allows treating a version as a range for comparisson
func (v *Version) LowerLimit() *Version {
	return v
}
