package semver

import "regexp"

func namedReEvaluate(re *regexp.Regexp, str string) (matched bool, mapping map[string]string) {
	result := re.FindStringSubmatch(str)
	mapping = make(map[string]string)

	if len(result) == 0 {
		return false, mapping
	}

	for i, name := range re.SubexpNames() {
		if name == "" {
			continue
		}
		// Don't overwrite already detected elements
		if v, ok := mapping[name]; ok && v != "" {
			continue
		}
		mapping[name] = result[i]
	}
	return true, mapping
}
