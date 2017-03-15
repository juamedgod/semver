package semver

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

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

func digitsRange(n int, up bool) [][]string {
	var list [][]string
	for _, d := range digits(n) {
		r := []string{}
		if up {
			for i := d; i < 10; i++ {
				r = append(r, strconv.Itoa(i))
			}
		} else {
			for i := 0; i <= d; i++ {
				r = append([]string{strconv.Itoa(i)}, r...)

			}
		}
		list = append(list, r)
	}
	return list
}

func digits(n int) (r []int) {
	for _, i := range strings.Split(strconv.Itoa(n), "") {
		v, _ := strconv.Atoi(i)
		r = append(r, v)
	}
	return r
}

func gt(n int) string {
	return fmt.Sprintf(`(%s)`, strings.Join(nPatterns(n, true)[1:], `|`))
}

func lt(n int) string {
	return fmt.Sprintf(`(%s)`, strings.Join(nPatterns(n, false)[1:], `|`))
}
func lte(n int) string {
	return fmt.Sprintf(`(%s)`, strings.Join(nPatterns(n, false)[0:], `|`))
}

func gte(n int) string {
	return fmt.Sprintf(`(%s)`, strings.Join(nPatterns(n, true)[0:], `|`))
}
func nPatterns(n int, up bool) []string {
	r := digitsRange(n, up)
	t := []string{strconv.Itoa(n)}

	for i := len(r) - 1; i >= 0; i-- {
		end := r[i]
		str := ""
		for j := 0; j < i; j++ {
			if str == "" && r[j][0] == "0" {
				continue
			}
			str += r[j][0]
		}
		captGrp := strings.Join(end[1:], `|`)
		if captGrp != "" && (str != "" || captGrp != "0") {
			str += `(` + captGrp + `)`
		}
		if extra := len(r) - i - 1; extra > 0 {
			str += fmt.Sprintf(`\d{%d}`, extra)
		}
		if str != "" {
			t = append(t, str)
		} else {
			t = append(t, `(0)`)
		}
	}
	if up {
		t = append(t, fmt.Sprintf(`\d{%d,}`, len(r)+1))
	} else if len(r) > 1 {
		t = append(t, fmt.Sprintf(`\d{,%d}`, len(r)-1))
	}
	return t
}

func isInt(str string) bool {
	if _, err := strconv.Atoi(str); err != nil {
		return false
	}
	return true
}

func toInt(str string) int64 {
	i, _ := strconv.Atoi(str)
	return int64(i)
}

var preReleaseMapping = map[string]*regexp.Regexp{
	"pre-alpha":         regexp.MustCompile(`^(?i)(pre-?alpha)[-\.\_]?(\d*)$`),
	"alpha":             regexp.MustCompile(`^(?i)(alpha|a)[-\.\_]?(\d*)$`),
	"beta":              regexp.MustCompile(`^(?i)(beta|b)[-\.\_]?(\d*)$`),
	"release_candidate": regexp.MustCompile(`^(?i)(rc)[-\.\_]?(\d*)$`),
	"final":             regexp.MustCompile(`^(?i)(final)[-\.\_]?(\d*)$`),
}

func preReleaseQuantifier(pr string) (int, int) {
	quantifier := -1
	level := 0
	for i, id := range []string{
		"pre-alpha", "alpha", "beta", "release_candidate", "final",
	} {
		match := preReleaseMapping[id].FindStringSubmatch(pr)
		if len(match) > 0 {
			quantifier = i
			if match[2] != "" {
				level = int(toInt(match[2]))
			}
			break
		}
	}
	return quantifier, level
}

func compareInt(i1, i2 int) int {
	switch {
	case i1 < i2:
		return -1
	case i1 > i2:
		return 1
	default:
		return 0
	}
}

func comparePreReleases(pr1, pr2 string) (res int, err error) {
	switch {
	case pr1 == pr2:
		res = 0
	case pr1 == "":
		res = 1
	case pr2 == "":
		res = -1
	default:
		q1, l1 := preReleaseQuantifier(pr1)
		q2, l2 := preReleaseQuantifier(pr2)
		switch {
		case q1 == -1 && q2 == -1:
			// Just compare the strings
			if pr1 < pr2 {
				res = -1
			} else {
				res = 1
			}
		case q1 != -1 && q2 != -1:
			res = compareInt(q1, q2)
			if res == 0 {
				res = compareInt(l1, l2)
			}
		default:
			err = fmt.Errorf("unreliable pre-release comparison")
			if pr1 < pr2 {
				return -1, err
			}
			return 1, err
		}
	}
	return res, nil
}
