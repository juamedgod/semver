package semver

import (
	"fmt"
	"regexp"
	"strconv"
)

type Comparable interface {
	UpperLimit() *Version
	LowerLimit() *Version
}
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

func NewVersion(major int64, minor int64, patch int64) *Version {
	return &Version{Major: major, Minor: minor, Patch: patch, majorPresent: true, minorPresent: true, patchPresent: true}
}

func (v *Version) UpperLimit() *Version {
	return v
}

func (v *Version) LowerLimit() *Version {
	return v
}

func (v *Version) Next() *Version {
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
func (v *Version) Compare(v2 *Version) int {
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
func (v *Version) Equal(v2 *Version) bool {
	if v2 == nil {
		return false
	}
	return v.Compare(v2) == 0
}
func (v *Version) Less(v2 *Version) bool {
	if v2 == nil {
		return true
	}
	return v.Compare(v2) < 0
}
func (v *Version) LessOrEqual(v2 *Version) bool {
	return v.Less(v2) || v.Equal(v2)
}
func (v *Version) Greater(v2 *Version) bool {
	if v2 == nil {
		return true
	}
	return v.Compare(v2) > 0
}

func (v *Version) GreaterOrEqual(v2 *Version) bool {
	return v.Greater(v2) || v.Equal(v2)
}

var idStr = `[a-zA-Z0-9-]+(\.[a-zA-Z0-9\.-]*[^\.])?`
var versionRe = regexp.MustCompile(
	fmt.Sprintf(
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
			`)?`, idStr, idStr,
	),
)

var relaxedVersionRe = regexp.MustCompile(
	fmt.Sprintf(
		// Major
		`\s*[=v]*(?P<major>\d[\da-zA-Z]+)`+
			// Followed by an optional .Minor.Patch-preRelease+build
			`([\._-]`+
			// Minor
			`(?P<minor>\d[\da-zA-Z]+)`+
			// Followed by an optional .Patch-preRelease+build
			`([\._-]`+
			// Patch
			`(?P<patch>\d+)`+
			// -preRelease
			`(-?(?P<preRelease>%s))?`+
			// build
			`(\+(?P<build>%s))?`+
			`)?`+
			`)?`, idStr, idStr,
	),
)

func toInt(str string) int64 {
	i, _ := strconv.Atoi(str)
	return int64(i)
}

func parseVersion(str string) (map[string]string, error) {
	result := versionRe.FindStringSubmatch(str)
	mapping := make(map[string]string)

	if len(result) == 0 {
		return mapping, fmt.Errorf("malformed version string")
	}
	for i, name := range versionRe.SubexpNames() {
		if name == "" {
			continue
		}
		mapping[name] = result[i]
	}
	return mapping, nil
}

func ParseVersion(str string) (*Version, error) {
	if mapping, err := parseVersion(str); err == nil {
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
	} else {
		return nil, err
	}
}

func MustParseVersion(str string) *Version {
	if v, err := ParseVersion(str); err != nil {
		panic(err)
	} else {
		return v
	}
}
