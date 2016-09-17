package semver

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var rangedVersionRe = regexp.MustCompile(
	fmt.Sprintf(
		// Major
		`\s*[=v]*(?P<major>(\d+|\*|[xX]))`+
			// Followed by an optional .Minor.Patch-preRelease+build
			`(\.`+
			// Minor
			`(?P<minor>(\d+|\*|[xX]))`+
			// Followed by an optional .Patch-preRelease+build
			`(\.`+
			// Patch
			`(?P<patch>(\d+|\*|[xX]))`+
			// -preRelease
			`(-?(?P<preRelease>%s))?`+
			// build
			`(\+(?P<build>%s))?`+
			`)?`+
			`)?`, idStr, idStr,
	),
)

var simpleRangeOperators = []string{`\^`, `~`, `=`, `\<=`, `\>=`, `\<`, `\>`, ``}

var simpleRangeExpr = regexp.MustCompile(fmt.Sprintf(
	`((?P<rangeOp>%s)\s*(?P<version1>%s))`,
	strings.Join(simpleRangeOperators, `|`), rangedVersionRe.String()))

var rangeExpr = regexp.MustCompile(fmt.Sprintf(
	`(?P<range>%s|%s)`,
	hyphenRangeExpr.String(), simpleRangeExpr.String()))

var hyphenRangeExpr = regexp.MustCompile(fmt.Sprintf(
	`((?P<version1>%s)\s*(?P<rangeOp>\-)\s*(?P<version2>%s))`,
	rangedVersionRe.String(), rangedVersionRe.String()))

var versionExpr = regexp.MustCompile(fmt.Sprintf(
	`(?P<version1>%s)`,
	rangedVersionRe.String()))

// GlobVersion defines a version supporting x-range elements
type GlobVersion struct {
	*Version
	anyMajor bool
	anyMinor bool
	anyPatch bool
}

// IsFixed returns true if the x-range does not contain any 'x' elements
func (v *GlobVersion) IsFixed() bool {
	return !v.anyMajor && !v.anyMinor && !v.anyPatch
}

func newGlobVersion(major int64, minor int64, patch int64) *GlobVersion {
	return &GlobVersion{
		anyMajor: major < 0,
		anyMinor: minor < 0,
		anyPatch: patch < 0,
		Version:  NewVersion(major, minor, patch),
	}
}

// TODO: This does not really works with GlobVersion
// Next returns the next version
func (v *GlobVersion) next() *GlobVersion {
	ver := v.Version.next()
	return newGlobVersion(ver.Major, ver.Minor, ver.Patch)
}

// MustParseGlobVersion parses x-range semver from str
// Pnics if it cannot be parsed
func MustParseGlobVersion(str string) *GlobVersion {
	if v, err := ParseGlobVersion(str); err != nil {
		panic(err)
	} else {
		return v
	}
}

// ParseGlobVersion parses x-range semver from str
// Returns an error if it cannot be parsed
func ParseGlobVersion(str string) (*GlobVersion, error) {
	isGlob := func(str string) bool {
		globs := map[string]bool{"*": true, "x": true, "X": true}
		_, ok := globs[str]
		return ok
	}
	mapping, err := parseVersion(str, rangedVersionRe)
	if err != nil {
		return nil, err
	}
	var major, minor, patch int64
	var anyMajor, anyMinor, anyPatch = false, false, false
	if isGlob(mapping["major"]) {
		anyMajor = true
		major = -1
	} else {
		major = toInt(mapping["major"])
	}

	if isGlob(mapping["minor"]) || (mapping["minor"] == "" && major == -1) {
		anyMinor = true
		minor = -1
	} else {
		minor = toInt(mapping["minor"])
		minor = toInt(mapping["minor"])
	}
	if isGlob(mapping["patch"]) || (mapping["patch"] == "" && minor == -1) {
		anyPatch = true
		patch = -1
	} else {
		patch = toInt(mapping["patch"])
	}

	return &GlobVersion{
		anyMajor: anyMajor,
		anyMinor: anyMinor,
		anyPatch: anyPatch,
		Version: &Version{
			Major:        major,
			Minor:        minor,
			Patch:        patch,
			PreRelease:   mapping["preRelease"],
			Build:        mapping["build"],
			majorPresent: mapping["major"] != "",
			minorPresent: mapping["minor"] != "",
			patchPresent: mapping["patch"] != "",
		},
	}, nil
}

// Range defines a semver range
type Range struct {
	Op               string
	MinVersion       *GlobVersion
	AllowMinEquality bool
	MaxVersion       *GlobVersion
	AllowMaxEquality bool
}

var infinity = int64(math.Inf(1))

// RegExp returns a pair of regexps corresponding to the lower and higher limits
// of the range. If a version string matches both regexps, is contained in the range.
func (r *Range) RegExp() []*regexp.Regexp {
	result := make([]*regexp.Regexp, 2)
	if v1 := r.MinVersion; v1 != nil {
		list := []string{
			fmt.Sprintf(`%d\.%d\.%s`, v1.Major, v1.Minor, gt(int(v1.Patch))),
			fmt.Sprintf(`%d\.%s\.\d+`, v1.Major, gt(int(v1.Minor))),
			fmt.Sprintf(`%s\.\d+\.\d+`, gt(int(v1.Major))),
		}
		if r.AllowMinEquality {
			eqlStrComp := []string{}
			for _, c := range []int64{v1.Major, v1.Minor, v1.Patch} {
				if c < 0 {
					eqlStrComp = append(eqlStrComp, `\d*`)
				} else {
					eqlStrComp = append(eqlStrComp, strconv.Itoa(int(c)))
				}
			}
			list = append(list, strings.Join(eqlStrComp, `.`))
		}
		result[0] = regexp.MustCompile(fmt.Sprintf(`(?:(%s))`, strings.Join(list, `|`)))
	} else {
		result[0] = regexp.MustCompile(`.*`)
	}
	if v2 := r.MaxVersion; v2 != nil {
		list := []string{}
		if v2.Patch > 0 {
			list = append(list, fmt.Sprintf(`%d\.%d\.%s`, v2.Major, v2.Minor, lt(int(v2.Patch))))
		}
		if v2.Minor > 0 {
			list = append(list, fmt.Sprintf(`%d\.%s\.\d+`, v2.Major, lt(int(v2.Minor))))
		}

		if v2.Major > 0 {
			list = append(list, fmt.Sprintf(`%s\.\d+\.\d+`, lt(int(v2.Major))))
		} else if v2.Major < 0 {
			// *.\d.\d is in-matcheable
			list = append(list, fmt.Sprintf(`x\.x\.x`))
		}
		if r.AllowMaxEquality {
			eqlStrComp := []string{}
			for _, c := range []int64{v2.Major, v2.Minor, v2.Patch} {
				if c < 0 {
					eqlStrComp = append(eqlStrComp, `\d*`)
				} else {
					eqlStrComp = append(eqlStrComp, strconv.Itoa(int(c)))
				}
			}
			list = append(list, strings.Join(eqlStrComp, `.`))
		}
		result[1] = regexp.MustCompile(fmt.Sprintf(`(?:(%s))`, strings.Join(list, `|`)))
	} else {
		result[1] = regexp.MustCompile(`.*`)
	}
	return result
}

// UpperLimit returns a version describing the upper limit of the range
func (r *Range) UpperLimit() *GlobVersion {
	mVersion := r.MaxVersion
	if mVersion == nil {
		return newGlobVersion(infinity, infinity, infinity)
	}
	if r.AllowMaxEquality {
		return mVersion
	}
	return newGlobVersion(mVersion.Major, mVersion.Minor-1, -infinity)
}

// LowerLimit returns a version describing the lower limit of the range
func (r *Range) LowerLimit() *GlobVersion {
	minVersion := r.MinVersion
	if minVersion == nil {
		return newGlobVersion(-infinity, -infinity, -infinity)
	}
	if r.AllowMinEquality {
		return minVersion
	}
	return newGlobVersion(minVersion.Major, minVersion.Minor+1, -infinity)
}

// MustParseRange creates a Range from a semver string
// It panics in case of error
func MustParseRange(str string) *Range {
	r, err := ParseRange(str)
	if err != nil {
		panic(err)
	}
	return r
}

// ParseRange creates a Range from a semver string
// It will return a non-nil error if it fails
func ParseRange(str string) (*Range, error) {
	matched, mapping := namedReEvaluate(rangeExpr, str)
	if !matched {
		return nil, fmt.Errorf("Malformed range expression %q", str)
	}
	v := MustParseGlobVersion(mapping["version1"])
	op := &Range{
		Op: mapping["rangeOp"],
	}
	var maxVersion, minVersion *GlobVersion

	minVersion = v
	op.AllowMinEquality = true
	op.AllowMaxEquality = false

	switch op.Op {
	case `-`:
		v2 := MustParseGlobVersion(mapping["version2"])
		op.AllowMinEquality = true

		if v2.patchPresent {
			maxVersion = v2
			op.AllowMaxEquality = true
		} else {
			maxVersion = v2.next()
			op.AllowMaxEquality = false
		}
	case `^`:
		d := make([]int64, 3)
		for i, c := range v.split() {
			if c != 0 {
				d[i] = int64(c) + 1
				break
			}
		}
		maxVersion = newGlobVersion(d[0], d[1], d[2])
		op.AllowMaxEquality = false
	case `>`:
		minVersion = v
		op.AllowMinEquality = false
		maxVersion = nil
	case `>=`:
		maxVersion = nil
	case `<`:
		minVersion = nil
		op.AllowMaxEquality = false
		maxVersion = v
	case `<=`:
		minVersion = nil
		op.AllowMaxEquality = true
	case `=`:
		fallthrough
	case ``:
		minVersion = v
		maxVersion = v
		op.AllowMaxEquality = true
		op.AllowMinEquality = true
	case `~`:
		switch {
		case v.minorPresent:
			maxVersion = newGlobVersion(v.Major, v.Minor+1, 0)
		default:
			maxVersion = newGlobVersion(v.Major+1, 0, 0)
		}
	default:
		return nil, fmt.Errorf(`Unknown range operator ` + op.Op)
	}
	op.MaxVersion = maxVersion
	op.MinVersion = minVersion
	return op, nil
}

// Contains checks if the provided version v is contained by the Range
func (r *Range) Contains(v *Version) bool {
	if (v.Greater(r.MinVersion) || (r.AllowMinEquality && v.Equal(r.MinVersion))) &&
		(v.Less(r.MaxVersion) || (r.AllowMaxEquality && v.Equal(r.MaxVersion))) {
		return true
	}
	return false
}

// Matches is equivalent to Contains and is provided to satisfy the Expression interface (a range is a simple expression)
func (r *Range) Matches(v *Version) bool {
	return r.Contains(v)
}

func (r *Range) evaluate(v *Version) bool {
	return r.Contains(v)
}
