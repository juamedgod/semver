package semver

// Hack defines a version modification procedure
type Hack func(v *Version) *Version

// WithPreReleaseHandler allows configuring how pre-releases in versions should be handled
func WithPreReleaseHandler(handler func(pre1, pre2 string) (int, error)) Hack {
	return func(v *Version) *Version {
		v.HonorPreRelease(true)
		v.preReleaseComparator = handler
		return v
	}
}

// SupportRevisionsInPreRelease hack make numeric pre-releases to act as revisions
var SupportRevisionsInPreRelease = WithPreReleaseHandler(func(pr1, pr2 string) (res int, err error) {
	if pr1 == pr2 {
		return 0, nil
	}
	// If we have any int, just compare (a number is always greater tha any pre-release indicator)
	if isInt(pr1) || isInt(pr2) {
		return compareInt(int(toInt(pr1)), int(toInt(pr2))), nil
	}
	return comparePreReleases(pr1, pr2)
})
