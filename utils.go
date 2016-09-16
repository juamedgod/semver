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
	list := make([][]string, 0)
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
