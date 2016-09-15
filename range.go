package semver

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

type Range struct {
	Op               string
	MinVersion       *Version
	AllowMinEquality bool
	MaxVersion       *Version
	AllowMaxEquality bool
}

var infinity = int64(math.Inf(1))

func (r *Range) UpperLimit() *Version {
	mVersion := r.MaxVersion
	if mVersion == nil {
		return NewVersion(infinity, infinity, infinity)
	}
	if r.AllowMaxEquality {
		return mVersion
	}
	return NewVersion(mVersion.Major, mVersion.Minor-1, -infinity)
}

func (r *Range) LowerLimit() *Version {
	minVersion := r.MinVersion
	if minVersion == nil {
		return NewVersion(-infinity, -infinity, -infinity)
	}
	if r.AllowMinEquality {
		return minVersion
	}
	return NewVersion(minVersion.Major, minVersion.Minor+1, -infinity)
}

var simpleRangeOperators = []string{`\^`, `~`, `=`, `\<=`, `\>=`, `\<`, `\>`, ``}

var simpleRangeExpr = regexp.MustCompile(fmt.Sprintf(
	`((?P<rangeOp>%s)\s*(?P<version1>%s))`,
	strings.Join(simpleRangeOperators, `|`), versionRe.String()))

var rangeExpr = regexp.MustCompile(fmt.Sprintf(
	`(?P<range>%s|%s)`,
	hyphenRangeExpr.String(), simpleRangeExpr.String()))

var hyphenRangeExpr = regexp.MustCompile(fmt.Sprintf(
	`((?P<version1>%s)\s*(?P<rangeOp>\-)\s*(?P<version2>%s))`,
	versionRe.String(), versionRe.String()))

var versionExpr = regexp.MustCompile(fmt.Sprintf(
	`(?P<version1>%s)`,
	versionRe.String()))

func NewRange(str string) *Range {
	matched, mapping := namedReEvaluate(rangeExpr, str)
	if !matched {
		return nil
	}
	v := MustParseVersion(mapping["version1"])
	op := &Range{
		Op: mapping["rangeOp"],
	}
	var maxVersion, minVersion *Version

	minVersion = v
	op.AllowMinEquality = true
	op.AllowMaxEquality = false

	switch op.Op {
	case `-`:
		v2 := MustParseVersion(mapping["version2"])
		op.AllowMinEquality = true

		if v2.patchPresent {
			maxVersion = v2
			op.AllowMaxEquality = true
		} else {
			maxVersion = v2.Next()
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
		maxVersion = &Version{Major: d[0], Minor: d[1], Patch: d[2]}
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
			maxVersion = &Version{Major: v.Major, Minor: v.Minor + 1, Patch: 0}
		default:
			maxVersion = &Version{Major: v.Major + 1, Minor: 0, Patch: 0}
		}
	default:
		panic(`Unknown operator ` + op.Op)
	}
	op.MaxVersion = maxVersion
	op.MinVersion = minVersion
	return op
}

func (r *Range) Evaluate(v *Version) bool {
	if (v.Greater(r.MinVersion) || (r.AllowMinEquality && v.Equal(r.MinVersion))) &&
		(v.Less(r.MaxVersion) || (r.AllowMaxEquality && v.Equal(r.MaxVersion))) {
		return true
	}
	return false
}
